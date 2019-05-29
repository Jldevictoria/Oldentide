// Filename:    dedicated_server.go
// Author:      Joseph DeVictoria
// Date:        June_13_2018
// Purpose:     The dedicated game server for Oldentide written in Go.

package main

import (
	"Oldentide/shared"
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/vmihailenco/msgpack"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/smtp"
	"os"
	_ "path/filepath"
	"runtime"
	"strconv"
	"time"
)

// Global program variables.
var err error
var gport int
var wport int
var everify bool
var webadd string
var email string
var epass string
var dbpath string
var eauth smtp.Auth
var db *sql.DB
var packet_count int

func init() {
	flag.IntVar(&gport, "gport", 0, "Port used for dedicated game server.")
	flag.IntVar(&wport, "wport", 0, "Port used for website.")
	flag.BoolVar(&everify, "everify", false, "Use an emailer to verify accounts?")
	flag.StringVar(&webadd, "webadd", "", "Public website root address where accounts will be created.")
	flag.StringVar(&email, "email", "", "Gmail email address used to send verification emails.")
	flag.StringVar(&epass, "epass", "", "Gmail email password used to send verification emails.")
	flag.StringVar(&dbpath, "dbpath", shared.DefaultGOPATH()+"/src/Oldentide/db/oldentide.db", "Path to oldentide.db")
	rand.Seed(time.Now().UTC().UnixNano())
}

func main() {
	// Extract command line input.
	flag.Parse()
	fmt.Println("\n-------------------------------------------------------")
	fmt.Println("Server Configurations from command line:")
	fmt.Println("-------------------------------------------------------")
	fmt.Println("gport:", gport)
	fmt.Println("wport:", wport)
	fmt.Println("webadd:", webadd)
	fmt.Println("everify:", everify)
	fmt.Println("email:", email)
	fmt.Println("epass:", epass)
	fmt.Println("dbpath:", dbpath)
	if gport == 0 {
		log.Fatal("Please provide a game port with the command line flag -gport=<number>")
	}
	if wport == 0 {
		log.Fatal("Please provide a website port with the command line flag -wport=<number>")
	}
	if webadd == "" {
		log.Fatal("Please provide the website address/ip with the command line flag -webadd=<www.address.domain>")
	}
	if everify {
		if email == "" {
			log.Fatal("Please provide a Gmail email account with the command line flag -email=<email@gmail.com>")
		}
		if epass == "" {
			log.Fatal("Please provide a Gmail email password with the command line flag -epass=<P@55word>")
		}
	} else {
		fmt.Println("Warning: website allowing account creation without email verification!")
		fmt.Println("To enable email verification please use the command line flag -everify")
	}
	eauth = smtp.PlainAuth("", email, epass, "smtp.gmail.com")
	fmt.Println("\n-------------------------------------------------------")
	fmt.Println("Starting Oldentide dedicated server!")
	fmt.Println("-------------------------------------------------------")

	// --------------------------------------------------------------------------------------------
	// Opening database.
	// --------------------------------------------------------------------------------------------
	_, err := os.Stat(dbpath)
	if err != nil {
		log.Fatal("Couldn't find a database file at: " + dbpath)
	}
	db, err = sql.Open("sqlite3", dbpath)
	shared.CheckErr(err)
	fmt.Println("* Database connected.\n")

	// Initialize the game state (populates all of the npcs, and game objects, etc).
	// --------------------------------------------------------------------------------------------
	race_templates := pullRaceTemplates()
	fmt.Println("\n* Race templates populated from database.\n")
	for _, race_template := range race_templates {
		fmt.Println(race_template)
	}

	profession_templates := pullProfessionTemplates()
	fmt.Println("\n* Profession templates populated from database.\n")
	for _, profession_template := range profession_templates {
		fmt.Println(profession_template)
	}

	item_templates := pullItemTemplates()
	fmt.Println("\n* Item templates populated from database:\n")
	for _, item_template := range item_templates {
		fmt.Println(item_template)
	}

	pcs := pullPcs()
	fmt.Println("* PCs listed in database:\n")
	for _, pc := range pcs {
		fmt.Println(pc)
	}

	npcs := pullNpcs()
	fmt.Println("\n* NPCs populated from database:\n")
	for _, npc := range npcs {
		fmt.Println(npc)
	}

	// inventories := pullInventories()

	// --------------------------------------------------------------------------------------------
	// Kick off http server for registration page.
	// --------------------------------------------------------------------------------------------
	mux := http.NewServeMux()
	mux.HandleFunc("/", routeWebTraffic)
	go http.ListenAndServe(":"+strconv.Itoa(wport), mux)

	// --------------------------------------------------------------------------------------------
	// Create and bind a udp socket descriptor.
	// --------------------------------------------------------------------------------------------
	server_address := net.UDPAddr{
		IP:   net.IP{0, 0, 0, 0},
		Port: gport,
	}
	socket, err := net.ListenUDP("udp", &server_address)
	shared.CheckErr(err)

	// --------------------------------------------------------------------------------------------
	// Start our collecter to pull in packets from the hardware socket.
	// --------------------------------------------------------------------------------------------
	RawPacketQueue := make(chan shared.Raw_packet, 100000)
	QuitChan := make(chan bool)
	go Collect(socket, RawPacketQueue, QuitChan)
	fmt.Println("\n* Collector Launched.\n")

	// --------------------------------------------------------------------------------------------
	// Start as many handler goroutines as we have processors.
	// Should result in non-thrashing handler thread concurrency.
	// --------------------------------------------------------------------------------------------
	for i := 0; i < runtime.NumCPU(); i++ {
		go Handle(RawPacketQueue, QuitChan, i)
		fmt.Println("* Handler ", i+1, " of ", runtime.NumCPU(), " Launched.")
	}

	// Wait for a message to the Quit Channel to halt execution.
	<-QuitChan

	// Close database.
	db.Close()
}

