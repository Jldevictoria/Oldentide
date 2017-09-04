// Filename:    Server.cpp
// Author:      Joseph DeVictoria
// Date:        Jan_31_2016
// Purpose:     Dedicated server class.

#include "Server.h"
#include "AdminShell.h"
#include "GameState.h"
#include <cstring>
#include <iostream>
#include <iterator>
#include <netinet/in.h>
#include <sys/socket.h>
#include <sys/types.h>
#include <thread>
#include <unistd.h>
#include <limits.h>
#include <msgpack.hpp>
#include <sstream>
#include "Utils.h"
#include <queue>
#include <mutex>
#include <chrono>

// TODO: Create a helper function to reduce try/catch code with messagepack

Server::Server(int port) {
    sql = new SQLConnector();
    gameState = new GameState(this, sql);
    adminshell = new AdminShell(this, sql, gameState);

    // Create server address struct.
    sockaddr_in server;
    server.sin_family = AF_INET;
    server.sin_addr.s_addr = htonl(INADDR_ANY);
    server.sin_port = htons(port);

    // Create socket for IP, UDP normal protocol.
    if ((sockfd = socket(AF_INET, SOCK_DGRAM, IPPROTO_UDP)) < 0) {
        std::cout << "Cannot create socket..." << std::endl;
        exit(EXIT_FAILURE);
    }

    // Bind socket to a port.
    if ((bind(sockfd, (struct sockaddr *)&server, sizeof(server))) < 0) {
        std::cout << "Cannot bind socket..." << std::endl;
        exit(EXIT_FAILURE);
    };

    std::cout << "Server socket: " << utils::GetIpAndPortFromSocket(&server) <<  std::endl;
}

Server::~Server() {
    delete sql;
    delete gameState;
    delete adminshell;
}

void Server::Run() {
    std::thread shell(*adminshell);

    // TODO: There should be 1 active thread per logical core to prevent thrashing
    int NUM_WORKER_THREADS = 4;

    // Initialize all the threads
    for (int i = 0; i < NUM_WORKER_THREADS; ++i) {
        // Start a thread on a Server member function
        // Funky syntax; server object needs to be second param
        // See http://stackoverflow.com/questions/10673585/start-thread-with-member-function
        std::thread t(&Server::WorkerThread, this, i);
        t.detach(); // Detach the threads so we don't need to join them manually
    }

    // // Set up a statistics thread, to monitor the server
    // std::thread statThread(&Server::StatisticsThread, this);
    // statThread.detach(); // Detach the threads so we don't need to join them manually

    bool listen = true;
    while(listen) {
        packets::packet_t p;

        // Wait for a packet
        utils::ReceivePacketFrom(sockfd, p.data, &(p.source));

        // Put it in the packetQueue
        packetQueueMutex.lock();
        packetQueue.push(p);
        packetQueueMutex.unlock();
        // Use a condition variable to notify the next worker thread to wake up
        // https://stackoverflow.com/questions/3513045/conditional-variable-vs-semaphore
        packetQueueSyncVar.notify_one();
    }

    shell.join();
}


void Server::WorkerThread(int id) {
    while(1){
        packets::packet_t packet;

        // The following 2 lines will lock, release lock, and block thread
        // until notify_one() is called after a packet arrives
        // http://en.cppreference.com/w/cpp/thread/condition_variable/wait
        std::unique_lock<std::mutex> ul(packetQueueMutex);
        packetQueueSyncVar.wait(ul);
        // std::cout << "Worker thread " << id << " handling the request!" << std::endl;
        // Retrieve a packet from the queue
        // create a copy of the packet
        packet = packetQueue.front();
        packetQueue.pop();
        packetQueueMutex.unlock();

        msgpack::object_handle deserialized_data;
        uint8_t packetType;
        deserialized_data = utils::GetDataFromPacket(packet.data, &packetType);

        switch (packetType) {
            case false:
                utils::SendErrorTo(sockfd, std::string("Invalid packetType. Ignoring packet."), &(packet.source));
                break;
            case packets::GENERIC:
                GenericHandler(&deserialized_data, &(packet.source));
                break;
            case packets::ACK:
                AckHandler(&deserialized_data, &(packet.source));
                break;
            case packets::CONNECT:
                ConnectHandler(&deserialized_data, &(packet.source));
                break;
            case packets::DISCONNECT:
                DisconnectHandler(&deserialized_data, &(packet.source));
                break;
            case packets::LISTCHARACTERS:
                ListCharactersHandler(&deserialized_data, &(packet.source));
                break;
            case packets::SELECTCHARACTER:
                SelectCharacterHandler(&deserialized_data, &(packet.source));
                break;
            case packets::DELETECHARACTER:
                DeleteCharacterHandler(&deserialized_data, &(packet.source));
                break;
            case packets::CREATECHARACTER:
                CreateCharacterHandler(&deserialized_data, &(packet.source));
                break;
            case packets::INITIALIZEGAME:
                InitializeGameHandler(&deserialized_data, &(packet.source));
                break;
            case packets::UPDATEPC:
                UpdatePcHandler(&deserialized_data, &(packet.source));
                break;
            case packets::UPDATENPC:
                UpdateNpcHandler(&deserialized_data, &(packet.source));
                break;
            case packets::SENDPLAYERCOMMAND:
                SendPlayerCommandHandler(&deserialized_data, &(packet.source));
                break;
            case packets::SENDPLAYERACTION:
                SendPlayerActionHandler(&deserialized_data, &(packet.source));
                break;
            case packets::SENDSERVERACTION:
                SendServerActionHandler(&deserialized_data, &(packet.source));
                break;
            case packets::UNITY:
                UnityHandler(&deserialized_data, &(packet.source));
                break;
            default:
                std::stringstream errSs;
                errSs << "Received unknown packet of type " << packetType;
                utils::SendErrorTo(sockfd, errSs.str(), &(packet.source));
                break;
        }
    }
}

