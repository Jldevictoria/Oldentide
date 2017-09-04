// Filename:    utils.cpp
// Author:      Joseph DeVictoria
// Date:        Mar_10_2016
// Purpose:     Miscellaneous functions.

#include "Utils.h"
#include "Packets.h"
#include <algorithm>
#include <regex>
#include <iostream>
#include <iterator>
#include <string.h>

#include <netinet/in.h>
#include <sys/socket.h>
#include <sys/types.h>

#include <netdb.h> // For getnameinfo() and related flags
#include <msgpack.hpp>

namespace utils {

    const int MIN_ACCOUNT_NAME_LENGTH = 3;
    const int MAX_ACCOUNT_NAME_LENGTH = 30;
    const int MIN_PASSWORD_LENGTH = 8;
    const int MAX_PASSWORD_LENGTH = 30;

    // Simple function that separates a string (s) into tokens based on a delimiter (delim).
    std::vector<std::string> Tokenfy(std::string s, char delim) {
        size_t start=0;
        size_t end=s.find_first_of(delim);
        std::vector<std::string> tokens;
        while (end <= std::string::npos) {
	        tokens.emplace_back(s.substr(start, end-start));
	        if (end == std::string::npos)
	    	    break;
    	    start=end+1;
    	    end = s.find_first_of(delim, start);
        }
        return tokens;
    }

    // CLIENT.
    bool CheckPasswordLength(char *password) {
        // Check to make sure password is reasonable
        if (strlen(password) > MAX_PASSWORD_LENGTH) {
            printf("Password is too large!\n");
            return 0;
        }
        else if (strlen(password) < MIN_PASSWORD_LENGTH) {
            printf("Password is too small!\n" );
            return 0;
        }
        else {
            // Password is of a good length
            return 1;
        }
    }

    // CLIENT/SERVER.
    bool CheckAccountNameLength(char *account_name) {
        // Check to make sure account_name is reasonable
        if (strlen(account_name) > MAX_ACCOUNT_NAME_LENGTH) {
            printf("Account name is too large!\n");
            return 0;
        }
        else if (strlen(account_name) < MIN_ACCOUNT_NAME_LENGTH) {
            printf("Account name is too small!\n" );
            return 0;
        }
        else {
            // Account name is of good length
            return 1;
        }
    }

    // CLIENT/SERVER. Checks to make sure the passed string looks like a valid account name
    // Checks the length of the string, and then makes sure it is alphanumeric and underscore only.
    bool SanitizeAccountName(char *account_name) {
        if (!utils::CheckAccountNameLength(account_name)) {
            return false;
        }
        if (!utils::SanitizeAlphanumeric(account_name)) {
            printf("Invalid account name! Account name must be only contain characters a-z, A-Z, 0-9, _\n");
            return false;
        }
        else {
            return true;
        }
    }

    // CLIENT/SERVER. Checks to makesure the passed string is only alphanumeric characters and underscore.
    // Use regex to check that the account name is only alpha-numeric
    // Regex: \^\w{3,30}$\
    // Tested regex at regex101.com using javascript (ECMAScript) flavor
    // NOTE: Account names are stored in a case-insensitive way in the db
    // NOTE: In C, adjacent string literals are concatenated (MACRO
    // TODO: Wow. Adding regex stuff adds a full three seconds to the compile time... What's the deal?
    bool SanitizeAlphanumeric(char *string) {
        std::stringstream ss;
        ss << "^\\w{" << MIN_ACCOUNT_NAME_LENGTH << "," << MAX_ACCOUNT_NAME_LENGTH << "}$";
        std::string regex_string();
        std::regex check_alpha_regex( ss.str() );
        if (!regex_match(string, check_alpha_regex)) {
            return false;
        }
        else {
            return true;
        }
    }

    // CLIENT/SERVER. Checks to make sure the passed string is only hex - a-f, A-F, 0-9.
    // Returns 1 if it's hex, 0 if not.
    bool SanitizeHexString(char *string) {
        // Check to see if the passed string is at least one character (+) of only hex characters
        std::regex check_hex_regex("^[a-fA-F0-9]+$");
        if (!regex_match(string, check_hex_regex)) {
            return false;
        }
        else {
            return true;
        }
    }


