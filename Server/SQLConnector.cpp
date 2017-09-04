// Filename:    SQLConnector.cpp
// Author:      Joseph DeVictoria
// Date:        Feb_29_2016
// Purpose:     SQL Connector class to handle SQLite3 queries.

#include "SQLConnector.h"
#include "Utils.h"
#include <iostream>
#include <string.h>
#include <sstream>

SQLConnector::SQLConnector() {
    sqls = sqlite3_open("db/Oldentide.db", &database);
    if (sqls) {
        std::cout << "Can't open database!" << std::endl;
        // Is this needed?
        sqlite3_close(database);
    }
    else {
        std::cout << "Oldentide database opened successfully." << std::endl;
    }
}

SQLConnector::~SQLConnector() {
    sqlite3_close(database);
}

// Executes a static SQL statement
int SQLConnector::Execute(std::string cmd) {
    char *errorMessage = NULL;
    sqls = sqlite3_exec(database, cmd.c_str(), &ExecuteCallback, 0, &errorMessage);
    if (sqls != SQLITE_OK) {
        std::cout << "Could not Execute SQL query! Return Code:" << sqls << std::endl;
    }
    if (errorMessage) {
        // Print out the error message if any
        std::cout << "SQL ERROR MESSAGE: " << errorMessage << std::endl;
        // Free the error message, since it was alloced in exec()
        sqlite3_free(errorMessage);
    }
    return 0;
}

// -------------------------------------- Parent Functions ----------------------------------------

// Insert a new account record.
bool SQLConnector::InsertAccount(char *accountName, char *key, char *salt) {
    std::stringstream query;
    // Sanitize key, salt, and account name
    if (!utils::SanitizeAccountName(accountName)) {
        std::cout << "Account Name is invalid! Cannot insert account record" << std::endl;
        return false;
    }
    if (!utils::SanitizeHexString(key)) {
        std::cout << "Key is invalid! Cannot insert account record" << std::endl;
        return false;
    }
    if (!utils::SanitizeHexString(salt)) {
        std::cout << "Salt is invalid! Cannot insert account record" << std::endl;
        return false;
    }
    query << "insert into accounts (accountname, key, salt) values (";
    query << "\"" << accountName << "\",";
    query << "\"" << key << "\",";
    query << "\"" << salt << "\")";
    Execute(query.str());
    if (sqls == SQLITE_OK) {
        return true;
    }
    else {
        return false;
    }
}

Account SQLConnector::GetAccountForSession(std::string sessionId) {
    // TODO Sanitize sessionId

    std::set<Account> accounts;
    std::stringstream query;
    char *errorMessage = NULL;

    query << "SELECT * FROM accounts";
    sqls = sqlite3_exec(database, query.str().c_str(), ExecuteCallback, (void*)&accounts, &errorMessage);
    if (sqls != SQLITE_OK) {
        std::cout << "Could not Execute SQL query! Return Code:" << sqls << std::endl;
    }
    if (errorMessage) {
        std::cout << "SQL ERROR MESSAGE: " << errorMessage << std::endl;
        sqlite3_free(errorMessage);
    }
    if (accounts.empty()) {
      throw std::runtime_error("No account with session");
    }
    return *accounts.begin();
}