// TODO: This can be used to eventually log things to a file
// OR, this thread could be in charge of sending PING requests to all
// actively connected users, and if a ACK request is not sent in a timely
// fashion (10 seconds?), then they are considered disconnected.
// // Print out the size of the packetQueue every so often
// void Server::StatisticsThread() {
//     while(1){
//         // Get the size
//         int size = GetPacketQueueSize();
//         // Print out the size of the packet queue
//         std::cout << "packetQueue size: " << size << std::endl;
//         std::this_thread::sleep_for(std::chrono::milliseconds(200));
//     }
// }

int Server::GetPacketQueueSize() {
    packetQueueMutex.lock();
    int size = packetQueue.size();
    packetQueueMutex.unlock();
    return size;
}

void Server::BroadcastToConnections(std::string msg, std::string fromUserSessionId){
    // For each connected client, send a SendServerCommand packet!
    bool isAdmin = false;
    std::string broadcast;
    if(fromUserSessionId == "admin"){
        isAdmin = true;
        broadcast = ">>>BROADCAST Oldentide Admin: " + msg + " <<<";
    }
    else {
        broadcast = ">>>BROADCAST Client " + fromUserSessionId + ": " + msg + " <<<";
        // Print the broadcast so the server admin can see it as well
        std::cout << broadcast << std::endl;
    }
    // std::cout << "Broadcasting to sessions:" << std::endl;


    std::lock_guard<std::mutex> lk(activeConnectionsMutex);
    for (std::map<std::string, sockaddr_in>::iterator it = activeConnections.begin(); it != activeConnections.end(); ++it){
        // Don't broadcast to the originator of the message
        if(!isAdmin && it->first == fromUserSessionId){
            continue;
        }
        // std::cout << it->first;

        packets::Sendservercommand packet;
        packet.sessionId = it->first;
        packet.packetId = 0; // n/a
        packet.command = broadcast;
        std::stringstream buffer;
        msgpack::pack(buffer, packet);
        utils::SendDataTo(sockfd, &buffer, packets::SENDSERVERCOMMAND, &(it->second));
    }
}


void Server::SendMessageToConnection(std::string msg, std::string fromUserSessionId, std::string toUserSessionId){
    bool isAdmin = false;
    std::string formattedMsg;

    // Format the message
    if(fromUserSessionId == "admin"){
        formattedMsg = ">>>MESSAGE Oldentide Admin: " + msg + " <<<";
    }
    else {
        formattedMsg = ">>>MESSAGE Client " + fromUserSessionId + ": " + msg + " <<<";
    }

    // Figure out who to send the message to
    if(toUserSessionId == "admin"){
        // If message is sent to the admin, simply print to the admin shell
        isAdmin = true;
        std::cout << formattedMsg << std::endl;
    }
    else {
        // user should be an int
        std::string sessionId = toUserSessionId;

        std::lock_guard<std::mutex> lk(activeConnectionsMutex);
        std::map<std::string, sockaddr_in>::iterator connection = activeConnections.find(sessionId);
        if(connection == activeConnections.end()) {
            std::cout << "Could not deliver message to session '" << toUserSessionId << ": does not exist" << std::endl;
            // TODO: Send "could not deliver" error message?
            // utils::SendErrorTo(sockfd, err, client);
            return;
        }

        // Send the message to the user
        packets::Sendservercommand packet;
        packet.sessionId = sessionId;
        packet.packetId = 0; // n/a
        packet.command = formattedMsg;
        std::stringstream buffer;
        msgpack::pack(buffer, packet);
        utils::SendDataTo(sockfd, &buffer, packets::SENDSERVERCOMMAND, &(connection->second));
    }
}




