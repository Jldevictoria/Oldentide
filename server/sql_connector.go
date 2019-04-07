// Filename:    sql_connector.go
// Author:      Joseph DeVictoria
// Date:        June_16_2018
// Purpose:     Database / sql abstraction for Oldentide dedicated server.

package main

import "Oldentide/shared"
import "database/sql"
import "fmt"
import "log"

// Check if the account name is already taken in the database.
func accountExists(a string) bool {
	rows, err := db.Query("SELECT accountname FROM accounts WHERE accountname='" + a + "'")
	return foundInRows(rows, err)
}

// Check if this email already has an account associated with it in the database.
func emailExists(e string) bool {
	rows, err := db.Query("SELECT accountname FROM accounts WHERE email='" + e + "'")
	return foundInRows(rows, err)
}

func createAccount(accountname string, email string, verify_key string, hashed_key string, salt_key string) bool {
	// Prepare insert statement.
	ins, err := db.Prepare("INSERT INTO accounts(valid, banned, accountname, email, gamesession, playing, verify, hash, salt) values(?, ?, ?, ?, ?, ?, ?, ?, ?)")
	shared.CheckErr(err)
	// Try to populate and execute an SQL statment.
	_, err = ins.Exec("0", "0", accountname, email, "0", "0", verify_key, hashed_key, salt_key)
	if err == nil {
		return true
	} else {
		return false
	}
}

func getAccountnameFromVerifyKey(k string) string {
	rows, _ := db.Query("SELECT accountname FROM accounts WHERE verify='" + k + "'")
	//ifErrPrintErr(err)
	var accountname string = ""
	if rows != nil {
		for rows.Next() {
			rows.Scan(&accountname)
		}
	}
	rows.Close()
	return accountname
}

func getSaltFromAccount(account string) string {
	rows, err := db.Query("SELECT salt FROM accounts WHERE accountname=?", account)
	shared.IfErrPrintErr(err)
	var salt string = ""
	if rows != nil {
		for rows.Next() {
			rows.Scan(&salt)
		}
	}
	rows.Close()
	return salt
}

func getHashFromAccount(account string) string {
	rows, err := db.Query("SELECT hash FROM accounts WHERE accountname=?", account)
	shared.IfErrPrintErr(err)
	var hash string = ""
	if rows != nil {
		for rows.Next() {
			rows.Scan(&hash)
		}
	}
	rows.Close()
	return hash
}

func setSessionId(account string, session int) bool {
	update, err := db.Prepare("UPDATE accounts SET gamesession=? WHERE accountname=?")
	_, err = update.Exec(session, account)
	if err == nil {
		return true
	} else {
		return false
	}
}

func activateAccount(a string) bool {
	update, err := db.Prepare("UPDATE accounts SET valid=? WHERE accountname=?")
	_, err = update.Exec("1", a)
	if err == nil {
		return true
	} else {
		return false
	}
}

func banAccount(a string) bool {
	ban, err := db.Prepare("UPDATE accounts SET banned=? WHERE accountname=?")
	_, err = ban.Exec("1", a)
	if err == nil {
		return true
	} else {
		return false
	}
}

// Queries the db to make sure that we are generating a truly unique verify key.
func generateUniqueVerify(n int) string {
	findKey := true
	for findKey {
		verify_key := shared.GenerateRandomLetters(n)
		rows, err := db.Query("SELECT accountname FROM accounts WHERE verify='" + verify_key + "'")
		//ifErrPrintErr(err)
		if !foundInRows(rows, err) {
			return verify_key
		}
	}
	return ""
}

// Queries the db to make sure that we are generating a truly unique salt.
func generateUniqueSalt(n int) string {
	findKey := true
	for findKey {
		salt_key := shared.GenerateRandomLetters(n)
		rows, err := db.Query("SELECT accountname FROM accounts WHERE salt='" + salt_key + "'")
		if !foundInRows(rows, err) {
			return salt_key
		}
	}
	return ""
}

