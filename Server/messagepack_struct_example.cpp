// https://github.com/msgpack/msgpack-c.git

/* To Compile:
    g++ -std=c++11 -I /home/michael/Documents/msgpack-c/include/ messagepack_struct_example.cpp
*/

// Simply clone the code and set the include path

#include <msgpack.hpp>
#include <string>
#include <iostream>
#include <sstream>

#include "Packets.h"

// Test struct
typedef struct {
    PTYPE packetType;
    int32_t packetId;
    std::string sessionId;
    std::string my_str;
    // Need to designate the values for messagepack
    MSGPACK_DEFINE(packetType, packetId, sessionId, my_str);
    // MSGPACK_DEFINE(packetType, packetId);
} foo_struct_t;
// Need to specify enums as a type
MSGPACK_ADD_ENUM(PTYPE);


// Copied from the message pack c++ example
int main(int argc, char const *argv[]) {
    // Can also try using tuples
    std::tuple<uint8_t, uint8_t, uint8_t> src1(
        15,
        31,
        255
    );

    // Cannot do this unless you take away the default values for struct type
    foo_struct_t src2 = {
        // Use designated/tagged initialization
        .packetType = SENDPLAYERACTION,
        .packetId = 17,
        .sessionId = 1337,
        .my_str = std::string("Hello!")
    };

    unsigned int src3 = 1024;

    // serialize the object into the buffer.
    // any classes that implements write(const char*,size_t) can be a buffer.
    std::stringstream buffer;
    // msgpack::pack(buffer, src1);
    msgpack::pack(buffer, src2);
    // msgpack::pack(buffer, src3);

    buffer.seekg(0);
    std::string str(buffer.str());

    // Print out how the message pack looks
    const char *msgpack_string = str.data();
    // Don't calculate length based on c_str, in case data has 0x00s fake out strlen(c_str)
    for (int i = 0; i < str.length(); ++i) {
        // Use hh to designate char-width, since it expects an int
        // See http://www.cplusplus.com/reference/cstdio/printf/
        // Use 02 to backfill with 0s up to 2 places
        // See http://stackoverflow.com/questions/8441257/why-does-c-print-my-hex-values-incorrectly
        printf("%02hhX", msgpack_string[i]);
    }
    printf("\n");

    // // Print out each char individually
    // for (int i = 0; i < str.length(); ++i) {
    //     printf("%d (%p): %02hhX\n", i, msgpack_string+i, msgpack_string[i]);
    // }

    // It appears that messagepack converts a struct into a fixarray

    // To manually decode, look at the spec:
    // https://github.com/msgpack/msgpack/blob/master/spec.md

    // For src2 struct, the messagepack hex string is:
    // 940F11CD0539A648656C6C6F21
    // 94 0F 11 CD0539 A6 48656C6C6F21

    // 94               -> fixarray with 4 elements
    // 0F               -> positive fixnum : 15
    // 11               -> positive fixnum : 17
    // CD 0539          -> uint 16 big endian : 1337
    // A6 48656C6C6F21  -> fixstr : Hello!

    // Messagepack doesn't rely on null termination, because it already knows how long the data strings are


    msgpack::object_handle oh = msgpack::unpack(str.data(), str.size());

    // deserialized object is valid during the msgpack::object_handle instance is alive.
    msgpack::object deserialized = oh.get();

    // msgpack::object supports ostream.
    std::cout << deserialized << std::endl;

    // convert msgpack::object instance into the original type.
    // if the type is mismatched, it throws msgpack::type_error exception.

    // std::tuple<uint8_t, uint8_t, uint8_t> dst1;
    // deserialized.convert(dst1);
    // // Use + to force uint8_t (unsigned char) to print as a number, not an ascii char
    // // See http://stackoverflow.com/questions/14644716/how-to-output-a-character-as-an-integer-through-cout
    // std::cout << +std::get<0>(dst1) << "|" << +std::get<1>(dst1) << "|" << +std::get<2>(dst1) << std::endl;

    foo_struct_t dst2;
    deserialized.convert(dst2);
    // Prove that I can access the struct like normal
    std::cout << dst2.packetType << "|" << dst2.packetId << "|" << dst2.sessionId << "|" << dst2.my_str << std::endl;

    // unsigned int dst3;
    // deserialized.convert(dst3);
    // std::cout << dst3 << std::endl;


    return 0;
}