// Inserts a new player into the database.
bool SQLConnector::InsertPlayer(Player p) {
    std::stringstream query;
    query << "INSERT INTO players (accountid, firstname, lastname, race, gender, face, skin, ";
    query << "profession, level, hp, maxhp, bp, maxbp, mp, maxmp, ep, maxep, strength, ";
    query << "constitution, intelligence, dexterity, axe, dagger, unarmed, hammer, polearm, ";
    query << "spear, staff, sword, archery, crossbow, sling, thrown, armor, dualweapon, shield, ";
    query << "bardic, conjuring, druidic, illusion, necromancy, sorcery, shamanic, summoning, ";
    query << "spellcraft, focus, armorsmithing, tailoring, fletching, weaponsmithing, ";
    query << "alchemy, lapidary, calligraphy, enchanting, herbalism, hunting, mining, ";
    query << "bargaining, camping, firstaid, lore, picklocks, scouting, search, stealth, ";
    query << "traps, aeolandis, hieroform, highgundis, oldpraxic, praxic, runic, head, ";
    query << "chest, arms, hands, legs, feet, cloak, necklace, ringone, ringtwo, righthand, ";
    query << "lefthand, zone, x, y, z, pitch, yaw) VALUES (";
    query << "'" << p.GetAccount() << "', ";
    query << "'" << p.GetName() << "', ";
    query << "'" << p.GetLastname() << "', ";
    query << "'" << p.GetRace() << "', ";
    query << "'" << p.GetGender() << "', ";
    query << "'" << p.GetFace() << "', ";
    query << "'" << p.GetSkin() << "', ";
    query << "'" << p.GetProfession() << "', ";
    query << p.GetLevel() << ", ";
    query << p.GetHp() << ", ";
    query << p.GetMaxHp() << ", ";
    query << p.GetBp() << ", ";
    query << p.GetMaxBp() << ", ";
    query << p.GetMp() << ", ";
    query << p.GetMaxMp() << ", ";
    query << p.GetEp() << ", ";
    query << p.GetMaxEp() << ", ";
    query << p.GetStrength() << ", ";
    query << p.GetConstitution() << ", ";
    query << p.GetIntelligence() << ", ";
    query << p.GetDexterity() << ", ";
    query << p.GetAxe() << ", ";
    query << p.GetDagger() << ", ";
    query << p.GetUnarmed() << ", ";
    query << p.GetHammer() << ", ";
    query << p.GetPolearm() << ", ";
    query << p.GetSpear() << ", ";
    query << p.GetStaff() << ", ";
    query << p.GetSword() << ", ";
    query << p.GetArchery() << ", ";
    query << p.GetCrossbow() << ", ";
    query << p.GetSling() << ", ";
    query << p.GetThrown() << ", ";
    query << p.GetArmor() << ", ";
    query << p.GetDualWeapon() << ", ";
    query << p.GetShield() << ", ";
    query << p.GetBardic() << ", ";
    query << p.GetConjuring() << ", ";
    query << p.GetDruidic() << ", ";
    query << p.GetIllusion() << ", ";
    query << p.GetNecromancy() << ", ";
    query << p.GetSorcery() << ", ";
    query << p.GetShamanic() << ", ";
    query << p.GetSummoning() << ", ";
    query << p.GetSpellcraft() << ", ";
    query << p.GetFocus() << ", ";
    query << p.GetArmorsmithing() << ", ";
    query << p.GetTailoring() << ", ";
    query << p.GetFletching() << ", ";
    query << p.GetWeaponsmithing() << ", ";
    query << p.GetAlchemy() << ", ";
    query << p.GetLapidary() << ", ";
    query << p.GetCalligraphy() << ", ";
    query << p.GetEnchanting() << ", ";
    query << p.GetHerbalism() << ", ";
    query << p.GetHunting() << ", ";
    query << p.GetMining() << ", ";
    query << p.GetBargaining() << ", ";
    query << p.GetCamping() << ", ";
    query << p.GetFirstAid() << ", ";
    query << p.GetLore() << ", ";
    query << p.GetPickLocks() << ", ";
    query << p.GetScouting() << ", ";
    query << p.GetSearch() << ", ";
    query << p.GetStealth() << ", ";
    query << p.GetTraps() << ", ";
    query << p.GetAeolandis() << ", ";
    query << p.GetHieroform() << ", ";
    query << p.GetHighGundis() << ", ";
    query << p.GetOldPraxic() << ", ";
    query << p.GetPraxic() << ", ";
    query << p.GetRunic() << ", ";
    query << "'" << p.GetHead() << "', ";
    query << "'" << p.GetChest() << "', ";
    query << "'" << p.GetArms() << "', ";
    query << "'" << p.GetHands() << "', ";
    query << "'" << p.GetLegs() << "', ";
    query << "'" << p.GetFeet() << "', ";
    query << "'" << p.GetCloak() << "', ";
    query << "'" << p.GetNecklace() << "', ";
    query << "'" << p.GetRingOne() << "', ";
    query << "'" << p.GetRingTwo() << "', ";
    query << "'" << p.GetRightHand() << "', ";
    query << "'" << p.GetLeftHand() << "', ";
    query << "'" << p.GetZone() << "', ";
    query << p.GetX() << ", ";
    query << p.GetY() << ", ";
    query << p.GetZ() << ", ";
    query << p.GetPitch() << ", ";
    query << p.GetYaw() << ")";
    std::cout << query.str() << std::endl;
    Execute(query.str());
    if (sqls == SQLITE_OK) {
        return true;
    }
    else {
        return false;
    }
}

// Lists all the accounts
void SQLConnector::ListAccounts() {
    std::stringstream query;
    query << "select * from accounts ORDER BY accountname";
    Execute(query.str());
}

