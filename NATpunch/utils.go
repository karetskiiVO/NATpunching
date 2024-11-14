package natpunch

import (
	"net"
)

type NATPunchigPacket struct {
	Name       string
	LocalAddr  *net.UDPAddr
	GlobalAddr *net.UDPAddr
}

var EmptyAddress, _ = net.ResolveUDPAddr("udp", "0.0.0.0:0")
