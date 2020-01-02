// Filename:    game_manager.go
// Author:      Joseph DeVictoria
// Date:        Sept_1_2018
// Purpose:     This file contains all of the tools we need for managing the game. (npcs, players actions etc)

package main

import "Oldentide/shared"

// Session is an integer that represents the session id of the player
type Session int64

var players = make(map[int64]*shared.Pc)