func foundInRows(rows *sql.Rows, err error) bool {
	found := false
	if err == sql.ErrNoRows {
		found = false
	} else if err != nil {
		log.Fatal(err)
	} else if rows.Next() {
		found = true
	}
	rows.Close()
	return found
}

func pullPcs() []shared.Pc {
	rows, err := db.Query("Select * FROM players")
	defer rows.Close()
	shared.CheckErr(err)
	return pcRowsToStruct(rows)
}

func getPlayers(account string) []shared.Pc {
	rows, err := db.Query("Select * FROM players WHERE account=?", account)
	defer rows.Close()
	shared.CheckErr(err)
	return pcRowsToStruct(rows)
}

func pcRowsToStruct(rows * sql.Rows) []shared.Pc {
	var pcs []shared.Pc
	for rows.Next() {
		var pc shared.Pc
		err = rows.Scan(
			&pc.Id,
			&pc.Accountid,
			&pc.Firstname,
			&pc.Lastname,
			&pc.Guild,
			&pc.Race,
			&pc.Gender,
			&pc.Face,
			&pc.Skin,
			&pc.Profession,
			&pc.Alive,
			&pc.Plevel,
			&pc.Dp,
			&pc.Hp,
			&pc.Maxhp,
			&pc.Bp,
			&pc.Maxbp,
			&pc.Mp,
			&pc.Maxmp,
			&pc.Ep,
			&pc.Maxep,
			&pc.Strength,
			&pc.Constitution,
			&pc.Intelligence,
			&pc.Dexterity,
			&pc.Axe,
			&pc.Dagger,
			&pc.Unarmed,
			&pc.Hammer,
			&pc.Polearm,
			&pc.Spear,
			&pc.Staff,
			&pc.Sword,
			&pc.Archery,
			&pc.Crossbow,
			&pc.Sling,
			&pc.Thrown,
			&pc.Armor,
			&pc.Dualweapon,
			&pc.Shield,
			&pc.Bardic,
			&pc.Conjuring,
			&pc.Druidic,
			&pc.Illusion,
			&pc.Necromancy,
			&pc.Sorcery,
			&pc.Shamanic,
			&pc.Spellcraft,
			&pc.Summoning,
			&pc.Focus,
			&pc.Armorsmithing,
			&pc.Tailoring,
			&pc.Fletching,
			&pc.Weaponsmithing,
			&pc.Alchemy,
			&pc.Lapidary,
			&pc.Calligraphy,
			&pc.Enchanting,
			&pc.Herbalism,
			&pc.Hunting,
			&pc.Mining,
			&pc.Bargaining,
			&pc.Camping,
			&pc.Firstaid,
			&pc.Lore,
			&pc.Picklocks,
			&pc.Scouting,
			&pc.Search,
			&pc.Stealth,
			&pc.Traps,
			&pc.Aeolandis,
			&pc.Hieroform,
			&pc.Highgundis,
			&pc.Oldpraxic,
			&pc.Praxic,
			&pc.Runic,
			&pc.Head,
			&pc.Chest,
			&pc.Arms,
			&pc.Hands,
			&pc.Legs,
			&pc.Feet,
			&pc.Cloak,
			&pc.Necklace,
			&pc.Ringone,
			&pc.Ringtwo,
			&pc.Righthand,
			&pc.Lefthand,
			&pc.Zone,
			&pc.X,
			&pc.Y,
			&pc.Z,
			&pc.Direction,
		)
		shared.CheckErr(err)
		pcs = append(pcs, pc)
	}
	return pcs
}