    // Prints the bytes of str in a hex format
    void PrintStringHex(std::string *str){
        printf("0x");
        const char *c_str = str->data();
        size_t size = str->size();
        // Don't calculate length based on c_str, in case data has 0x00s fake out strlen(c_str)
        for (size_t i = 0; i < size; ++i) {
            printf("%02hhX", c_str[i]);
        }
        printf("\n");
    }


    /*

    Oldentide Packet Header
    =======================

    The following packet header shall be appended to the beginning
    of every packet sent to and from the server:

    one byte:
    +--------+
    |        |
    +--------+

    a variable number of bytes:
    +========+
    |        |
    +========+

    type is a uint8_t
    size is a uint16_t (stored little-endian)
    l = lower byte (LSB)
    h = higher/upper byte (MSB)

    +----------+----------+----------+=========================+
    |   type   |  size(l) |  size(h) |  msgpck_data...         |
    +----------+----------+----------+=========================+

    To safely and consistently access this data, use the helper functions in Utils.cpp.



    How to install and set up msgpck-c in the project:
    ==================================================

    See https://github.com/msgpack/msgpack-c#building-and-installing

        git clone https://github.com/msgpack/msgpack-c.git
        cd msgpack-c
        cmake .
        make
        sudo make install

    Note that sudo make install should copy the msgpck-c headers to /usr/local/include,
    and this should automatically be part of the include path to g++

    */







    // Adds in the packetType and the size of the msgpack data
    // TODO: Make this more efficient by using memcpy instead of insert?
    void PrependPacketHeader(std::string *str, uint8_t packetType){
        uint16_t size = str->size();
        if(size > packets::MSGPCK_MAX_PAYLOAD_SIZE){
            std::cout << "ERROR: MessagePack data is larger than " << packets::MSGPCK_MAX_PAYLOAD_SIZE << " bytes!!!!" << std::endl;
        }

        std::string size_str = std::string((char *)&size,2);
        // Add in header info - prepend the packet type and messagepack data size
        str->insert(0, size_str);
        str->insert(0, 1, (char)packetType);
    }

    // Returns the PTYPE of the packet
    uint8_t GetPacketTypeFromPacket(char *packet_buffer){
        // Simply return the first byte of the packet, since that is where the PTYPE value is
        return packet_buffer[0];
    }

    // Returns the size of the msgpack data payload
    // Accesses the udp packet from position 1, grabs 2 bytes, and stores them in a short
    uint16_t GetMsgpckSizeFromPacket(char *packet_buffer){
        // TODO: Make it so that there are no unaligned accesses -
        // i.e. no short access starting at arr[1] (RISC machines - ARM, etc)
        uint16_t *shorty_array = (uint16_t *)(&packet_buffer[1]);
        if(shorty_array[0] > packets::MSGPCK_MAX_PAYLOAD_SIZE){
            std::cout << "ERROR: Msgpck Size is too large (> " << packets::MSGPCK_MAX_PAYLOAD_SIZE << ")to be correct!" << std::endl;
        }

        return shorty_array[0];
    }



    std::string GetMsgpckDataFromPacket(char *packetBuffer){
        uint16_t msgpckSize = utils::GetMsgpckSizeFromPacket(packetBuffer);

        // Check to make sure we don't accidentally access data outside the packet!
        if(msgpckSize > packets::MSGPCK_MAX_PAYLOAD_SIZE){
            std::cout << "Msgpck Size is too large to be correct! Returning empty string" << std::endl;
            return std::string();
        }

        // Return the msgpck payload, which starts at index 3 and is msgpckSize in length
        return std::string(packetBuffer+packets::PACKET_HEADER_SIZE, msgpckSize);
    };





    // Returns the ip address and port number of the socket
    std::string GetIpAndPortFromSocket(sockaddr_in *socket){
        char host[100];
        char service[100];
        getnameinfo((sockaddr *)socket, sizeof(sockaddr), host, sizeof(host), service, sizeof(service), NI_NUMERICHOST);
        std::stringstream ss;
        ss << host << ":" << ntohs(socket->sin_port);
        return ss.str();
    }


