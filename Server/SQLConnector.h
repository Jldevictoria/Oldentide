// Filename:    SQLConnector.h
// Author:      Joseph DeVictoria
// Date:        Feb_29_2016
// Purpose:     Header for SQL Connector class.

#ifndef OLDENTIDE_SQLCONNECTOR_H
#define OLDENTIDE_SQLCONNECTOR_H

#include "Account.h"
#include "Player.h"
#include "Npc.h"
#include <string>
#include <sqlite3.h>
#include <vector>
#include <set>

class SQLConnector{
    public:
        SQLConnector();
        ~SQLConnector();
        int Execute(std::string cmd);
        bool InsertAccount(char *, char *, char *);
        bool InsertPlayer(Player newPlayer);
        Account GetAccountForSession(std::string sessionId);
        void ListAccounts();
        std::vector<std::string> GetPlayerList(std::string account);
        std::set<Npc> PopulateNpcList();
        bool GetAccountSalt(char *, char *);
        int GetAccountKey(char *, char *);
    private:
        sqlite3 * database;
        int sqls;
};

// This can't be a private method, because it needs to be passed as a pointer to sqlite3_exec()
// and methods of a c++ class have a hidden function pointer to the class instance
// See http://stackoverflow.com/a/8045331
static int ExecuteCallback(void *, int, char **, char **);
static int ReturnStringCallback(void *, int, char **, char **);
static int ParsePlayerList(void *, int, char **, char **);
static int ParseNpcs(void *, int, char **, char **);

#endif // OLDENTIDE_SQLCONNECTOR_H