// Places all UDP packets that arrive on the hardware socket into a queue for handling.
func Collect(connection *net.UDPConn, RawPacketQueue chan shared.Raw_packet, QuitChan chan bool) {
	for {
		buffer := make([]byte, 512) //65507) // Max IPv4 UDP packet size.
		n, remote_address, err := connection.ReadFromUDP(buffer)
		shared.CheckErr(err)
		RawPacketQueue <- shared.Raw_packet{n, remote_address, buffer}
		packet_count++
		fmt.Println("PC:", packet_count)
	}
}

// Handle all arriving packets based on which opcode they are.
func Handle(RawPacketQueue chan shared.Raw_packet, QuitChan chan bool, rid int) {
	var packet shared.Raw_packet
	for {
		select {
		// This case will run when there is a packet available at the front of the packet queue.
		case packet = <-RawPacketQueue:
			//fmt.Println("Goroutine ID:", rid, "Size:", packet.Size, "Sender:", packet.Client, "Payload:", packet.Payload[:packet.Size]) //debug
			var op shared.Opcode_packet
			err = msgpack.Unmarshal(packet.Payload, &op)
			shared.IfErrPrintErr(err)
			// Depending on what packet opcode we recieved, handle the data accordingly.
			switch op.Opcode {
			case shared.EMPTY:
				fmt.Println("Handling an EMPTY packet.")
				continue
			case shared.GENERIC:
				// fmt.Println("Handling a GENERIC packet.")
				continue
			case shared.ACK:
				fmt.Println("Handling an ACK packet.")
				var decpac shared.Ack_packet
				err = msgpack.Unmarshal(packet.Payload, &decpac)
				shared.CheckErr(err)
				continue
			case shared.ERROR:
				fmt.Println("Handling an ERROR packet.")
				var decpac shared.Error_packet
				err = msgpack.Unmarshal(packet.Payload, &decpac)
				shared.CheckErr(err)
				continue
			case shared.REQCLIST:
				fmt.Println("Handling a REQCLIST packet.")
				var decpac shared.Req_clist_packet
				err = msgpack.Unmarshal(packet.Payload, &decpac)
				fmt.Println(decpac)
				var retpac shared.Send_clist_packet
				retpac.Opcode = shared.SENDCLIST
				retpac.Characters = getCharacterList(decpac.Account)
				fmt.Println(retpac)
				continue
			case shared.CREATEPLAYER:
				fmt.Println("Handling a CREATEPLAYER packet.")
				var decpac shared.Create_player_packet
				err = msgpack.Unmarshal(packet.Payload, &decpac)
				fmt.Println(decpac)
				// Need to get the account name by session id.
				account_name := "test"
				player_name := decpac.Pc.Firstname
				decpac.Pc.Accountid = getAccountIdFromAccountName(account_name)
				fmt.Println(decpac.Pc.Accountid)
				if getRemainingPlayerSlots(account_name, 10) == 0 {
					log.Println("Account tried to create too many players.")
					continue
				}
				if playerFirstNameTaken(player_name) {
					log.Println("Account tried to create a player whose name was already taken.")
					continue
				}
				if validNewPlayer(decpac.Pc) {
					addNewPlayer(decpac.Pc)
					log.Println("Account <account> created a new player \"<player>\".")
				} else {
					log.Println("Account is trying something fraudulent during account creation!")
					//banAccount()
					continue
				}
			case shared.CONNECT:
				fmt.Println("Handling a CONNECT packet.")
				var decpac shared.Connect_packet
				err = msgpack.Unmarshal(packet.Payload, &decpac)
				shared.CheckErr(err)
				if p, ok := players[decpac.Session_id]; ok {
					fmt.Println("Player", *p, "tried to connect twice.  Forcing full disconnect.")
					delete(players, decpac.Session_id)
					// Send a force disconnect packet.
				} else {
					fmt.Println("I should be adding a player to the map here...")
				}
				continue
			case shared.DISCONNECT:
				fmt.Println("Handling a DISCONNECT packet.")
				var decpac shared.Disconnect_packet
				err = msgpack.Unmarshal(packet.Payload, &decpac)
				shared.CheckErr(err)
				if _, ok := players[decpac.Session_id]; ok {
					delete(players, decpac.Session_id)
					// Send a force disconnect packet.
				} else {
					fmt.Println("Player with session", decpac.Session_id, "tried to disconnect, but was never connected...")
				}
				continue
			case shared.MOVEPLAYER:
				fmt.Println("Handling a MOVEPLAYER packet.")
				var decpac shared.Move_player_packet
				err = msgpack.Unmarshal(packet.Payload, &decpac)
				shared.CheckErr(err)
				var player *shared.Pc
				players[decpac.Session_id] = player
				if player, ok := players[decpac.Session_id]; ok {
					fmt.Println(*player)
					player.X = decpac.X
					player.Y = decpac.Y
					player.Z = decpac.Z
					player.Direction = decpac.Direction
					fmt.Println(*player)
				} else {
					fmt.Println("Player did not exist in MOVEPLAYER case for session", decpac.Session_id)
				}
				continue
			case shared.SPENDDP:
				fmt.Println("Handling a SPENDDP packet.")
				var decpac shared.Spend_dp_packet
				err = msgpack.Unmarshal(packet.Payload, &decpac)
				shared.CheckErr(err)
				continue
			case shared.TALKCMD:
				fmt.Println("Handling a TALKCMD packet.")
				var decpac shared.Talk_cmd_packet
				err = msgpack.Unmarshal(packet.Payload, &decpac)
				shared.CheckErr(err)
				continue
			case shared.ATTACKCMD:
				fmt.Println("Handling a ATTACKCMD packet.")
				var decpac shared.Attack_cmd_packet
				err = msgpack.Unmarshal(packet.Payload, &decpac)
				shared.CheckErr(err)
				continue
			case shared.TRADECMD:
				fmt.Println("Handling a TRADECMD packet.")
				var decpac shared.Trade_cmd_packet
				err = msgpack.Unmarshal(packet.Payload, &decpac)
				shared.CheckErr(err)
				continue
			case shared.INVITECMD:
				fmt.Println("Handling a INVITECMD packet.")
				var decpac shared.Invite_cmd_packet
				err = msgpack.Unmarshal(packet.Payload, &decpac)
				shared.CheckErr(err)
				continue
			case shared.GINVITECMD:
				fmt.Println("Handling a GINVITECMD packet.")
				var decpac shared.Guild_invite_cmd_packet
				err = msgpack.Unmarshal(packet.Payload, &decpac)
				shared.CheckErr(err)
				continue
			case shared.GKICK:
				fmt.Println("Handling a GKICK packet.")
				var decpac shared.Guild_kick_cmd_packet
				err = msgpack.Unmarshal(packet.Payload, &decpac)
				shared.CheckErr(err)
				continue
			case shared.GPROMOTE:
				fmt.Println("Handling a GPROMOTE packet.")
				var decpac shared.Guild_promote_cmd_packet
				err = msgpack.Unmarshal(packet.Payload, &decpac)
				shared.CheckErr(err)
				continue
			case shared.SAYCMD:
				fmt.Println("Handling a SAYCMD packet.")
				var decpac shared.Say_cmd_packet
				err = msgpack.Unmarshal(packet.Payload, &decpac)
				shared.CheckErr(err)
				handleSayMessage(decpac)
				continue
			case shared.YELLCMD:
				fmt.Println("Handling a YELLCMD packet.")
				var decpac shared.Yell_cmd_packet
				err = msgpack.Unmarshal(packet.Payload, &decpac)
				shared.CheckErr(err)
				handleYellMessage(decpac)
				continue
			case shared.OOCCMD:
				fmt.Println("Handling a OOCCMD packet.")
				var decpac shared.Ooc_cmd_packet
				err = msgpack.Unmarshal(packet.Payload, &decpac)
				shared.CheckErr(err)
				handleOocMessage(decpac)
				continue
			case shared.HELPCMD:
				fmt.Println("Handling a HELPCMD packet.")
				var decpac shared.Help_cmd_packet
				err = msgpack.Unmarshal(packet.Payload, &decpac)
				shared.CheckErr(err)
				handleHelpMessage(decpac)
				continue
			case shared.PCHATCMD:
				fmt.Println("Handling a PCHATCMD packet.")
				var decpac shared.Pchat_cmd_packet
				err = msgpack.Unmarshal(packet.Payload, &decpac)
				shared.CheckErr(err)
				handlePartyMessage(decpac)
				continue
			case shared.GCHATCMD:
				fmt.Println("Handling a GCHATCMD packet.")
				var decpac shared.Gchat_cmd_packet
				err = msgpack.Unmarshal(packet.Payload, &decpac)
				shared.CheckErr(err)
				handleGuildMessage(decpac)
				continue
			case shared.WHISPERCMD:
				fmt.Println("Handling a WHISPERCMD packet.")
				var decpac shared.Whisper_cmd_packet
				err = msgpack.Unmarshal(packet.Payload, &decpac)
				shared.CheckErr(err)
				handleWhisperMessage(decpac)
				continue
			case shared.ACTIVATECMD:
				fmt.Println("Handling a ACTIVATECMD packet.")
				var decpac shared.Activate_cmd_packet
				err = msgpack.Unmarshal(packet.Payload, &decpac)
				shared.CheckErr(err)
				continue
			case shared.DIALOGCMD:
				fmt.Println("Handling a DIALOGUECMD packet.")
				var decpac shared.Dialogue_cmd_packet
				err = msgpack.Unmarshal(packet.Payload, &decpac)
				shared.CheckErr(err)
				continue
			case shared.BUYITEM:
				fmt.Println("Handling a BUYITEM packet.")
				var decpac shared.Buy_item_packet
				err = msgpack.Unmarshal(packet.Payload, &decpac)
				shared.CheckErr(err)
				continue
			case shared.TAKELOOT:
				fmt.Println("Handling a TAKELOOT packet.")
				var decpac shared.Take_loot_packet
				err = msgpack.Unmarshal(packet.Payload, &decpac)
				shared.CheckErr(err)
				continue
			case shared.OFFERITEM:
				fmt.Println("Handling a OFFERITEM packet.")
				var decpac shared.Offer_item_packet
				err = msgpack.Unmarshal(packet.Payload, &decpac)
				shared.CheckErr(err)
				continue
			case shared.PULLITEM:
				fmt.Println("Handling a PULLITEM packet.")
				var decpac shared.Pull_item_packet
				err = msgpack.Unmarshal(packet.Payload, &decpac)
				shared.CheckErr(err)
				continue
			case shared.ACCTRADE:
				fmt.Println("Handling a ACCTRADE packet.")
				var decpac shared.Accept_trade_packet
				err = msgpack.Unmarshal(packet.Payload, &decpac)
				shared.CheckErr(err)
				continue
			case shared.UNACCTRADE:
				fmt.Println("Handling a UNACCTRADE packet.")
				var decpac shared.Unaccept_trade_packet
				err = msgpack.Unmarshal(packet.Payload, &decpac)
				shared.CheckErr(err)
				continue
			default:
				continue
			}
		}
	}
}
