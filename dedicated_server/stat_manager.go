// Filename:    stat_manager.go
// Author:      Joseph DeVictoria
// Date:        Sept_1_2018
// Purpose:     The stat management tools and character building tools.

package main

import "Oldentide/shared"

// DPPerLevel is a list of how much experience is granted at each level (starting from level 1).
// I would like to do this via formula, but I'm having trouble fitting a common equation to this line.
var DPPerLevel []int64 = []int64{
	1500,
	1575,
	1729,
	1965,
	2289,
	2704,
	3214,
	3823,
	4534,
	5353,
	6283,
	7328,
	8492,
	9778,
	11192,
	12737,
	14417,
	16236,
	18197,
	20306,
	22566,
	24981,
	27555,
	30291,
	33195,
	36270,
	39520,
	42949,
	46560,
	50359,
	54349,
	58534,
	62918,
	67504,
	72298,
	77303,
	82523,
	87962,
	93623,
	99512,
	105632,
	111987,
	118581,
	125417,
	132501,
	139836,
	147426,
	155275,
	163386,
	171765,
	178685,
}

func checkDp(su shared.SkillUpdate) bool {
	if su.Predp > 100 {
		return true
	}
	return false
}

func validNewPlayer(p shared.Pc) bool {
	// Need to implement this stuff still...
	return true
}