func pullNpcs() []shared.Npc {
	rows, err := db.Query("Select * FROM npcs")
	defer rows.Close()
	var npcs []shared.Npc
	for rows.Next() {
		var npc shared.Npc
		err = rows.Scan(
			&npc.Id,
			&npc.Firstname,
			&npc.Lastname,
			&npc.Guild,
			&npc.Race,
			&npc.Gender,
			&npc.Face,
			&npc.Skin,
			&npc.Profession,
			&npc.Alive,
			&npc.Nlevel,
			&npc.Hp,
			&npc.Maxhp,
			&npc.Bp,
			&npc.Maxbp,
			&npc.Mp,
			&npc.Maxmp,
			&npc.Ep,
			&npc.Maxep,
			&npc.Strength,
			&npc.Constitution,
			&npc.Intelligence,
			&npc.Dexterity,
			&npc.Head,
			&npc.Chest,
			&npc.Arms,
			&npc.Hands,
			&npc.Legs,
			&npc.Feet,
			&npc.Cloak,
			&npc.Righthand,
			&npc.Lefthand,
			&npc.Zone,
			&npc.X,
			&npc.Y,
			&npc.Z,
			&npc.Direction,
		)
		shared.CheckErr(err)
		npcs = append(npcs, npc)
	}
	return npcs
}

func pullItemTemplates() []shared.Item_template {
	rows, err := db.Query("Select * FROM item_templates")
	defer rows.Close()
	var item_templates []shared.Item_template
	for rows.Next() {
		var item_template shared.Item_template
		err = rows.Scan(
			&item_template.Id,
			&item_template.Name,
			&item_template.True_name,
			&item_template.Lore_level,
			&item_template.Item_type,
			&item_template.Slot,
			&item_template.Icon,
			&item_template.Weight,
			&item_template.Encumbrance,
			&item_template.Dyeable,
			&item_template.Stackable,
			&item_template.Stack_size,
			&item_template.Usable,
			&item_template.Equipable,
			&item_template.Base_price,
			&item_template.Strength_requirement,
			&item_template.Constitution_requirement,
			&item_template.Intelligence_requirement,
			&item_template.Dexterity_requirement,
			&item_template.Skill_type_0,
			&item_template.Skill_requirement_0,
			&item_template.Skill_type_1,
			&item_template.Skill_requirement_1,
			&item_template.Skill_type_2,
			&item_template.Skill_requirement_2,
			&item_template.Skill_type_3,
			&item_template.Skill_requirement_3,
			&item_template.Skill_type_4,
			&item_template.Skill_requirement_4,
			&item_template.Description,
			&item_template.True_description,
		)
		fmt.Println(&item_template.Name)
		shared.CheckErr(err)
		item_templates = append(item_templates, item_template)
	}
	return item_templates
}

func pullRaceTemplates() []shared.Race_template {
	rows, err := db.Query("Select * FROM race_templates")
	defer rows.Close()
	var race_templates []shared.Race_template
	for rows.Next() {
		var race_template shared.Race_template
		err = rows.Scan(
			&race_template.Id,
			&race_template.Race,
			&race_template.Strength_mod,
			&race_template.Constitution_mod,
			&race_template.Intelligence_mod,
			&race_template.Dexterity_mod,
			&race_template.Axe_mod,
			&race_template.Dagger_mod,
			&race_template.Unarmed_mod,
			&race_template.Hammer_mod,
			&race_template.Polearm_mod,
			&race_template.Spear_mod,
			&race_template.Staff_mod,
			&race_template.Sword_mod,
			&race_template.Archery_mod,
			&race_template.Crossbow_mod,
			&race_template.Sling_mod,
			&race_template.Thrown_mod,
			&race_template.Armor_mod,
			&race_template.Dualweapon_mod,
			&race_template.Shield_mod,
			&race_template.Bardic_mod,
			&race_template.Conjuring_mod,
			&race_template.Druidic_mod,
			&race_template.Illusion_mod,
			&race_template.Necromancy_mod,
			&race_template.Sorcery_mod,
			&race_template.Shamanic_mod,
			&race_template.Spellcraft_mod,
			&race_template.Summoning_mod,
			&race_template.Focus_mod,
			&race_template.Armorsmithing_mod,
			&race_template.Tailoring_mod,
			&race_template.Fletching_mod,
			&race_template.Weaponsmithing_mod,
			&race_template.Alchemy_mod,
			&race_template.Lapidary_mod,
			&race_template.Calligraphy_mod,
			&race_template.Enchanting_mod,
			&race_template.Herbalism_mod,
			&race_template.Hunting_mod,
			&race_template.Mining_mod,
			&race_template.Bargaining_mod,
			&race_template.Camping_mod,
			&race_template.Firstaid_mod,
			&race_template.Lore_mod,
			&race_template.Picklocks_mod,
			&race_template.Scouting_mod,
			&race_template.Search_mod,
			&race_template.Stealth_mod,
			&race_template.Traps_mod,
			&race_template.Aeolandis_mod,
			&race_template.Hieroform_mod,
			&race_template.Highgundis_mod,
			&race_template.Oldpraxic_mod,
			&race_template.Praxic_mod,
			&race_template.Runic_mod,
			&race_template.Description,
		)
		shared.CheckErr(err)
		race_templates = append(race_templates, race_template)
	}
	return race_templates
}

