// Filename:    test_client.go
// Author:      Joseph DeVictoria
// Date:        February_26_2019
// Purpose:     A command line based test client for Oldentide written in Go.

package main

import (
	"Oldentide/shared"
	"flag"
	"fmt"
	"github.com/vmihailenco/msgpack"
	"math/rand"
	"net"
	"strconv"
	"time"
)

// Global program variables.
var err error
var sadd string
var sport int
var test int

func init() {
	flag.StringVar(&sadd, "server", "0.0.0.0", "Dedicated game server address.")
	flag.IntVar(&sport, "port", 1337, "Port used for dedicated game server.")
	flag.IntVar(&test, "test", 0, "Test number within the test_client that we want to call.")
	rand.Seed(time.Now().UTC().UnixNano())
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
	fmt.Println("Starting Oldentide test client!")
	fmt.Println("-------------------------------------------------------")
	// // Listener.
	// client_address := net.UDPAddr{
	// 	IP:   net.IP{0, 0, 0, 0},
	// 	Port: sport,
	// }
	// listen_socket, err := net.ListenUDP("udp", &client_address)
	// defer listen_socket.Close()
	// shared.CheckErr(err)

	// Set up server connection.
	// Create udp socket description struct.
	server_connection, err := net.Dial("udp", sadd+":"+strconv.Itoa(sport))
	shared.CheckErr(err)
	defer server_connection.Close()

	fmt.Println("Executing test option ", test)

	switch test {
	case 0: // default case.
		fmt.Println("You probably meant to give me a test number.  (-test=[number])")
		break
	case 1: // SPAM
		pac := shared.Opcode_packet{Opcode: shared.GENERIC}
		reqpac, err := msgpack.Marshal(pac)
		shared.CheckErr(err)
		for i := 0; i < 10000000; i++ {
			server_connection.Write(reqpac)
			fmt.Println(i)
		}
		break
	case 2: // create a character.
		p := make_player("Joe")
		pac := shared.Create_player_packet{Opcode: shared.CREATEPLAYER, Pc: p}
		reqpac, err := msgpack.Marshal(pac)
		shared.CheckErr(err)
		server_connection.Write(reqpac)
		break
	case 3: // Request haracter list.
		pac := shared.Req_clist_packet{Opcode: shared.REQCLIST, Account: "test"}
		reqpac, err := msgpack.Marshal(pac)
		shared.CheckErr(err)
		server_connection.Write(reqpac)
		break
	default:
		fmt.Println("You need to give a valid test number.  (-test=[number])")
	}
}

// Simple function to check the error status of an operation.
func make_player(name string) shared.Pc {
	return shared.Pc{
		Id:             0,
		Accountid:      0,
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
		Ep:             150,
		Maxep:          150,
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
