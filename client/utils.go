package main

import "net"

const (
	PacketTypeClientDeclaration = "ClientDecl"
	PacketTypeClientRequest     = "ClientReq"
	PacketTypeServerResponse    = "ServerResp"
	PacketTypeUnknownUser       = "UnknownUser"
)

type NATPunchigPacket struct {
	Type   string
	Name   string
	Paylad []net.TCPAddr
}
