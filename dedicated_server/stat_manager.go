// Filename:    stat_manager.go
// Author:      Joseph DeVictoria
// Date:        Sept_1_2018
// Purpose:     The stat management tools and character building tools.

package main

import (
	"Oldentide/shared"
	"log"
)

// calculateAwardedDP will return the amount of DP to add to a player based on the level they have reached.
// I decided to make this recursive since that is what matched the original model.
func calculateAwardedDP(level int32) int32 {
	if level <= 1 {
		return 1500
	} else if level > 51 {
		return 0
	}
	return (calculateAwardedDP(level-1) + (75 * (level - 1)) + (2 * (level - 1) * (level - 2)))
}

// checkDP will see if the SkillUpdate that was requested by the player is legal according to the game rules.
// If this check passes, we will update the player in the game model and database.
func checkDP(su shared.SkillUpdate) bool {
	_, err := getPlayerByFirstname(su.Playername)
	if err != nil {
		log.Println("Someone attempted to update the skills of a player who was not currently playing.")
	}

	return false
}

func validNewPlayer(p shared.Pc) bool {
	// Need to implement this stuff still...
	return true
}
