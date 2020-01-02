// Filename:    game_manager.go
// Author:      Joseph DeVictoria
// Date:        Sept_1_2018
// Purpose:     This file contains all of the tools we need for managing the game. (npcs, players actions etc)

package main

import (
	"Oldentide/shared"
	"fmt"
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
var Pcs []shared.Pc

// Npcs is the object that hold all the Npcs pulled from the database.
var Npcs []shared.Npc

// Sessions is a map that connects session IDs to player objects for maintaining the gamestate.
var Sessions = make(map[int64]*shared.Pc)

// SessionConnectionInfo is a map that connects session IDs to connection info for tracking players IP.
var SessionConnectionInfo = make(map[int64]*net.UDPAddr)

// PlayerNameSessions is a map that connects player names to their session ID, which is useful for sending commands to a given player by name.
var PlayerNameSessions = make(map[string]int64)

// GetPlayerByFirstname will look through all the recorded players and return the Pc object found under the given name.
func GetPlayerByFirstname(firstname string) (shared.Pc, error) {
	for _, tp := range Pcs {
		if tp.Firstname == firstname {
			return tp, nil
		}
	}
	return *new(shared.Pc), fmt.Errorf("player %s does not exist", firstname)
}