    // Wrapper for recvfrom that returns msgpack data (object handle).
    // The packetType is returned in the uint8_t pointer, and sender info
    // is returned in the sockaddr_in pointer.
    // If packetType is 0, then an error occured.
    msgpack::object_handle ReceiveDataFrom(int sockfd, uint8_t *packetTypeOut, sockaddr_in *sourceOut){
        // Create a temporary internal buffer
        char packet[packets::PACKET_MAX_SIZE];
        // Wait for and get the packet and source of the packet
        ReceivePacketFrom(sockfd, packet, sourceOut);
        // Extract the packetType and messagepack data
        return GetDataFromPacket(packet, packetTypeOut);
    }



    // Receives a single packet
    // packetBufferOut should point to PACKET_MAX_SIZE allocated bytes
    void ReceivePacketFrom(int sockfd, char *packetBufferOut, sockaddr_in *sourceOut){
        static socklen_t LEN = sizeof(sockaddr_in);
        int n = recvfrom(sockfd, packetBufferOut, packets::PACKET_MAX_SIZE, 0, (struct sockaddr *)sourceOut, &LEN);
    }



    // packetBufferIn must be of length PACKET_MAX_SIZE
    msgpack::object_handle GetDataFromPacket(char *packetBufferIn, uint8_t *packetTypeOut){
        *packetTypeOut = utils::GetPacketTypeFromPacket(packetBufferIn);
        uint16_t msgpckSize = utils::GetMsgpckSizeFromPacket(packetBufferIn);

        // std::cout << "Packet type: " << (unsigned int) *packetTypeOut << std::endl;
        // std::cout << "Msgpack Size: " << msgpckSize << std::endl;

        // Check to make sure we don't accidentally access data outside the packet!
        if(msgpckSize > packets::MSGPCK_MAX_PAYLOAD_SIZE){
            std::cout << "Msgpack Size is too large to be correct! Ignoring packet and continuing to avoid buffer overflow..." << std::endl;
            *packetTypeOut = 0;
        }

        // Get msgpack data
        std::string msgpack_data = utils::GetMsgpckDataFromPacket(packetBufferIn);
        // std::cout << "Msgpack data: ";
        // utils::PrintStringHex(&msgpack_data);

        // Use MessagePack to deserialize the payload, and return it in dataOut
        msgpack::object_handle dataOut = msgpack::unpack(msgpack_data.data(), msgpack_data.size());

        // std::cout << "Msgpack deserialized: ";
        // std::cout << dataOut.get() << std::endl;

        return dataOut;
    }
    // NOTE: The reason object_handle is being returned instead of object is because
    // if object_handle goes out of scope, then the object itself will be freed,
    // and strange things will happen (e.g. if stringstream is used before convert, convert will fail)
    // Also, the = operator doesn't seem to work with object_handle, so I need to return it instead of
    // using a pointer.
    // See https://github.com/msgpack/msgpack-c/wiki/v2_0_cpp_unpacker





    // Wrapper for sendto. Takes data msgpacked into dataIn, prepends a custom header to include the
    // packet type, and sends data in a UDP packet to destIn.
    // Returns the return value of sendto.
    int SendDataTo(int sockfd, std::stringstream *dataIn, uint8_t packetTypeIn, sockaddr_in *destIn){
        // std::cout << "Sending packetType " << (int) packetTypeIn << std::endl;
        static socklen_t LEN = sizeof(sockaddr_in);
        // Make sure that the buffer is at the 0th position
        dataIn->seekg(0);
        std::string str(dataIn->str());
        // std::cout << "Sending msgpack data:" << std::endl;
        // utils::PrintStringHex(&str);
        // Add in header info - prepend the packet type and messagepack data size
        utils::PrependPacketHeader(&str, packetTypeIn);
        // std::cout << "Sending header + msgpack data:" << std::endl;
        // utils::PrintStringHex(&str);
        // Send the packet
        return sendto(sockfd, str.data(), str.size(), 0, (struct sockaddr *)destIn, LEN);
    }

    // TODO: Pass in string pointer? Will it make any difference?
    void SendErrorTo(int sockfd, std::string errorMsg, sockaddr_in *destIn){
        // Return error packet
        packets::Error returnPacket;
        // Since we couldn't convert/cast msgpack data, we don't know what sessionId or packetId it had
        returnPacket.sessionId = "";
        returnPacket.packetId = 0;
        returnPacket.errorMsg = errorMsg;
        std::stringstream buffer;
        msgpack::pack(buffer, returnPacket);
        utils::SendDataTo(sockfd, &buffer, packets::ERROR, destIn);
    }



}
