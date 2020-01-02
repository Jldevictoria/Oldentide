// Filename:    sql_connector.go
// Author:      Joseph DeVictoria
// Date:        June_16_2018
// Purpose:     Database / sql abstraction for Oldentide dedicated server.

package main

import (
	"Oldentide/shared"
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

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

func createAccount(accountname string, email string, verifyKey string, hashedKey string, saltKey string) bool {
	// Prepare insert statement.
	ins, err := db.Prepare("INSERT INTO accounts(valid, banned, accountname, email, gamesession, playing, verify, hash, salt) values(?, ?, ?, ?, ?, ?, ?, ?, ?)")
	shared.CheckErr(err)
	// Try to populate and execute an SQL statment.
	_, err = ins.Exec("0", "0", accountname, email, "0", "0", verifyKey, hashedKey, saltKey)
	if err == nil {
		return true
	}
	return false
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

func setSessionID(account string, sessionID int64) bool {
	update, err := db.Prepare("UPDATE accounts SET gamesession=? WHERE accountname=?")
	_, err = update.Exec(sessionID, account)
	if err == nil {
		return true
	}
	return false
}

func activateAccount(a string) bool {
	update, err := db.Prepare("UPDATE accounts SET valid=? WHERE accountname=?")
	_, err = update.Exec("1", a)
	if err == nil {
		return true
	}
	return false
}

func banAccount(a string) bool {
	ban, err := db.Prepare("UPDATE accounts SET banned=? WHERE accountname=?")
	_, err = ban.Exec("1", a)
	if err == nil {
		return true
	}
	return false
}

// Queries the db to make sure that we are generating a truly unique verify key.
func generateUniqueVerify(n int) string {
	findKey := true
	for findKey {
		verifyKey := shared.GenerateRandomLetters(n)
		rows, err := db.Query("SELECT accountname FROM accounts WHERE verify='" + verifyKey + "'")
		//ifErrPrintErr(err)
		if !foundInRows(rows, err) {
			return verifyKey
		}
	}
	return ""
}

// Queries the db to make sure that we are generating a truly unique salt.
func generateUniqueSalt(n int) string {
	findKey := true
	for findKey {
		saltKey := shared.GenerateRandomLetters(n)
		rows, err := db.Query("SELECT accountname FROM accounts WHERE salt='" + saltKey + "'")
		if !foundInRows(rows, err) {
			return saltKey
		}
	}
	return ""
}

func generateUniqueSessionID() int64 {
	findSession := true
	for findSession {
		randomSession := strconv.FormatInt(rand.Int63(), 10)
		rows, err := db.Query("SELECT accountname FROM accounts WHERE gamesession='" + randomSession + "'")
		if !foundInRows(rows, err) {
			sessionID, err := strconv.ParseInt(randomSession, 10, 64)
			shared.CheckErr(err)
			return sessionID
		}
	}
	return 0
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

func pcRowsToStruct(rows *sql.Rows) []shared.Pc {
	var pcs []shared.Pc
	for rows.Next() {
		var pc shared.Pc
		err = rows.Scan(
			&pc.ID,
			&pc.AccountID,
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
			&npc.ID,
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

func pullItemTemplates() []shared.ItemTemplate {
	rows, err := db.Query("Select * FROM item_templates")
	defer rows.Close()
	var itemTemplates []shared.ItemTemplate
	for rows.Next() {
		var itemTemplate shared.ItemTemplate
		err = rows.Scan(
			&itemTemplate.ID,
			&itemTemplate.Name,
			&itemTemplate.TrueName,
			&itemTemplate.LoreLevel,
			&itemTemplate.ItemType,
			&itemTemplate.Slot,
			&itemTemplate.Icon,
			&itemTemplate.Weight,
			&itemTemplate.Encumbrance,
			&itemTemplate.Dyeable,
			&itemTemplate.Stackable,
			&itemTemplate.StackSize,
			&itemTemplate.Usable,
			&itemTemplate.Equipable,
			&itemTemplate.BasePrice,
			&itemTemplate.StrengthRequirement,
			&itemTemplate.ConstitutionRequirement,
			&itemTemplate.IntelligenceRequirement,
			&itemTemplate.DexterityRequirement,
			&itemTemplate.SkillType0,
			&itemTemplate.SkillRequirement0,
			&itemTemplate.SkillType1,
			&itemTemplate.SkillRequirement1,
			&itemTemplate.SkillType2,
			&itemTemplate.SkillRequirement2,
			&itemTemplate.SkillType3,
			&itemTemplate.SkillRequirement3,
			&itemTemplate.SkillType4,
			&itemTemplate.SkillRequirement4,
			&itemTemplate.Description,
			&itemTemplate.TrueDescription,
		)
		shared.CheckErr(err)
		itemTemplates = append(itemTemplates, itemTemplate)
	}
	return itemTemplates
}

func pullSpellTemplates() []shared.SpellTemplate {
	rows, err := db.Query("Select * FROM spell_templates")
	defer rows.Close()
	var spellTemplates []shared.SpellTemplate
	for rows.Next() {
		var spellTemplate shared.SpellTemplate
		err = rows.Scan(
			&spellTemplate.ID,
			&spellTemplate.Spellname,
			&spellTemplate.School,
			&spellTemplate.Level,
			&spellTemplate.Type,
			&spellTemplate.Target,
			&spellTemplate.Range,
			&spellTemplate.Accuracy,
			&spellTemplate.PreparationTime,
			&spellTemplate.RecoveryTime,
			&spellTemplate.Effect1,
			&spellTemplate.Effect2,
			&spellTemplate.Effect3,
			&spellTemplate.Effect4,
			&spellTemplate.Effect5,
			&spellTemplate.Description,
		)
		shared.CheckErr(err)
		spellTemplates = append(spellTemplates, spellTemplate)
	}
	return spellTemplates
}

func pullRaceTemplates() []shared.RaceTemplate {
	rows, err := db.Query("Select * FROM race_templates")
	defer rows.Close()
	var raceTemplates []shared.RaceTemplate
	for rows.Next() {
		var raceTemplate shared.RaceTemplate
		err = rows.Scan(
			&raceTemplate.ID,
			&raceTemplate.Race,
			&raceTemplate.StrengthMod,
			&raceTemplate.ConstitutionMod,
			&raceTemplate.IntelligenceMod,
			&raceTemplate.DexterityMod,
			&raceTemplate.AxeMod,
			&raceTemplate.DaggerMod,
			&raceTemplate.UnarmedMod,
			&raceTemplate.HammerMod,
			&raceTemplate.PolearmMod,
			&raceTemplate.SpearMod,
			&raceTemplate.StaffMod,
			&raceTemplate.SwordMod,
			&raceTemplate.ArcheryMod,
			&raceTemplate.CrossbowMod,
			&raceTemplate.SlingMod,
			&raceTemplate.ThrownMod,
			&raceTemplate.ArmorMod,
			&raceTemplate.DualweaponMod,
			&raceTemplate.ShieldMod,
			&raceTemplate.BardicMod,
			&raceTemplate.ConjuringMod,
			&raceTemplate.DruidicMod,
			&raceTemplate.IllusionMod,
			&raceTemplate.NecromancyMod,
			&raceTemplate.SorceryMod,
			&raceTemplate.ShamanicMod,
			&raceTemplate.SpellcraftMod,
			&raceTemplate.SummoningMod,
			&raceTemplate.FocusMod,
			&raceTemplate.ArmorsmithingMod,
			&raceTemplate.TailoringMod,
			&raceTemplate.FletchingMod,
			&raceTemplate.WeaponsmithingMod,
			&raceTemplate.AlchemyMod,
			&raceTemplate.LapidaryMod,
			&raceTemplate.CalligraphyMod,
			&raceTemplate.EnchantingMod,
			&raceTemplate.HerbalismMod,
			&raceTemplate.HuntingMod,
			&raceTemplate.MiningMod,
			&raceTemplate.BargainingMod,
			&raceTemplate.CampingMod,
			&raceTemplate.FirstaidMod,
			&raceTemplate.LoreMod,
			&raceTemplate.PicklocksMod,
			&raceTemplate.ScoutingMod,
			&raceTemplate.SearchMod,
			&raceTemplate.StealthMod,
			&raceTemplate.TrapsMod,
			&raceTemplate.AeolandisMod,
			&raceTemplate.HieroformMod,
			&raceTemplate.HighgundisMod,
			&raceTemplate.OldpraxicMod,
			&raceTemplate.PraxicMod,
			&raceTemplate.RunicMod,
			&raceTemplate.Description,
		)
		shared.CheckErr(err)
		raceTemplates = append(raceTemplates, raceTemplate)
	}
	return raceTemplates
}

func pullProfessionTemplates() []shared.ProfessionTemplate {
	rows, err := db.Query("Select * FROM profession_templates")
	defer rows.Close()
	var professionTemplates []shared.ProfessionTemplate
	for rows.Next() {
		var professionTemplate shared.ProfessionTemplate
		err = rows.Scan(
			&professionTemplate.ID,
			&professionTemplate.Profession,
			&professionTemplate.Hppl,
			&professionTemplate.Mppl,
			&professionTemplate.StrengthMod,
			&professionTemplate.ConstitutionMod,
			&professionTemplate.IntelligenceMod,
			&professionTemplate.DexterityMod,
			&professionTemplate.AxeMod,
			&professionTemplate.DaggerMod,
			&professionTemplate.UnarmedMod,
			&professionTemplate.HammerMod,
			&professionTemplate.PolearmMod,
			&professionTemplate.SpearMod,
			&professionTemplate.StaffMod,
			&professionTemplate.SwordMod,
			&professionTemplate.ArcheryMod,
			&professionTemplate.CrossbowMod,
			&professionTemplate.SlingMod,
			&professionTemplate.ThrownMod,
			&professionTemplate.ArmorMod,
			&professionTemplate.DualweaponMod,
			&professionTemplate.ShieldMod,
			&professionTemplate.BardicMod,
			&professionTemplate.ConjuringMod,
			&professionTemplate.DruidicMod,
			&professionTemplate.IllusionMod,
			&professionTemplate.NecromancyMod,
			&professionTemplate.SorceryMod,
			&professionTemplate.ShamanicMod,
			&professionTemplate.SpellcraftMod,
			&professionTemplate.SummoningMod,
			&professionTemplate.FocusMod,
			&professionTemplate.ArmorsmithingMod,
			&professionTemplate.TailoringMod,
			&professionTemplate.FletchingMod,
			&professionTemplate.WeaponsmithingMod,
			&professionTemplate.AlchemyMod,
			&professionTemplate.LapidaryMod,
			&professionTemplate.CalligraphyMod,
			&professionTemplate.EnchantingMod,
			&professionTemplate.HerbalismMod,
			&professionTemplate.HuntingMod,
			&professionTemplate.MiningMod,
			&professionTemplate.BargainingMod,
			&professionTemplate.CampingMod,
			&professionTemplate.FirstaidMod,
			&professionTemplate.LoreMod,
			&professionTemplate.PicklocksMod,
			&professionTemplate.ScoutingMod,
			&professionTemplate.SearchMod,
			&professionTemplate.StealthMod,
			&professionTemplate.TrapsMod,
			&professionTemplate.AeolandisMod,
			&professionTemplate.HieroformMod,
			&professionTemplate.HighgundisMod,
			&professionTemplate.OldpraxicMod,
			&professionTemplate.PraxicMod,
			&professionTemplate.RunicMod,
			&professionTemplate.SkillMulti1,
			&professionTemplate.SkillNames1,
			&professionTemplate.SkillValue1,
			&professionTemplate.SkillMulti2,
			&professionTemplate.SkillNames2,
			&professionTemplate.SkillValue2,
			&professionTemplate.SkillMulti3,
			&professionTemplate.SkillNames3,
			&professionTemplate.SkillValue3,
			&professionTemplate.SkillMulti4,
			&professionTemplate.SkillNames4,
			&professionTemplate.SkillValue4,
			&professionTemplate.SkillMulti5,
			&professionTemplate.SkillNames5,
			&professionTemplate.SkillValue5,
			&professionTemplate.Description,
		)
		shared.CheckErr(err)
		professionTemplates = append(professionTemplates, professionTemplate)
	}
	return professionTemplates
}

func pushNpcs([]shared.Npc) {
	fmt.Println("Not yet implemented")
}

func getCharacterList(accountName string) []string {
	rows, err := db.Query("SELECT firstname FROM players INNER JOIN accounts ON players.accountID=accounts.id WHERE accountname=?", accountName)
	shared.CheckErr(err)
	defer rows.Close()
	var accountCharacters []string
	for rows.Next() {
		var characterName string
		err = rows.Scan(&characterName)
		shared.CheckErr(err)
		accountCharacters = append(accountCharacters, characterName)
	}
	return accountCharacters
}

func getRemainingPlayerSlots(accountName string, maxPlayerSlots int) int {
	rows, err := db.Query("SELECT * FROM players INNER JOIN accounts ON players.accountID=accounts.id WHERE accountname=?", accountName)
	shared.CheckErr(err)
	defer rows.Close()
	numPlayers := maxPlayerSlots
	for rows.Next() {
		numPlayers--
	}
	return numPlayers
}

func playerFirstNameTaken(playerFirstname string) bool {
	rows, err := db.Query("SELECT * FROM players WHERE firstname=?", playerFirstname)
	shared.CheckErr(err)
	return foundInRows(rows, err)
}

func getAccountIDFromAccountName(accountName string) int32 {
	rows, err := db.Query("SELECT id FROM accounts WHERE accountname=?", accountName)
	shared.CheckErr(err)
	defer rows.Close()
	var accountID int32
	for rows.Next() {
		err = rows.Scan(&accountID)
		shared.CheckErr(err)
	}
	return accountID
}

func addNewPlayer(player shared.Pc) {
	// Need to add this...
	ins, err := db.Prepare("INSERT INTO players(accountID, firstname, lastname, guild, race, gender, face, skin, profession, alive, level, dp, hp, maxhp, bp, maxbp, mp, maxmp, ep, maxep, strength, constitution, intelligence, dexterity, axe, dagger, unarmed, hammer, polearm, spear, staff, sword, archery, crossbow, sling, thrown, armor, dualweapon, shield, bardic, conjuring, druidic, illusion, necromancy, sorcery, shamanic, spellcraft, summoning, focus, armorsmithing, tailoring, fletching, weaponsmithing, alchemy, lapidary, calligraphy, enchanting, herbalism, hunting, mining, bargaining, camping, firstaid, lore, picklocks, scouting, search, stealth, traps, aeolandis, hieroform, highgundis, oldpraxic, praxic, runic, head, chest, arms, hands, legs, feet, cloak, necklace, ringone, ringtwo, righthand, lefthand, zone, x, y, z, direction) values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	shared.CheckErr(err)
	// Try to populate and execute an SQL statment.
	_, err = ins.Exec(
		player.AccountID,
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
