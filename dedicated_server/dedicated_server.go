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
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/smtp"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/vmihailenco/msgpack"
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
var debug bool
var eauth smtp.Auth
var db *sql.DB
var packetCount int

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
	flag.IntVar(&gport, "gport", 0, "Port used for dedicated game server.")
	flag.IntVar(&wport, "wport", 0, "Port used for website.")
	flag.BoolVar(&everify, "everify", false, "Use an emailer to verify accounts?")
	flag.StringVar(&webadd, "webadd", "", "Public website root address where accounts will be created.")
	flag.StringVar(&email, "email", "", "Gmail email address used to send verification emails.")
	flag.StringVar(&epass, "epass", "", "Gmail email password used to send verification emails.")
	flag.StringVar(&dbpath, "dbpath", shared.DefaultGOPATH()+"/src/Oldentide/db/oldentide.db", "Path to oldentide.db")
	flag.BoolVar(&debug, "debug", false, "Turn on debugging prints")
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
	fmt.Println("debug:", debug)
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
			log.Fatal("Please provide a Gmail email password with the command line flag -epass=<P@$$word>")
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
	shared.IfErrPrintErr(err)
	fmt.Println("* Database connected.")

	// Initialize the game state (populates all of the npcs, and game objects, etc).
	// --------------------------------------------------------------------------------------------
	RaceTemplates = pullRaceTemplates()
	fmt.Println("\n* Race templates populated from database.")
	if debug {
		shared.PrettyPrint(RaceTemplates)
	}

	ProfessionTemplates = pullProfessionTemplates()
	fmt.Println("\n* Profession templates populated from database.")
	if debug {
		shared.PrettyPrint(ProfessionTemplates)
	}

	ItemTemplates = pullItemTemplates()
	fmt.Println("\n* Item templates populated from database.")
	if debug {
		shared.PrettyPrint(ItemTemplates)
	}

	SpellTemplates = pullSpellTemplates()
	fmt.Println("\n* Spell templates populated from database.")
	if debug {
		shared.PrettyPrint(SpellTemplates)
	}

	Npcs = pullNpcs()
	fmt.Println("\n* NPCs populated from database.")
	if debug {
		shared.PrettyPrint(Npcs)
	}

	// --------------------------------------------------------------------------------------------
	// Kick off http server for registration page.
	// --------------------------------------------------------------------------------------------
	mux := http.NewServeMux()
	mux.HandleFunc("/", routeWebTraffic)
	go http.ListenAndServe(":"+strconv.Itoa(wport), mux)

	// --------------------------------------------------------------------------------------------
	// Create and bind a udp socket descriptor.
	// --------------------------------------------------------------------------------------------
	serverAddress := net.UDPAddr{
		IP:   net.IP{0, 0, 0, 0},
		Port: gport,
	}
	socket, err := net.ListenUDP("udp", &serverAddress)
	shared.IfErrPrintErr(err)

	// --------------------------------------------------------------------------------------------
	// Start our collecter to pull in packets from the hardware socket.
	// --------------------------------------------------------------------------------------------
	rawPacketQueue := make(chan shared.RawPacket)
	quitChan := make(chan bool)
	go Collect(socket, rawPacketQueue, quitChan)
	fmt.Println("\n* Collector Launched.")

	// --------------------------------------------------------------------------------------------
	// Start as many handler goroutines as we have processors.
	// Should result in non-thrashing handler thread concurrency.
	// --------------------------------------------------------------------------------------------
	for i := 0; i < runtime.NumCPU(); i++ {
		go Handle(rawPacketQueue, quitChan, i)
		fmt.Println("* Handler ", i+1, " of ", runtime.NumCPU(), " Launched.")
	}

	// Wait for a message to the Quit Channel to halt execution.
	<-quitChan

	// Close database.
	db.Close()
}

