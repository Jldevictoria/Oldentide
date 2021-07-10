// Filename:    oldentide.go (Formerly test_client.go)
// Author:      Joseph DeVictoria
// Date:        February_26_2019
// Purpose:     A command line based test client for Oldentide written in Go.

package main

import (
	"Oldentide/shared"
	"bufio"
	"errors"
	"flag"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/vmihailenco/msgpack"
)

// Global program variables.
var err error
var sadd string
var sport int
var cport int
var test int
var sid uint64
var gui bool
var serverConnection net.Conn
var packetCount int

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
	flag.StringVar(&sadd, "server", "0.0.0.0", "Dedicated game server address.")
	flag.IntVar(&sport, "sport", 1337, "Port used for dedicated game server.")
	flag.IntVar(&cport, "cport", 1338, "Port used for client listener.")
	flag.IntVar(&test, "test", 0, "Test number within the test_client that we want to call. If not given, it will default to a sample \"game\".")
	flag.Uint64Var(&sid, "session", rand.Uint64(), "Session will allow you to force a SessionID for your packets.")
	flag.BoolVar(&gui, "gui", false, "Define whether you want to use the gui option!")
}

func main() {
	// Extract command line input.
	flag.Parse()
	fmt.Println("\n-------------------------------------------------------")
	fmt.Println("Server Configurations from command line:")
	fmt.Println("-------------------------------------------------------")
	fmt.Println("Server Address:", sadd)
	fmt.Println("Server Port:", sport)
	fmt.Println("\n-------------------------------------------------------")
	fmt.Println("Starting Oldentide command line client!")
	fmt.Println("-------------------------------------------------------")
	// Listener.
	clientAddress := net.UDPAddr{
		IP:   net.IP{0, 0, 0, 0},
		Port: cport,
	}
	listenSocket, err := net.ListenUDP("udp", &clientAddress)
	shared.CheckErr(err)
	go collect(listenSocket)

	// Set up server connection through udp socket descriptor struct.
	serverConnection, err = net.Dial("udp", sadd+":"+strconv.Itoa(sport))
	shared.CheckErr(err)
	defer serverConnection.Close()

	if gui {
		fmt.Println("Will add in the future...")
	} else {
		cline := bufio.NewReader(os.Stdin)
		for {
			fmt.Print("$ ")
			command, err := cline.ReadString('\n')
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
			err = runCommand(command)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		}
	}
}

