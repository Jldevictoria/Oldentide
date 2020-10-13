// Filename:    game_manager.go
// Author:      Joseph DeVictoria
// Date:        Sept_1_2018
// Purpose:     This file contains all of the tools we need for managing the game. (npcs, players actions etc)

package main

import (
	"Oldentide/shared"
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

// Npcs is the object that hold all the Npcs pulled from the database.
var Npcs []shared.Npc

// SessionsPlayers is a map that connects session IDs to player objects for maintaining the gamestate.
var SessionsPlayers = make(map[uint64]string)

// SessionConnectionInfo is a map that connects session IDs to connection info for tracking players IP.
var SessionConnectionInfo = make(map[uint64]*net.UDPAddr)

// PlayerNameSessions is a map that connects player names to their session ID, which is useful for sending commands to a given player by name.
var PlayerNameSessions = make(map[string]uint64)