// Collect places all UDP packets that arrive on the hardware socket into a queue for handling.
func Collect(connection *net.UDPConn, rawPacketQueue chan shared.RawPacket, quitChan chan bool) {
	for {
		buffer := make([]byte, 4096) // Max IPv4 UDP packet size.
		n, remoteAddress, err := connection.ReadFromUDP(buffer)
		shared.IfErrPrintErr(err)
		rawPacketQueue <- shared.RawPacket{Size: n, Client: remoteAddress, Payload: buffer}
		packetCount++
		fmt.Println("PC:", packetCount)
	}
}

// Handle all arriving packets based on which opcode they are.
func Handle(rawPacketQueue chan shared.RawPacket, quitChan chan bool, rid int) {
	var packet shared.RawPacket
	for {
		select {
		// This case will run when there is a packet available at the front of the packet queue.
		case packet = <-rawPacketQueue:
			//fmt.Println("Goroutine ID:", rid, "Size:", packet.Size, "Sender:", packet.Client, "Payload:", packet.Payload[:packet.Size]) //debug
			var op shared.OpcodePacket
			err = msgpack.Unmarshal(packet.Payload, &op)
			shared.IfErrPrintErr(err)
			// Depending on what packet opcode we recieved, handle the data accordingly.
			switch op.Opcode {
			case shared.EMPTY:
				fmt.Println("Handling an EMPTY packet.")
				continue
			case shared.GENERIC:
				fmt.Println("Handling a GENERIC packet.")
				continue
			case shared.ACK:
				fmt.Println("Handling an ACK packet.")
				var decpac shared.AckPacket
				err = msgpack.Unmarshal(packet.Payload, &decpac)
				shared.IfErrPrintErr(err)
				err = VerifySession(decpac.SessionID, packet.Client)
				if err != nil {
					fmt.Println(err)
					continue
				}
				continue
			case shared.ERROR:
				fmt.Println("Handling an ERROR packet.")
				var decpac shared.ErrorPacket
				err = msgpack.Unmarshal(packet.Payload, &decpac)
				if err != nil {
					fmt.Println(err)
					continue
				}
				err = VerifySession(decpac.SessionID, packet.Client)
				if err != nil {
					fmt.Println(err)
					continue
				}
				continue
			case shared.REQCLIST:
				fmt.Println("Handling a REQCLIST packet.")
				var decpac shared.ReqClistPacket
				err = msgpack.Unmarshal(packet.Payload, &decpac)
				fmt.Println(decpac)
				err = VerifySession(decpac.SessionID, packet.Client)
				if err != nil {
					fmt.Println(err)
					continue
				}
				var retpac shared.SendClistPacket
				retpac.Opcode = shared.SENDCLIST
				retpac.Characters = getCharacterList(decpac.Account)
				fmt.Println(retpac)
				continue
			case shared.CREATEPLAYER:
				fmt.Println("Handling a CREATEPLAYER packet.")
				var decpac shared.CreatePlayerPacket
				err = msgpack.Unmarshal(packet.Payload, &decpac)
				fmt.Println(decpac)
				err = VerifySession(decpac.SessionID, packet.Client)
				if err != nil {
					fmt.Println(err)
					continue
				}
				// Need to get the account name by session id.
				accountName := "test"
				playerName := decpac.Pc.Firstname
				decpac.Pc.AccountID = getAccountIDFromAccountName(accountName)
				fmt.Println(decpac.Pc.AccountID)
				if getRemainingPlayerSlots(accountName, 10) == 0 {
					log.Println("Account tried to create too many players.")
					continue
				}
				if playerFirstNameTaken(playerName) {
					log.Println("Account tried to create a player whose name was already taken.")
					continue
				}
				if validNewPlayer(decpac.Pc) {
					addNewPlayer(decpac.Pc)
					log.Println("Account \"", accountName, "\" created a new player \"", playerName, "\".")
				} else {
					log.Println("Account is trying something fraudulent during account creation!")
					//banAccount()
					continue
				}
			case shared.CONNECT:
				fmt.Println("Handling a CONNECT packet.")
				var decpac shared.ConnectPacket
				err = msgpack.Unmarshal(packet.Payload, &decpac)
				shared.IfErrPrintErr(err)
				err = VerifySession(decpac.SessionID, packet.Client)
				if err != nil {
					fmt.Println(err)
					continue
				}
				if p, ok := SessionsPlayers[decpac.SessionID]; ok {
					fmt.Println("Player" + p + "tried to connect twice.  Forcing full disconnect.")
					DisconnectSession(decpac.SessionID, true)
				} else if shared.CountUintStringMapInstances(SessionsPlayers, decpac.Firstname) > 0 {
					fmt.Println("Session" + strconv.FormatUint(decpac.SessionID, 10) + "tried to connect to an already existing player.  Hacking!")
				} else {
					var pc, err = getPlayerByFirstname(decpac.Firstname)
					fmt.Println(pc)
					if err != nil {
						fmt.Println(err)
					} else {
						SessionsPlayers[decpac.SessionID] = pc.Firstname
						Pcs[pc.Firstname] = &pc
						fmt.Println(SessionsPlayers)
					}
				}
				continue
			case shared.DISCONNECT:
				fmt.Println("Handling a DISCONNECT packet.")
				var decpac shared.DisconnectPacket
				err = msgpack.Unmarshal(packet.Payload, &decpac)
				shared.IfErrPrintErr(err)
				err = VerifySession(decpac.SessionID, packet.Client)
				if err != nil {
					fmt.Println(err)
					continue
				}
				if _, ok := SessionsPlayers[decpac.SessionID]; ok {
					DisconnectSession(decpac.SessionID, true)
				} else {
					fmt.Println("Player with session", decpac.SessionID, "tried to disconnect, but was never connected...")
				}
				continue
			case shared.MOVEPLAYER:
				fmt.Println("Handling a MOVEPLAYER packet.")
				var decpac shared.MovePlayerPacket
				err = msgpack.Unmarshal(packet.Payload, &decpac)
				shared.IfErrPrintErr(err)
				err = VerifySession(decpac.SessionID, packet.Client)
				if err != nil {
					fmt.Println(err)
					continue
				}
				if player, ok := Pcs[SessionsPlayers[decpac.SessionID]]; ok {
					player.X = decpac.X
					player.Y = decpac.Y
					player.Z = decpac.Z
					player.Direction = decpac.Direction
				} else {
					fmt.Println("Player did not exist in MOVEPLAYER case for session", decpac.SessionID)
				}
				continue
			case shared.SPENDDP:
				fmt.Println("Handling a SPENDDP packet.")
				var decpac shared.SpendDPPacket
				err = msgpack.Unmarshal(packet.Payload, &decpac)
				shared.IfErrPrintErr(err)
				err = VerifySession(decpac.SessionID, packet.Client)
				if err != nil {
					fmt.Println(err)
					continue
				}
				continue
			case shared.TALKCMD:
				fmt.Println("Handling a TALKCMD packet.")
				var decpac shared.TalkCmdPacket
				err = msgpack.Unmarshal(packet.Payload, &decpac)
				shared.IfErrPrintErr(err)
				err = VerifySession(decpac.SessionID, packet.Client)
				if err != nil {
					fmt.Println(err)
					continue
				}
				continue
			case shared.ATTACKCMD:
				fmt.Println("Handling a ATTACKCMD packet.")
				var decpac shared.AttackCmdPacket
				err = msgpack.Unmarshal(packet.Payload, &decpac)
				shared.IfErrPrintErr(err)
				err = VerifySession(decpac.SessionID, packet.Client)
				if err != nil {
					fmt.Println(err)
					continue
				}
				continue
			case shared.TRADECMD:
				fmt.Println("Handling a TRADECMD packet.")
				var decpac shared.TradeCmdPacket
				err = msgpack.Unmarshal(packet.Payload, &decpac)
				shared.IfErrPrintErr(err)
				err = VerifySession(decpac.SessionID, packet.Client)
				if err != nil {
					fmt.Println(err)
					continue
				}
				continue
			case shared.INVITECMD:
				fmt.Println("Handling a INVITECMD packet.")
				var decpac shared.InviteCmdPacket
				err = msgpack.Unmarshal(packet.Payload, &decpac)
				shared.IfErrPrintErr(err)
				err = VerifySession(decpac.SessionID, packet.Client)
				if err != nil {
					fmt.Println(err)
					continue
				}
				continue
			case shared.GINVITECMD:
				fmt.Println("Handling a GINVITECMD packet.")
				var decpac shared.GuildInviteCmdPacket
				err = msgpack.Unmarshal(packet.Payload, &decpac)
				shared.IfErrPrintErr(err)
				err = VerifySession(decpac.SessionID, packet.Client)
				if err != nil {
					fmt.Println(err)
					continue
				}
				continue
			case shared.GKICK:
				fmt.Println("Handling a GKICK packet.")
				var decpac shared.GuildKickCmdPacket
				err = msgpack.Unmarshal(packet.Payload, &decpac)
				shared.IfErrPrintErr(err)
				err = VerifySession(decpac.SessionID, packet.Client)
				if err != nil {
					fmt.Println(err)
					continue
				}
				continue
			case shared.GPROMOTE:
				fmt.Println("Handling a GPROMOTE packet.")
				var decpac shared.GuildPromoteCmdPacket
				err = msgpack.Unmarshal(packet.Payload, &decpac)
				shared.IfErrPrintErr(err)
				err = VerifySession(decpac.SessionID, packet.Client)
				if err != nil {
					fmt.Println(err)
					continue
				}
				continue
			case shared.GDEMOTE:
				fmt.Println("Handling a GDEMOTE packet.")
				var decpac shared.GuildPromoteCmdPacket
				err = msgpack.Unmarshal(packet.Payload, &decpac)
				shared.IfErrPrintErr(err)
				err = VerifySession(decpac.SessionID, packet.Client)
				if err != nil {
					fmt.Println(err)
					continue
				}
				continue
			case shared.SAYCMD:
				fmt.Println("Handling a SAYCMD packet.")
				var decpac shared.SayCmdPacket
				err = msgpack.Unmarshal(packet.Payload, &decpac)
				shared.IfErrPrintErr(err)
				err = VerifySession(decpac.SessionID, packet.Client)
				if err != nil {
					fmt.Println(err)
					continue
				}
				handleSayMessage(decpac)
				continue
			case shared.YELLCMD:
				fmt.Println("Handling a YELLCMD packet.")
				var decpac shared.YellCmdPacket
				err = msgpack.Unmarshal(packet.Payload, &decpac)
				shared.IfErrPrintErr(err)
				err = VerifySession(decpac.SessionID, packet.Client)
				if err != nil {
					fmt.Println(err)
					continue
				}
				handleYellMessage(decpac)
				continue
			case shared.OOCCMD:
				fmt.Println("Handling a OOCCMD packet.")
				var decpac shared.OocCmdPacket
				err = msgpack.Unmarshal(packet.Payload, &decpac)
				shared.IfErrPrintErr(err)
				err = VerifySession(decpac.SessionID, packet.Client)
				if err != nil {
					fmt.Println(err)
					continue
				}
				handleOocMessage(decpac)
				continue
			case shared.HELPCMD:
				fmt.Println("Handling a HELPCMD packet.")
				var decpac shared.HelpCmdPacket
				err = msgpack.Unmarshal(packet.Payload, &decpac)
				shared.IfErrPrintErr(err)
				err = VerifySession(decpac.SessionID, packet.Client)
				if err != nil {
					fmt.Println(err)
					continue
				}
				handleHelpMessage(decpac)
				continue
			case shared.PCHATCMD:
				fmt.Println("Handling a PCHATCMD packet.")
				var decpac shared.PchatCmdPacket
				err = msgpack.Unmarshal(packet.Payload, &decpac)
				shared.IfErrPrintErr(err)
				err = VerifySession(decpac.SessionID, packet.Client)
				if err != nil {
					fmt.Println(err)
					continue
				}
				handlePartyMessage(decpac)
				continue
			case shared.GCHATCMD:
				fmt.Println("Handling a GCHATCMD packet.")
				var decpac shared.GchatCmdPacket
				err = msgpack.Unmarshal(packet.Payload, &decpac)
				shared.IfErrPrintErr(err)
				err = VerifySession(decpac.SessionID, packet.Client)
				if err != nil {
					fmt.Println(err)
					continue
				}
				handleGuildMessage(decpac)
				continue
			case shared.WHISPERCMD:
				fmt.Println("Handling a WHISPERCMD packet.")
				var decpac shared.WhisperCmdPacket
				err = msgpack.Unmarshal(packet.Payload, &decpac)
				shared.IfErrPrintErr(err)
				err = VerifySession(decpac.SessionID, packet.Client)
				if err != nil {
					fmt.Println(err)
					continue
				}
				handleWhisperMessage(decpac)
				continue
			case shared.ACTIVATECMD:
				fmt.Println("Handling a ACTIVATECMD packet.")
				var decpac shared.ActivateCmdPacket
				err = msgpack.Unmarshal(packet.Payload, &decpac)
				shared.IfErrPrintErr(err)
				err = VerifySession(decpac.SessionID, packet.Client)
				if err != nil {
					fmt.Println(err)
					continue
				}
				continue
			case shared.DIALOGCMD:
				fmt.Println("Handling a DIALOGUECMD packet.")
				var decpac shared.DialogueCmdPacket
				err = msgpack.Unmarshal(packet.Payload, &decpac)
				shared.IfErrPrintErr(err)
				err = VerifySession(decpac.SessionID, packet.Client)
				if err != nil {
					fmt.Println(err)
					continue
				}
				continue
			case shared.BUYITEM:
				fmt.Println("Handling a BUYITEM packet.")
				var decpac shared.BuyItemPacket
				err = msgpack.Unmarshal(packet.Payload, &decpac)
				shared.IfErrPrintErr(err)
				err = VerifySession(decpac.SessionID, packet.Client)
				if err != nil {
					fmt.Println(err)
					continue
				}
				continue
			case shared.TAKELOOT:
				fmt.Println("Handling a TAKELOOT packet.")
				var decpac shared.TakeLootPacket
				err = msgpack.Unmarshal(packet.Payload, &decpac)
				shared.IfErrPrintErr(err)
				err = VerifySession(decpac.SessionID, packet.Client)
				if err != nil {
					fmt.Println(err)
					continue
				}
				continue
			case shared.OFFERITEM:
				fmt.Println("Handling a OFFERITEM packet.")
				var decpac shared.OfferItemPacket
				err = msgpack.Unmarshal(packet.Payload, &decpac)
				shared.IfErrPrintErr(err)
				err = VerifySession(decpac.SessionID, packet.Client)
				if err != nil {
					fmt.Println(err)
					continue
				}
				continue
			case shared.PULLITEM:
				fmt.Println("Handling a PULLITEM packet.")
				var decpac shared.PullItemPacket
				err = msgpack.Unmarshal(packet.Payload, &decpac)
				shared.IfErrPrintErr(err)
				err = VerifySession(decpac.SessionID, packet.Client)
				if err != nil {
					fmt.Println(err)
					continue
				}
				continue
			case shared.ACCTRADE:
				fmt.Println("Handling a ACCTRADE packet.")
				var decpac shared.AcceptTradePacket
				err = msgpack.Unmarshal(packet.Payload, &decpac)
				shared.IfErrPrintErr(err)
				err = VerifySession(decpac.SessionID, packet.Client)
				if err != nil {
					fmt.Println(err)
					continue
				}
				continue
			case shared.UNACCTRADE:
				fmt.Println("Handling a UNACCTRADE packet.")
				var decpac shared.UnacceptTradePacket
				err = msgpack.Unmarshal(packet.Payload, &decpac)
				shared.IfErrPrintErr(err)
				err = VerifySession(decpac.SessionID, packet.Client)
				if err != nil {
					fmt.Println(err)
					continue
				}
				continue
			default:
				continue
			}
		}
	}
}
