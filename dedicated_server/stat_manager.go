// Filename:    stat_manager.go
// Author:      Joseph DeVictoria
// Date:        Sept_1_2018
// Purpose:     The stat management tools and character building tools.

package main

import "Oldentide/shared"

func checkDp(su shared.Skillupdate) bool {
	if su.Predp > 100 {
		return true
	} else {
		return false
	}
}

func validNewPlayer(p shared.Pc) bool {
	// Ned to implement this stuff still...
	return true
}