func pullProfessionTemplates() []shared.Profession_template {
	rows, err := db.Query("Select * FROM Profession_templates")
	defer rows.Close()
	var profession_templates []shared.Profession_template
	for rows.Next() {
		var profession_template shared.Profession_template
		err = rows.Scan(
			&profession_template.Id,
			&profession_template.Profession,
			&profession_template.Hppl,
			&profession_template.Mppl,
			&profession_template.Strength_mod,
			&profession_template.Constitution_mod,
			&profession_template.Intelligence_mod,
			&profession_template.Dexterity_mod,
			&profession_template.Axe_mod,
			&profession_template.Dagger_mod,
			&profession_template.Unarmed_mod,
			&profession_template.Hammer_mod,
			&profession_template.Polearm_mod,
			&profession_template.Spear_mod,
			&profession_template.Staff_mod,
			&profession_template.Sword_mod,
			&profession_template.Archery_mod,
			&profession_template.Crossbow_mod,
			&profession_template.Sling_mod,
			&profession_template.Thrown_mod,
			&profession_template.Armor_mod,
			&profession_template.Dualweapon_mod,
			&profession_template.Shield_mod,
			&profession_template.Bardic_mod,
			&profession_template.Conjuring_mod,
			&profession_template.Druidic_mod,
			&profession_template.Illusion_mod,
			&profession_template.Necromancy_mod,
			&profession_template.Sorcery_mod,
			&profession_template.Shamanic_mod,
			&profession_template.Spellcraft_mod,
			&profession_template.Summoning_mod,
			&profession_template.Focus_mod,
			&profession_template.Armorsmithing_mod,
			&profession_template.Tailoring_mod,
			&profession_template.Fletching_mod,
			&profession_template.Weaponsmithing_mod,
			&profession_template.Alchemy_mod,
			&profession_template.Lapidary_mod,
			&profession_template.Calligraphy_mod,
			&profession_template.Enchanting_mod,
			&profession_template.Herbalism_mod,
			&profession_template.Hunting_mod,
			&profession_template.Mining_mod,
			&profession_template.Bargaining_mod,
			&profession_template.Camping_mod,
			&profession_template.Firstaid_mod,
			&profession_template.Lore_mod,
			&profession_template.Picklocks_mod,
			&profession_template.Scouting_mod,
			&profession_template.Search_mod,
			&profession_template.Stealth_mod,
			&profession_template.Traps_mod,
			&profession_template.Aeolandis_mod,
			&profession_template.Hieroform_mod,
			&profession_template.Highgundis_mod,
			&profession_template.Oldpraxic_mod,
			&profession_template.Praxic_mod,
			&profession_template.Runic_mod,
			&profession_template.Skill_1_multi,
			&profession_template.Skill_1_names,
			&profession_template.Skill_1_value,
			&profession_template.Skill_2_multi,
			&profession_template.Skill_2_names,
			&profession_template.Skill_2_value,
			&profession_template.Skill_3_multi,
			&profession_template.Skill_3_names,
			&profession_template.Skill_3_value,
			&profession_template.Skill_4_multi,
			&profession_template.Skill_4_names,
			&profession_template.Skill_4_value,
			&profession_template.Skill_5_multi,
			&profession_template.Skill_5_names,
			&profession_template.Skill_5_value,
			&profession_template.Description,
		)
		shared.CheckErr(err)
		profession_templates = append(profession_templates, profession_template)
	}
	return profession_templates
}

