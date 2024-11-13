package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"

	"github.com/jessevdk/go-flags"
)

func main() {
	var options struct {
		Args struct {
			Name       string
			Mode       string
			Port       string
			ServerAddr string

			OptionalNextAddress []string
		}
	}

	parser := flags.NewParser(&options, flags.Default&(^flags.PrintErrors))
	_, err := parser.Parse()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// регистририруемся на сервере
	conn, err := net.Dial("tcp", options.Args.ServerAddr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	encoder := json.NewEncoder(conn)
	decoder := json.NewDecoder(conn)

	currTCPArddr, err := net.ResolveTCPAddr("tcp", conn.LocalAddr().String())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	encoder.Encode(NATPunchigPacket{
		Name: options.Args.Name,
		Type: PacketTypeClientDeclaration,
		Paylad: []net.TCPAddr{
			*currTCPArddr,
		},
	})

	var response NATPunchigPacket
	err = decoder.Decode(&response)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	destaddr := ""

	if (len(options.Args.OptionalNextAddress)) == 0 {
		for {
			var destname string
			fmt.Printf("Input destination login: ")
			fmt.Scan(&destname)

			encoder.Encode(NATPunchigPacket{
				Name: destname,
				Type: PacketTypeClientRequest,
			})

			err = decoder.Decode(&response)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			if response.Type == PacketTypeUnknownUser {
				fmt.Println("Unknown login")
			} else {
				destaddr = response.Paylad[0].String()
				break
			}
		}
	}

	conn.Close()

	if (len(options.Args.OptionalNextAddress)) == 0 {

		err = ServerMode(options.Args.Port)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		err = ClientMode(destaddr)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}
