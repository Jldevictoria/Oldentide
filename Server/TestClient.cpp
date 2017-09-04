// Filename:    TestClient.cpp
// Author:      Joseph DeVictoria
// Date:        2_10_2013
// Purpose:     Act as an intermediate test platform for proving server functionality.

#include "Packets.h"
#include <arpa/inet.h>
#include <cstring>
#include <cstdlib>
#include <iostream>
#include <netdb.h>
#include <netinet/in.h>
#include <stdlib.h>
#include <string>
#include <sys/types.h>
#include <sys/socket.h>
#include <thread>
#include "Utils.h"
#include <msgpack.hpp>
#include <sstream>
#include <queue>
#include <mutex>
#include <condition_variable>


std::condition_variable isNotEmpty;
std::mutex packetQueueMutex;
std::queue<packets::packet_t> packetQueue; // Queue for all other packets in main thread
// std::mutex packetQueueSpecialMutex;
// std::queue<packets::packet_t> packetQueueSpecial; // Queue for server broadcasts

void PacketListener(int sockfd){
    sockaddr_in servret;
    uint8_t packetType;
    packets::packet_t p;

    while(1) {
        // Block and wait for a packet
        utils::ReceivePacketFrom(sockfd, p.data, &(p.source));

        msgpack::object_handle deserialized_data;
        uint8_t packetType;
        deserialized_data = utils::GetDataFromPacket(p.data, &packetType);

        // Do a quick packet check
        if(packetType == packets::SENDSERVERCOMMAND){
            // For now, skip the queue and service the server broadcast packet immediately
            packets::Sendservercommand data;
            deserialized_data.get().convert(data);
            std::cout << data.command << std::endl;

            // TODO: Break this out into it's own thread, so that it doesn't mess with getting the other packets
            // TODO: Put it in the special server broadcast packet queue
            // packetQueueSpecialMutex.lock();
            // packetQueueSpecial.push(p);
            // packetQueueSpecialMutex.unlock();
        }
        else {
            // Copy the packet into the packet queue
            packetQueueMutex.lock();
            packetQueue.push(p);
            packetQueueMutex.unlock();
            // Let the main thread know that there is a packet ready now
            isNotEmpty.notify_one();
        }
    }
}


// A helper function for waiting on a packet from the packet queue
// It blocks while waiting for the packetQueue to get filled up by
// the PacketListener thread.
// Yes, I know, it looks ugly, but it works for now
packets::packet_t WaitOnPacketQueue(){
    // The following 2 lines will lock, release lock, and block thread
    // until notify_one() is called after a packet arrives
    // http://en.cppreference.com/w/cpp/thread/condition_variable/wait
    std::unique_lock<std::mutex> ul(packetQueueMutex);
    isNotEmpty.wait(ul);
    // Wait until the packet queue has something in it

    // Now we own the lock again
    // Consume packet
    packets::packet_t retPacket = packetQueue.front();
    packetQueue.pop();
    packetQueueMutex.unlock();
    return retPacket;
}



