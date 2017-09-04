// Filename:    Server.h
// Author:      Joseph DeVictoria
// Date:        Jan_31_2016
// Purpose:     Header for dedicated server class.
//              This server is a traditional Universal Datagram Protocol
//              (UDP) Server which will receive and handle all game
//              packets sent from client.

#ifndef OLDENTIDE_OLDENTIDESERVER_H
#define OLDENTIDE_OLDENTIDESERVER_H

#include "Npc.h"
#include "Packets.h"
#include "Player.h"
#include "SQLConnector.h"
#include <arpa/inet.h>
#include <set>
#include <string>
#include <mutex>
#include <queue>
#include <condition_variable>

// forward declaration so AdminShell and Server can see each other
class AdminShell;
class GameState;

class Server{
    public:
        Server(int port);
        ~Server();
        void Run();
        int GetPacketQueueSize();
        void BroadcastToConnections(std::string msg, std::string fromUserSessionId);
        void SendMessageToConnection(std::string msg, std::string fromUserSessionId, std::string toUserSessionId);
    private:
        int sockfd;
        SQLConnector * sql;
        GameState * gameState;
        AdminShell * adminshell;
        std::queue<packets::packet_t> packetQueue;
        // Session to socket mapping
        std::map<std::string, sockaddr_in> activeConnections;

        // Wrap all globals accessed in WorkerThread with mutexes
        std::mutex gameStateMutex;
        std::mutex packetQueueMutex;
        std::mutex activeConnectionsMutex;
        // NOTE: Sockets are atomic, so they are thread safe
        // https://stackoverflow.com/questions/1981372/are-parallel-calls-to-send-recv-on-the-same-socket-valid

        std::condition_variable packetQueueSyncVar;

        void WorkerThread(int id);
        // void StatisticsThread();
        void GenericHandler(msgpack::object_handle *data, sockaddr_in *client);
        void AckHandler(msgpack::object_handle *data, sockaddr_in *client);
        void ConnectHandler(msgpack::object_handle *data, sockaddr_in *client);
        void DisconnectHandler(msgpack::object_handle *data, sockaddr_in *client);
        void ListCharactersHandler(msgpack::object_handle *data, sockaddr_in *client);
        void SelectCharacterHandler(msgpack::object_handle *data, sockaddr_in *client);
        void DeleteCharacterHandler(msgpack::object_handle *data, sockaddr_in *client);
        void CreateCharacterHandler(msgpack::object_handle *data, sockaddr_in *client);
        void InitializeGameHandler(msgpack::object_handle *data, sockaddr_in *client);
        void UpdatePcHandler(msgpack::object_handle *data, sockaddr_in *client);
        void UpdateNpcHandler(msgpack::object_handle *data, sockaddr_in *client);
        void SendPlayerCommandHandler(msgpack::object_handle *data, sockaddr_in *client);
        void SendPlayerActionHandler(msgpack::object_handle *data, sockaddr_in *client);
        void SendServerActionHandler(msgpack::object_handle *data, sockaddr_in *client);
        void UnityHandler(msgpack::object_handle *data, sockaddr_in *client);
};

#endif //OLDENTIDE_SERVER_H