func runCommand(command string) error {
	command = strings.TrimSuffix(command, "\n")
	commandTokens := strings.Fields(command)
	switch commandTokens[0] {
	case "/exit":
		os.Exit(0)
		break
	case "/help":
		fmt.Println("/exit - Quits this application.")
		fmt.Println("/help - Prints this help test.")
		fmt.Println("/spam <count> - Sends <count> number of generic packets at the server.")
		fmt.Println("/newcharacter <account> <firstname> - Attempts to create a new character with specified <firstname> on the server for <account>.")
		fmt.Println("/requestcharacterlist <account> - Gets the characters on the server for the specified <account>.")
		fmt.Println("/s <any message> - Sends <any message> as a \"Say\" command.")
		fmt.Println("/y <any message> - Sends <any message> as a \"Yell\" command.")
		fmt.Println("/ooc <any message> - Sends <any message> as an \"Out of Character\" command.")
		fmt.Println("/h <any message> - Sends <any message> as a \"Help\" command.")
		fmt.Println("/p <any message> - Sends <any message> as a \"Party Chat\" command.")
		fmt.Println("/g <any message> - Sends <any message> as a \"Guild Chat\" command.")
		fmt.Println("/w <target_player> <any message> - Sends <any message> directly to <target_player> as a \"Whisper\" command.")
		fmt.Println("/move <x> <y> <z> <direction> - Moves the specified player character to <x>, <y>, <z>, <direction>.")
		fmt.Println("/connect <target_player> - Attempts to connect <target_player> to the server.")
		break
	case "/spam":
		if len(commandTokens) != 2 {
			return errors.New("wrong arguments to /spam")
		}
		pac := shared.OpcodePacket{Opcode: shared.GENERIC}
		reqpac, err := msgpack.Marshal(pac)
		shared.CheckErr(err)
		count, err := strconv.Atoi(commandTokens[1])
		shared.CheckErr(err)
		for i := 0; i < count; i++ {
			serverConnection.Write(reqpac)
		}
		break
	case "/newcharacter":
		if len(commandTokens) != 3 {
			return errors.New("wrong arguments to /newcharacter")
		}
		p := makePlayer(commandTokens[2])
		pac := shared.CreatePlayerPacket{Opcode: shared.CREATEPLAYER, Pc: p}
		shared.MarshallAndSendPacket(pac, serverConnection)
		break
	case "/requestcharacterlist":
		if len(commandTokens) != 2 {
			return errors.New("wrong arguments to /requestcharacterlist")
		}
		pac := shared.ReqClistPacket{Opcode: shared.REQCLIST, Account: commandTokens[1]}
		shared.MarshallAndSendPacket(pac, serverConnection)
		break
	case "/s":
		pac := shared.SayCmdPacket{Opcode: shared.SAYCMD, SessionID: sid, Message: command[2:]}
		shared.MarshallAndSendPacket(pac, serverConnection)
		break
	case "/y":
		pac := shared.YellCmdPacket{Opcode: shared.YELLCMD, SessionID: sid, Message: command[2:]}
		shared.MarshallAndSendPacket(pac, serverConnection)
		break
	case "/ooc":
		pac := shared.OocCmdPacket{Opcode: shared.OOCCMD, SessionID: sid, Message: command[2:]}
		shared.MarshallAndSendPacket(pac, serverConnection)
		break
	case "/h":
		pac := shared.HelpCmdPacket{Opcode: shared.HELPCMD, SessionID: sid, Message: command[2:]}
		shared.MarshallAndSendPacket(pac, serverConnection)
		break
	case "/p":
		pac := shared.PchatCmdPacket{Opcode: shared.PCHATCMD, SessionID: sid, Message: command[2:]}
		shared.MarshallAndSendPacket(pac, serverConnection)
		break
	case "/g":
		pac := shared.GchatCmdPacket{Opcode: shared.GCHATCMD, SessionID: sid, Message: command[2:]}
		shared.MarshallAndSendPacket(pac, serverConnection)
		break
	case "/w":
		if len(commandTokens) < 3 {
			return errors.New("not enough arguments to /w")
		}
		pac := shared.WhisperCmdPacket{
			Opcode:    shared.WHISPERCMD,
			SessionID: sid,
			Target:    commandTokens[1],
			Message:   strings.Replace(command[2:], " "+commandTokens[1], "", -1),
		}
		shared.MarshallAndSendPacket(pac, serverConnection)
		break
	case "/move":
		if len(commandTokens) != 5 {
			return errors.New("wrong enough arguments to /move")
		}
		pac := shared.MovePlayerPacket{Opcode: shared.MOVEPLAYER, SessionID: sid}
		x, err := strconv.ParseFloat(commandTokens[1], 32)
		shared.CheckErr(err)
		y, err := strconv.ParseFloat(commandTokens[2], 32)
		shared.CheckErr(err)
		z, err := strconv.ParseFloat(commandTokens[3], 32)
		shared.CheckErr(err)
		direction, err := strconv.ParseFloat(commandTokens[4], 32)
		pac.X = float32(x)
		pac.Y = float32(y)
		pac.Z = float32(z)
		pac.Direction = float32(direction)
		shared.MarshallAndSendPacket(pac, serverConnection)
		break
	case "/movespam":
		if len(commandTokens) != 2 {
			return errors.New("wrong enough arguments to /movespam")
		}
		numPackets, err := strconv.Atoi(commandTokens[1])
		shared.CheckErr(err)
		for i := 0; i < numPackets; i++ {
			pac := shared.MovePlayerPacket{Opcode: shared.MOVEPLAYER, SessionID: sid, X: rand.Float32(), Y: rand.Float32(), Z: rand.Float32(), Direction: rand.Float32()}
			shared.MarshallAndSendPacket(pac, serverConnection)
		}
	case "/connect":
		if len(commandTokens) != 2 {
			return errors.New("not enough arguments to /connect")
		}
		pac := shared.ConnectPacket{Opcode: shared.CONNECT, SessionID: sid, Firstname: commandTokens[1]}
		shared.MarshallAndSendPacket(pac, serverConnection)
		break
	case "/disconnect":
		pac := shared.DisconnectPacket{Opcode: shared.DISCONNECT, SessionID: sid}
		shared.MarshallAndSendPacket(pac, serverConnection)
		break
	default:
		return errors.New("target command \"" + commandTokens[0] + "\" is not a valid command.")
	}
	return nil
}

