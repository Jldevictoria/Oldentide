// Filename:    GameState.h
// Author:      Joseph DeVictoria
// Date:        Mar_2_2016
// Purpose:     Header for Game State class.
//              This class will be responsible for maintaining the Players,
//              NPCs necessary to preserve the game state.  It will also
//              perform the bulk of database interactions.

#ifndef OLDENTIDE_GAMESTATE_H
#define OLDENTIDE_GAMESTATE_H

#include "Account.h"
#include "Npc.h"
#include "Player.h"
#include "SQLConnector.h"
#include <set>
#include <string>
#include <vector>
#include <map>

class Server;

class GameState{
    public:
        GameState(Server * server, SQLConnector * sql);
        ~GameState();
        bool VerifySession(std::string sessionId);
        bool VerifyActiveSession(std::string sessionId);
        void DisconnectSession(std::string sessionId);
        void PlayerCommand(std::string command, std::string sessionId);
        void SelectPlayer(std::string sessionId);

        void SetSessionAccountName(char * accountName, std::string sessionId);
        std::string GetSessionAccountName(std::string sessionId);
        std::set<Player> getPlayers();
        std::set<Npc> getNPCs();

    private:
        SQLConnector * sql;
        Server * server;
        std::set<Player> players;
        std::set<Npc> npcs;
        std::map<std::string, std::string> sessionAccounts;
        int curSession;
        //Player ReadPlayer(std::string name);
        void StorePlayer(std::string name);
};

#endif // OLDENTIDE_GAMESTATE_H
