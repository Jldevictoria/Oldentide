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

func handleSayMessage(packet shared.Say_packet) {
	fmt.Print("[/s](", packet.Session_id, "): ", packet.Text)
}

func handleYellMessage(packet shared.Yell_packet) {
	fmt.Print("[/y](", packet.Session_id, "): ", packet.Text)
}

func handleOocMessage(packet shared.Ooc_packet) {
	fmt.Print("[/ooc](", packet.Session_id, "): ", packet.Text)
}

func handleHelpMessage(packet shared.Help_packet) {
	fmt.Print("[/h](", packet.Session_id, "): ", packet.Text)
}

func handlePartyMessage(packet shared.Pchat_packet) {
	fmt.Print("[/p](", packet.Session_id, "): ", packet.Text)
}

func handleGuildMessage(packet shared.Gchat_packet) {
	fmt.Print("[/g](", packet.Session_id, "): ", packet.Text)
}

func handleWhisperMessage(packet shared.Whisper_packet) {
	fmt.Print("[/w target](", packet.Session_id, " -> ", packet.Target, "): ", packet.Text)
}