// Invisible packet case, simply ignore.  We don't want the client to be able to send a generic packet...
void Server::GenericHandler(msgpack::object_handle * deserialized_data, sockaddr_in *client) {
    packets::Generic packet;
    // Wrap this in a try catch, so bad cast doesn't crash the whole server
    try {
        deserialized_data->get().convert(packet);
    } catch (const std::exception& e) {
        utils::SendErrorTo(sockfd, std::string("Failed to convert/cast msgpack object"), client);
        return;
    }

    std::lock_guard<std::mutex> lk(gameStateMutex);
    if(!gameState->VerifySession(packet.sessionId)){
        // Quit early if invalid session
        std::cout << "Invalid session. Not doing anything" << std::endl;
        return;
    }
}

// Respond to any packet that does not have an associated server action.  Those other packets will be acked by response.
void Server::AckHandler(msgpack::object_handle * deserialized_data, sockaddr_in *client) {
    packets::Ack packet;
    // Wrap this in a try catch, so bad cast doesn't crash the whole server
    try {
        deserialized_data->get().convert(packet);
    } catch (const std::exception& e) {
        utils::SendErrorTo(sockfd, std::string("Failed to convert/cast msgpack object"), client);
        return;
    }

    std::lock_guard<std::mutex> lk(gameStateMutex);
    if(!gameState->VerifySession(packet.sessionId)){
        // Quit early if invalid session
        std::cout << "Invalid session. Not doing anything" << std::endl;
        return;
    }
}

// Connect a host to the server by generating a session for it, and adding it to the gamestate sessions.
void Server::ConnectHandler(msgpack::object_handle * deserialized_data, sockaddr_in *client) {
    packets::Connect packet;
    // Wrap this in a try catch, so bad cast doesn't crash the whole server
    try {
        deserialized_data->get().convert(packet);
    } catch (const std::exception& e) {
        utils::SendErrorTo(sockfd, std::string("Failed to convert/cast msgpack object"), client);
        return;
    }

    // Check to make sure the given session exists
    if (!gameState->VerifySession(packet.sessionId)) {
        utils::SendErrorTo(sockfd, std::string("Failed to convert/cast msgpack object"), client);
        return;
    }

    // Save a copy of the connection information
    activeConnectionsMutex.lock();
    activeConnections[packet.sessionId] = *client;
    activeConnectionsMutex.unlock();
    // TODO: Figure out a way to disconnect clients if unresponsive

    packets::Connect returnPacket;
    returnPacket.sessionId = packet.sessionId;
    returnPacket.packetId = 0;
    // Use MessagePack to serialize data
    std::stringstream buffer;
    msgpack::pack(buffer, returnPacket);

    // Send the packet
    utils::SendDataTo(sockfd, &buffer, packets::CONNECT, client);

    // TODO: Sending gamestate to connected client

    std::cout << "\nNew connection started, session id " << returnPacket.sessionId << " sent to client!" << std::endl;
}

// Remove the session for a given user, effectively disconnecting it from the server.
void Server::DisconnectHandler(msgpack::object_handle * deserialized_data, sockaddr_in *client) {
    packets::Disconnect packet;
    // Wrap this in a try catch, so bad cast doesn't crash the whole server
    try {
        deserialized_data->get().convert(packet);
    } catch (const std::exception& e) {
        utils::SendErrorTo(sockfd, std::string("Failed to convert/cast msgpack object"), client);
        return;
    }

    std::lock_guard<std::mutex> lk(gameStateMutex);
    if(!gameState->VerifySession(packet.sessionId)){
        // Quit early if invalid session
        std::cout << "Invalid session. Not doing anything" << std::endl;
        return;
    }
    gameState->DisconnectSession(packet.sessionId);
}

void Server::ListCharactersHandler(msgpack::object_handle * deserialized_data, sockaddr_in *client) {
    packets::Listcharacters packet;
    // Wrap this in a try catch, so bad cast doesn't crash the whole server
    try {
        deserialized_data->get().convert(packet);
    } catch (const std::exception& e) {
        utils::SendErrorTo(sockfd, std::string("Failed to convert/cast msgpack object"), client);
        return;
    }

    std::lock_guard<std::mutex> lk(gameStateMutex);
    if(!gameState->VerifySession(packet.sessionId)){
        // Quit early if invalid session
        utils::SendErrorTo(sockfd, std::string("Invalid session"), client);
        return;
    }
    std::string account = gameState->GetSessionAccountName(packet.sessionId);
    std::cout << "Account: " << account << " requested their character list." << std::endl;

    // TODO: Need to get an actual list of the player's characters

    packets::Listcharacters returnPacket;
    returnPacket.sessionId = packet.sessionId;
    // TODO: return the same packetId that it came with?
    returnPacket.packetId = packet.packetId;
    // TODO: For testing only
        returnPacket.characters.push_back(std::string("Grapthar")); // Test
        returnPacket.characters.push_back(std::string("Harry Potter")); // Test
        returnPacket.characters.push_back(std::string("Willy Wonka")); // Test
    // Use MessagePack to serialize data
    std::stringstream buffer;
    msgpack::pack(buffer, returnPacket);
    // Send the packet
    utils::SendDataTo(sockfd, &buffer, packets::LISTCHARACTERS, client);
}

