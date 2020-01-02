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

func handleSayMessage(packet shared.SayCmdPacket) {
	fmt.Println("[/s](", packet.SessionID, "): ", packet.Text)
}

func handleYellMessage(packet shared.YellCmdPacket) {
	fmt.Println("[/y](", packet.SessionID, "): ", packet.Text)
}

func handleOocMessage(packet shared.OocCmdPacket) {
	fmt.Println("[/ooc](", packet.SessionID, "): ", packet.Text)
}

func handleHelpMessage(packet shared.HelpCmdPacket) {
	fmt.Println("[/h](", packet.SessionID, "): ", packet.Text)
}

func handlePartyMessage(packet shared.PchatCmdPacket) {
	fmt.Println("[/p](", packet.SessionID, "): ", packet.Text)
}

func handleGuildMessage(packet shared.GchatCmdPacket) {
	fmt.Println("[/g](", packet.SessionID, "): ", packet.Text)
}

func handleWhisperMessage(packet shared.WhisperCmdPacket) {
	fmt.Println("[/w target](", packet.SessionID, " -> ", packet.Target, "): ", packet.Text)

}