// collect will simply listen for and print any incoming packets.
func collect(connection *net.UDPConn) {
	for {
		buffer := make([]byte, 4096) // Max IPv4 UDP packet size.
		n, remoteAddress, err := connection.ReadFromUDP(buffer)
		shared.IfErrPrintErr(err)
		fmt.Println(shared.RawPacket{Size: n, Client: remoteAddress, Payload: buffer})
		packetCount++
		fmt.Println("PC:", packetCount)
	}
}

// makePlayer creates a sample player for testing the server.
func makePlayer(name string) shared.Pc {
	return shared.Pc{
		ID:             0,
		AccountID:      0,
		Firstname:      name,
		Lastname:       "Mc" + name + "face",
		Guild:          "Gremlins",
		Race:           "Human",
		Gender:         "Male",
		Face:           "White",
		Skin:           "White",
		Profession:     "Engineer",
		Alive:          true,
		Plevel:         51,
		Dp:             12000,
		Hp:             450,
		Maxhp:          450,
		Bp:             250,
		Maxbp:          250,
		Mp:             300,
		Maxmp:          300,
		Sp:             150,
		Maxsp:          150,
		Strength:       65,
		Constitution:   45,
		Intelligence:   50,
		Dexterity:      50,
		Axe:            0,
		Dagger:         0,
		Unarmed:        499,
		Hammer:         0,
		Polearm:        0,
		Spear:          0,
		Staff:          0,
		Sword:          0,
		Archery:        0,
		Crossbow:       0,
		Sling:          0,
		Thrown:         0,
		Armor:          300,
		Dualweapon:     499,
		Shield:         0,
		Bardic:         0,
		Conjuring:      0,
		Druidic:        0,
		Illusion:       0,
		Necromancy:     0,
		Sorcery:        0,
		Shamanic:       0,
		Spellcraft:     0,
		Summoning:      0,
		Focus:          0,
		Armorsmithing:  0,
		Tailoring:      0,
		Fletching:      0,
		Weaponsmithing: 0,
		Alchemy:        0,
		Lapidary:       0,
		Calligraphy:    0,
		Enchanting:     0,
		Herbalism:      0,
		Hunting:        0,
		Mining:         0,
		Bargaining:     0,
		Camping:        0,
		Firstaid:       0,
		Lore:           0,
		Picklocks:      0,
		Scouting:       0,
		Search:         0,
		Stealth:        0,
		Traps:          0,
		Aeolandis:      0,
		Hieroform:      0,
		Highgundis:     0,
		Oldpraxic:      100,
		Praxic:         100,
		Runic:          0,
		Head:           "None",
		Chest:          "None",
		Arms:           "None",
		Hands:          "None",
		Legs:           "None",
		Feet:           "None",
		Cloak:          "None",
		Necklace:       "None",
		Ringone:        "None",
		Ringtwo:        "None",
		Righthand:      "None",
		Lefthand:       "None",
		Zone:           "Iskirrian Plains",
		X:              0,
		Y:              0,
		Z:              0,
		Direction:      47.3,
	}
}
