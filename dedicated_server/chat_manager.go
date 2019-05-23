// Filename:    chat_manager.go
// Author:      Joseph DeVictoria
// Date:        February_27_2018
// Purpose:     This file is where we will handle chat interactions between players.

package main

import (
	"Oldentide/shared"
	//"database/sql"
	//"flag"
	"fmt"
	//_ "github.com/mattn/go-sqlite3"
	//"github.com/vmihailenco/msgpack"
	//"log"
	//"math/rand"
	//"net"
	//"net/http"
	//"net/smtp"
	//"runtime"
	//"strconv"
	//"time"
)

func handleHelpMessage(packet shared.Raw_packet) {
	fmt.Println("Sending a help message to all players.")
	//    for i, pc := range players {
	//       // send back a RELAYHELP packet.
	//       fmt.Println(pc)
	//    }
}

func handleSayMessage(packet shared.Raw_packet) {
	fmt.Println("Sending a say message to nearby players. (10 yards?)")
}

func handleGuildMessage(packet shared.Raw_packet) {
	fmt.Println("Sending a guild message to fellow online guildmates.")
}

func handleOocMessage(packet shared.Raw_packet) {
	fmt.Println("Sending an OOC message to everyone in the current zone.")
}

func handleWhisperMessage(packet shared.Raw_packet) {
	fmt.Println("Sending a whisper message to a target player.")
}

func handlePartyMessage(packet shared.Raw_packet) {
	fmt.Println("Sending a party message to players in the same party.")
}

func handleYellMessage(packet shared.Raw_packet) {
	fmt.Println("Sending a yell message to nearby players. (50 yards?)")
}
