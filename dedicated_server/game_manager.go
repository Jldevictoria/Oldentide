// Filename:    game_manager.go
// Author:      Joseph DeVictoria
// Date:        Sept_1_2018
// Purpose:     This file contains all of the tools we need for managing the game. (npcs, players actions etc)

package main

import (
	"Oldentide/shared"
	"errors"
	"net"
)

// RaceTemplates is the object that hold all the race templates pulled from the database.
var RaceTemplates []shared.RaceTemplate

// ProfessionTemplates is the object that hold all the profession templates pulled from the database.
var ProfessionTemplates []shared.ProfessionTemplate

// ItemTemplates is the object that hold all the item templates pulled from the database.
var ItemTemplates []shared.ItemTemplate

// SpellTemplates is the object that hold all the spell templates pulled from the database.
var SpellTemplates []shared.SpellTemplate

// Pcs is the object that hold all the Pcs pulled from the database.
var Pcs = make(map[string]*shared.Pc)

// SessionsPlayers is a map that connects session IDs to player objects for maintaining the gamestate.
var SessionsPlayers = make(map[uint64]string)

// PlayerSessions is a map that connects player names to their session ID, which is useful for sending commands to a given player by name.
var PlayerSessions = make(map[string]uint64)

// SessionAddressInfo is a map that connects session IDs to connection info for tracking players IP.
var SessionAddressInfo = make(map[uint64]*net.UDPAddr)

// SessionConnectionInfo is a map that connects sessions to live connections.
var SessionConnectionInfo = make(map[uint64]net.Conn)

// Npcs is the object that hold all the Npcs pulled from the database.
var Npcs []shared.Npc

// VerifySession will check if a session already has an associated client address, assign it if not, and throw error if session doesn't match client.
func VerifySession(sessionID uint64, client *net.UDPAddr) error {
	client.Port = 1338 // Hard coded for now...
	if previousClient, ok := SessionAddressInfo[sessionID]; !ok {
		SessionAddressInfo[sessionID] = client
		clientConnection, err := net.Dial("udp", client.String())
		shared.CheckErr(err)
		if _, ok := SessionConnectionInfo[sessionID]; ok {
			SessionConnectionInfo[sessionID].Close() // Close any previous connection that was here?
			return errors.New("Error, we didnt have this address, but we did have a connection during session verification")
		}
		SessionConnectionInfo[sessionID] = clientConnection
		defer SessionConnectionInfo[sessionID].Close()
	} else if previousClient.String() != client.String() {
		return errors.New("Error, different computer tried to hijack player's sessionID")
	}
	return nil
}

// DisconnectSession will completely remove the target session ID from the game (with a disconnect packet!)
func DisconnectSession(sessionID uint64, sendPac bool) {
	delete(PlayerSessions, SessionsPlayers[sessionID])
	delete(SessionsPlayers, sessionID)
	if sendPac {
		pac := shared.DisconnectPacket{Opcode: shared.DISCONNECT, SessionID: sessionID}
		shared.MarshallAndSendPacket(pac, SessionConnectionInfo[sessionID])
	}
	delete(SessionAddressInfo, sessionID)
	delete(SessionConnectionInfo, sessionID)
}

// GetPlayersInZone will return all of the players in a given zone.
//func GetPlayersInZone()