int main(int argc, char * argv[]) {

    int sockfd;
    struct sockaddr_in servaddr;
    char * server_address;
    std::string session = "";
    int packetNumber = 1;
    bool running = true;

    // TODO: Parameter checking
    // Have parameter checking and exit gracefully if server address and port aren't specified
    if (argc != 3) {
        std::cout << "Invalid number of arguments passed to " << argv[0] << "; Exiting..." << std::endl;
        return 1;
    }

    // Read in server address.
    server_address = argv[1];
    std::cout << "Trying to connect to: " << server_address << std::endl;

    // For debug....
    struct hostent * he = gethostbyname(server_address);
    struct in_addr ** addr_list;
    addr_list = (struct in_addr **)he->h_addr_list;
    server_address = inet_ntoa(*addr_list[0]);

    int port = atoi(argv[2]);
    std::cout << server_address << std::endl;
    std::cout << port << std::endl;

    sockfd = socket(AF_INET, SOCK_DGRAM, IPPROTO_UDP);

    servaddr.sin_family = AF_INET;
    servaddr.sin_addr.s_addr = inet_addr(server_address);
    servaddr.sin_port = htons(port);

    // Spin off a thread that continually waits for packets and puts them in a queue
    std::thread t(PacketListener, sockfd);
    t.detach(); // Detach the threads so we don't need to join them manually

    int clientState = 0;

    while (running) {
        switch (clientState) {
            // Initial State.
            case 0: {
                std::cout << "Connect? (Y/n) " << std::endl;
                std::string response;
                getline (std::cin, response);
                if ((response.compare("y") == 0) || (response.compare("Y") == 0) || (response.compare("") == 0)) {
                    packets::Connect packet;
                    packet.packetId = packetNumber;
                    packetNumber++;
                    packet.sessionId = session;

                    // Use MessagePack to serialize data
                    std::stringstream buffer;
                    msgpack::pack(buffer, packet);
                    utils::SendDataTo(sockfd, &buffer, packets::CONNECT, &servaddr);

                    // Wait for a response
                    packets::packet_t retPacket = WaitOnPacketQueue();

                    std::cout << "Received a packet from " << utils::GetIpAndPortFromSocket(&retPacket.source) << std::endl;

                    uint8_t packetType;
                    msgpack::object_handle returnData = utils::GetDataFromPacket(retPacket.data, &packetType);
                    // We are assuming that the return packet will always be on type CONNECT
                    if(packetType != packets::CONNECT){
                        std::cout << "ERROR: Received packet other than CONNECT!! " << std::endl;
                        break;
                    }

                    // Use MessagePack to Deserialize the data
                    packets::Connect returnConnectPacket;
                    returnData.get().convert(returnConnectPacket);

                    std::cout << "Connected! Given the session id: " << returnConnectPacket.sessionId << std::endl;
                    // Set the session for the client
                    session = returnConnectPacket.sessionId;
                    clientState = 1;
                }
                else {
                    std::cout << "Shutting down!" << std::endl;
                    running = false;
                }
                break;
            }
            // Connected.
            case 1: {
                clientState = 2;
                std::cout << "Client State 2..." << std::endl;
                break;
            }
            // Logged In.
            case 2: {
                bool hasNoCharacter = true;

                while(hasNoCharacter) {
                    // Request a list of characters
                    packets::Listcharacters p;
                    p.packetId = packetNumber;
                    packetNumber++;
                    p.sessionId = session;

                    std::stringstream buffer;
                    msgpack::pack(buffer, p);
                    utils::SendDataTo(sockfd, &buffer, packets::LISTCHARACTERS, &servaddr);

                    // Wait for the response
                    packets::packet_t retPacket = WaitOnPacketQueue();
                    uint8_t packetType;
                    msgpack::object_handle deserialized = utils::GetDataFromPacket(retPacket.data, &packetType);
                    packets::Listcharacters characterList;
                    deserialized.get().convert(characterList);

                    std::cout << "Number of available characters: " << characterList.characters.size() << std::endl;
                    if (characterList.characters.size() == 0) {
                        std::cout << "You do not have any characters selected on this account" << std::endl;
                        std::cout << "Please give me a first name for your new character: ";
                        std::string firstName;
                        std::getline(std::cin, firstName);
                        std::cout << firstName << std::endl;
                        std::cout << "Please give me a last name for your new character: ";
                        std::string lastName;
                        std::getline(std::cin, lastName);

                        packets::Createcharacter newCharacter;
                        newCharacter.packetId = packetNumber;
                        packetNumber++;
                        newCharacter.sessionId = session;
                        newCharacter.firstName = firstName;
                        newCharacter.lastName = lastName;

                        std::stringstream buffer2;
                        msgpack::pack(buffer2, newCharacter);
                        utils::SendDataTo(sockfd, &buffer, packets::CREATECHARACTER, &servaddr);

                        // After creating the new character, loop back to the top
                        // and send a new request for the list of players.
                        // If the create was successful, then there should be a haracter
                    }
                    else {
                        hasNoCharacter = false;
                        for (int i = 0; i < characterList.characters.size(); ++i) {
                            if(characterList.characters.at(i) != ""){
                                std::cout << "Character " << i << ": " << characterList.characters.at(i) << std::endl;
                            }
                        }
                    }
                } // End has no character while loop

                // std::cout << "Please select a character: " << std::endl;
                // std::string name;
                // getline(std::cin, name);
                // std::cout << "TODO: Attempting to select character " << name << ": " << std::endl;

                // packets::Selectcharacter characterToSelect;
                // characterToSelect.packetId = packetNumber;
                // packetNumber++;
                // characterToSelect.sessionId = session;
                // characterToSelect.character = name;

                // std::stringstream buffer;
                // msgpack::pack(buffer, characterToSelect);
                // utils::SendDataTo(sockfd, &buffer, packets::SELECTCHARACTER, &servaddr);

                clientState = 3;
                std::cout << "Now, send commands to the server!" << std::endl;
                break;
            }
            // In game...
            case 3: {

                std::string command;
                getline(std::cin, command);
                if (command.empty()) {
                    break;
                }

                // if (utils::Tokenfy(command, ' ')[0] == "/stress") {
                // if (utils::Tokenfy(command, ' ')[0] == "/s") {
                // }

                packets::Sendplayercommand playerCommand;
                playerCommand.packetId = packetNumber;
                packetNumber++;
                playerCommand.sessionId = session;
                // playerCommand.command = std::string("/s " + command);
                playerCommand.command = std::string(command);

                std::stringstream buffer;
                msgpack::pack(buffer, playerCommand);
                utils::SendDataTo(sockfd, &buffer, packets::SENDPLAYERCOMMAND, &servaddr);


                packets::packet_t retPacket = WaitOnPacketQueue();
                uint8_t packetType;
                msgpack::object_handle deserialized = utils::GetDataFromPacket(retPacket.data, &packetType);
                packets::Sendplayercommand playerCommandResponse;
                deserialized.get().convert(playerCommandResponse);
                std::cout << playerCommandResponse.command << std::endl;

                break;
            }
        }
    }
}
