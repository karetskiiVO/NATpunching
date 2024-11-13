package main

import (
	"encoding/json"
	"log"
	"net"
	"slices"

	"github.com/jessevdk/go-flags"
)

func main() {
	var options struct {
		Args struct {
			Port string
		} `positional-args:"yes" required:"1"`
	}

	parser := flags.NewParser(&options, flags.Default&(^flags.PrintErrors))
	_, err := parser.Parse()
	if err != nil {
		log.Fatal(err)
	}

	tcpListner, err := net.Listen("tcp", ":"+options.Args.Port)
	if err != nil {
		log.Fatal(err)
	}
	defer tcpListner.Close()

	clientsMapChan := make(chan map[string]([]net.TCPAddr))

	for {
		conn, err := tcpListner.Accept()
		if err != nil {
			log.Print(err)
			continue
		}

		go func() {
			decoder := json.NewDecoder(conn)
			encoder := json.NewEncoder(conn)

			for {
				var packet NATPunchigPacket
				err := decoder.Decode(&packet)
				if err != nil {
					log.Print(err)
					continue
				}

				log.Printf("from %v recieved: %v", conn.RemoteAddr().String(), packet)

				switch packet.Type {
				case PacketTypeClientDeclaration:
					clients := <-clientsMapChan

					clients[packet.Name] = make([]net.TCPAddr, 2)
					addr, _ := net.ResolveTCPAddr("tcp", conn.RemoteAddr().String())
					clients[packet.Name][0] = *addr
					clients[packet.Name][1] = packet.Paylad[0]

					clientsMapChan <- clients
				case PacketTypeClientRequest:
					clients := <-clientsMapChan

					response := NATPunchigPacket{}
					if info, ok := clients[packet.Name]; ok {
						response = NATPunchigPacket{
							Type:   PacketTypeServerResponse,
							Name:   packet.Name,
							Paylad: slices.Clone(info),
						}
					} else {
						response.Type = PacketTypeUnknownUser
					}

					clientsMapChan <- clients

					err = encoder.Encode(response)
					if err != nil {
						log.Print(err)
					}
				}
			}
		}()
	}

}