void Server::SelectCharacterHandler(msgpack::object_handle * deserialized_data, sockaddr_in *client) {
    packets::Selectcharacter packet;
    // Wrap this in a try catch, so bad cast doesn't crash the whole server
    try {
        deserialized_data->get().convert(packet);
    } catch (const std::exception& e) {
        utils::SendErrorTo(sockfd, std::string("Failed to convert/cast msgpack object"), client);
        return;
    }

    std::lock_guard<std::mutex> lk(gameStateMutex);
    if(!gameState->VerifySession(packet.sessionId)){
        utils::SendErrorTo(sockfd, std::string("Invalid session"), client);
        return;
    }
    std::cout << "TODO: Need to select character " << packet.character << std::endl;
    // TODO: Select the requested character
    gameState->SelectPlayer(packet.sessionId);
}

void Server::DeleteCharacterHandler(msgpack::object_handle * deserialized_data, sockaddr_in *client) {
    std::cout << "TODO: Delete Character handler!" << std::endl;
}

void Server::CreateCharacterHandler(msgpack::object_handle * deserialized_data, sockaddr_in *client) {
    packets::Createcharacter packet;
    // Wrap this in a try catch, so bad cast doesn't crash the whole server
    try {
        deserialized_data->get().convert(packet);
    } catch (const std::exception& e) {
        utils::SendErrorTo(sockfd, std::string("Failed to convert/cast msgpack object"), client);
        return;
    }

    std::lock_guard<std::mutex> lk(gameStateMutex);
    if(!gameState->VerifySession(packet.sessionId)){
        utils::SendErrorTo(sockfd, std::string("Invalid session"), client);
        return;
    }

    // TODO:
    std::cout << "TODO: Create character " << packet.lastName << ", " <<  packet.firstName << std::endl;
}

void Server::InitializeGameHandler(msgpack::object_handle * deserialized_data, sockaddr_in *client) {
    std::cout << "InitializeGameHandler!" << std::endl;
}

void Server::UpdatePcHandler(msgpack::object_handle * deserialized_data, sockaddr_in *client) {
    std::cout << "UpdatePcHandler!" << std::endl;
}

void Server::UpdateNpcHandler(msgpack::object_handle * deserialized_data, sockaddr_in *client) {
    std::cout << "UpdateNpcHandler!" << std::endl;
}

void Server::SendPlayerCommandHandler(msgpack::object_handle * deserialized_data, sockaddr_in *client) {
    packets::Sendplayercommand packet;
    try {
        deserialized_data->get().convert(packet);
    } catch (const std::exception& e) {
        utils::SendErrorTo(sockfd, std::string("Failed to convert/cast msgpack object"), client);
        return;
    }
    std::lock_guard<std::mutex> lk(gameStateMutex);
    if(!gameState->VerifySession(packet.sessionId)){
        utils::SendErrorTo(sockfd, std::string("Invalid session"), client);
        return;
    }

    // TODO: Get active session stuff working
    // if (gameState->VerifyActiveSession(packet.sessionId)) {
        gameState->PlayerCommand(packet.command, packet.sessionId);
    // }
    // else {
    //     std::cout << "Nonactive session requested to send a player command..." << packet.sessionId << std::endl;
    // }

    std::stringstream server_response_s;
    server_response_s << "Command successfully received";

    packets::Sendplayercommand returnPacket;
    returnPacket.sessionId = packet.sessionId;
    returnPacket.packetId = packet.packetId;
    returnPacket.command = server_response_s.str();

    std::stringstream buffer;
    msgpack::pack(buffer, returnPacket);
    utils::SendDataTo(sockfd, &buffer, packets::SENDPLAYERCOMMAND, client);
}

void Server::SendPlayerActionHandler(msgpack::object_handle * deserialized_data, sockaddr_in *client) {
    std::cout << "SendPlayerActionHandler!" << std::endl;
}

void Server::SendServerActionHandler(msgpack::object_handle * deserialized_data, sockaddr_in *client) {
    std::cout << "SendServerActionHandler!" << std::endl;
}


void Server::UnityHandler(msgpack::object_handle * deserialized_data, sockaddr_in *client) {
    std::cout << "UnityHandler!" << std::endl;
}
