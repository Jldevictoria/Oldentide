// Filename:    network.go
// Author:      Joseph DeVictoria
// Date:        July_11_2021
// Purpose:     Networking functions used by Oldentide dedicated server.

package shared

import (
	"fmt"
	"net"

	"github.com/vmihailenco/msgpack"
)

// MarshallAndSendPacket will take in a generic packet, marshall it, and send it to the specified UDP address.
func MarshallAndSendPacket(v interface{}, target net.Conn) {
	reqpac, err := msgpack.Marshal(v)
	CheckErr(err)
	fmt.Println(reqpac)
	fmt.Println(target)
	target.Write(reqpac)
}