std::vector<std::string> SQLConnector::GetPlayerList(std::string account) {
    std::vector<std::string> players;
    std::stringstream query;
    char *errorMessage = NULL;
    query << "SELECT firstname, lastname ";
    query << "FROM players JOIN accounts ";
    query << "ON players.accountid = accounts.id ";
    query << "WHERE accountname = \"" << account << "\"";
    sqls = sqlite3_exec(database, query.str().c_str(), ParsePlayerList, (void*)&players, &errorMessage);
    if (sqls != SQLITE_OK) {
        std::cout << "Could not Execute SQL query! Return Code:" << sqls << std::endl;
    }
    if (errorMessage) {
        std::cout << "SQL ERROR MESSAGE: " << errorMessage << std::endl;
        sqlite3_free(errorMessage);
    }
    return players;
}

std::set<Npc> SQLConnector::PopulateNpcList() {
    std::set<Npc> npcs;
    std::stringstream query;
    char *errorMessage = NULL;
    query << "SELECT * FROM npcs";
    sqls = sqlite3_exec(database, query.str().c_str(), ExecuteCallback, (void*)&npcs, &errorMessage);
    if (sqls != SQLITE_OK) {
        std::cout << "Could not Execute SQL query! Return Code:" << sqls << std::endl;
    }
    if (errorMessage) {
        std::cout << "SQL ERROR MESSAGE: " << errorMessage << std::endl;
        sqlite3_free(errorMessage);
    }
    return npcs;
}

// Return the salt for the passed account
bool SQLConnector::GetAccountSalt(char *accountName, char *saltStringHex) {
    // Sanitize the account name before preceeding
    if (!utils::SanitizeAccountName(accountName)) {
        return false;
    }
    char *errorMessage = NULL;
    std::stringstream query;
    query << "select salt from accounts where accountname = \"" << accountName << "\"";
    // The fourth param is passed to the callback function as a void pointer to the first param
    sqls = sqlite3_exec(database, query.str().c_str(), &ReturnStringCallback, saltStringHex, &errorMessage);
    if (sqls != SQLITE_OK) {
        std::cout << "Could not Execute SQL query! Return Code:" << sqls << std::endl;
    }
    if (errorMessage) {
        // Print out the error message if any
        std::cout << "SQL ERROR MESSAGE: " << errorMessage << std::endl;
        // Free the error message, since it was alloced in exec()
        sqlite3_free(errorMessage);
        return false;
    }
    // Check to see if the salt was retrieved
    if (utils::SanitizeHexString(saltStringHex)) {
        return true;
    }
    else {
        return false;
    }

}

// Returns the key to the passed account
int SQLConnector::GetAccountKey(char *accountName, char *keyStringHex) {
    // Sanitize the account name before preceeding
    if (!utils::SanitizeAccountName(accountName)) {
        return 0;
    }
    char *errorMessage = NULL;
    std::stringstream query;
    query << "select key from accounts where accountname = \"" << accountName << "\"";
    //std::cout << query.str() << std::endl;
    // The fourth param is passed to the callback function as a void pointer to the first param
    sqls = sqlite3_exec(database, query.str().c_str(), &ReturnStringCallback, keyStringHex, &errorMessage);
    if (sqls != SQLITE_OK) {
        std::cout << "Could not Execute SQL query! Return Code:" << sqls << std::endl;
    }
    if (errorMessage) {
        // Print out the error message if any
        std::cout << "SQL ERROR MESSAGE: " << errorMessage << std::endl;
        // Free the error message, since it was alloced in exec()
        sqlite3_free(errorMessage);
    }
    return 1;
}

// -------------------------------------- Callback Functions ---------------------------------------

// A generic callback function to sqlite3_exec() that copies a c string in the first column
// of the returned row into stringToReturn. This function assumes that only one record will be
// in the result set, or else the return value will be the value in the last row processed.
static int ReturnStringCallback(void *stringToReturn, int argc, char **argv, char **azColName) {
    strcpy((char *)stringToReturn, argv[0]);
    return 0;
}

// This is a generic callback that simply prints to screen all the values
// of the returned row after a query. Callback adapted from https://www.sqlite.org/quickstart.html
static int ExecuteCallback(void *NotUsed, int argc, char **argv, char **azColName) {
    int i;
    for(i = 0; i < argc; i++) {
        std::cout << azColName[i] << " = " << (argv[i] ? argv[i] : "NULL");
        std::cout << " | ";
    }
    std::cout << std::endl;
    return 0;
}

static int ParsePlayerList(void * players, int argc, char ** argv, char ** azColName) {
    std::vector<std::string> * playerList = (std::vector<std::string> *)players;
    playerList->push_back(std::string(argv[0]).append(" ").append(std::string(argv[1])));
    return 0;
}

static int ParseNpcs(void * npcs, int argc, char ** argv, char ** azColName) {
    std::set<Npc> * npcset = (std::set<Npc> *)npcs;
    return 0;
}