func pushNpcs([]shared.Npc) {
	fmt.Println("Not yet implemented")
}

func getRemainingPlayerSlots(account_name string, max_player_slots int) int {
	rows, err := db.Query("Select * FROM players INNER JOIN accounts ON players.account_id=accounts.id WHERE accountname=?", account_name)
	shared.CheckErr(err)
	defer rows.Close()
	num_players := max_player_slots
	for rows.Next() {
		num_players--
	}
	return num_players
}

func playerFirstNameTaken(player_firstname string) bool {
	rows, err := db.Query("Select * FROM players WHERE firstname=?", player_firstname)
	shared.CheckErr(err)
	return foundInRows(rows, err)
}

func addNewPlayer(player shared.Pc) {
	// Need to add this...
	ins, err := db.Prepare("INSERT INTO players(account_id, firstname, lastname, guild, race, gender, face, skin, profession, alive, level, dp, hp, maxhp, bp, maxbp, mp, maxmp, ep, maxep, strength, constitution, intelligence, dexterity, axe, dagger, unarmed, hammer, polearm, spear, staff, sword, archery, crossbow, sling, thrown, armor, dualweapon, shield, bardic, conjuring, druidic, illusion, necromancy, sorcery, shamanic, spellcraft, summoning, focus, armorsmithing, tailoring, fletching, weaponsmithing, alchemy, lapidary, calligraphy, enchanting, herbalism, hunting, mining, bargaining, camping, firstaid, lore, picklocks, scouting, search, stealth, traps, aeolandis, hieroform, highgundis, oldpraxic, praxic, runic, head, chest, arms, hands, legs, feet, cloak, necklace, ringone, ringtwo, righthand, lefthand, zone, x, y, z, direction) values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	shared.CheckErr(err)
	// Try to populate and execute an SQL statment.
	_, err = ins.Exec(
		player.Accountid,
		player.Firstname,
		player.Lastname,
		player.Guild,
		player.Race,
		player.Gender,
		player.Face,
		player.Skin,
		player.Profession,
		player.Alive,
		player.Plevel,
		player.Dp,
		player.Hp,
		player.Maxhp,
		player.Bp,
		player.Maxbp,
		player.Mp,
		player.Maxmp,
		player.Ep,
		player.Maxep,
		player.Strength,
		player.Constitution,
		player.Intelligence,
		player.Dexterity,
		player.Axe,
		player.Dagger,
		player.Unarmed,
		player.Hammer,
		player.Polearm,
		player.Spear,
		player.Staff,
		player.Sword,
		player.Archery,
		player.Crossbow,
		player.Sling,
		player.Thrown,
		player.Armor,
		player.Dualweapon,
		player.Shield,
		player.Bardic,
		player.Conjuring,
		player.Druidic,
		player.Illusion,
		player.Necromancy,
		player.Sorcery,
		player.Shamanic,
		player.Spellcraft,
		player.Summoning,
		player.Focus,
		player.Armorsmithing,
		player.Tailoring,
		player.Fletching,
		player.Weaponsmithing,
		player.Alchemy,
		player.Lapidary,
		player.Calligraphy,
		player.Enchanting,
		player.Herbalism,
		player.Hunting,
		player.Mining,
		player.Bargaining,
		player.Camping,
		player.Firstaid,
		player.Lore,
		player.Picklocks,
		player.Scouting,
		player.Search,
		player.Stealth,
		player.Traps,
		player.Aeolandis,
		player.Hieroform,
		player.Highgundis,
		player.Oldpraxic,
		player.Praxic,
		player.Runic,
		player.Head,
		player.Chest,
		player.Arms,
		player.Hands,
		player.Legs,
		player.Feet,
		player.Cloak,
		player.Necklace,
		player.Ringone,
		player.Ringtwo,
		player.Righthand,
		player.Lefthand,
		player.Zone,
		player.X,
		player.Y,
		player.Z,
		player.Direction,
	)
	shared.CheckErr(err)
}
