package natpunch

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strings"
)

type MeetupServer struct {
	conn    net.PacketConn
	clients map[string]NATPunchigPacket
}

func NewMeetupServer(port string) (*MeetupServer, error) {
	conn, err := net.ListenPacket("udp", ":"+port)

	if err != nil {
		log.Fatal(err)
	}
	return &MeetupServer{
		conn:    conn,
		clients: make(map[string]NATPunchigPacket),
	}, nil
}

func (serv *MeetupServer) Finalize() error {
	return serv.conn.Close()
}

func (serv *MeetupServer) Run() error {
	for {
		var packet NATPunchigPacket
		inputBuffer := make([]byte, 2048)

		n, addr, err := serv.conn.ReadFrom(inputBuffer)
		if err != nil {
			log.Print(err)
			continue
		}

		err = json.Unmarshal(inputBuffer[:n], &packet)
		if err != nil {
			log.Print(err)
			continue
		}

		log.Printf("from %v: %v", addr.String(), packet)

		// проверить бы
		if packet.LocalAddr.String() == EmptyAddress.String() {
			err = serv.resolveRequest(addr, packet)
			if err != nil {
				log.Print(err)
			}
		} else {
			err = serv.resolveRegistration(addr, packet)
			if err != nil {
				log.Print(err)
			}
		}
	}
}

func (serv *MeetupServer) resolveRegistration(addr net.Addr, packet NATPunchigPacket) error {
	packet.GlobalAddr = addr.(*net.UDPAddr)
	serv.clients[packet.Name] = packet

	log.Printf("Client[%v] registration at %v, local addres %v", packet.Name, packet.GlobalAddr, packet.LocalAddr)
	outputBuffer, err := json.Marshal(packet)
	if err != nil {
		return err
	}

	_, err = serv.conn.WriteTo(outputBuffer, addr)
	if err != nil {
		return err
	}

	return nil
}

func (serv *MeetupServer) resolveRequest(addr net.Addr, packet NATPunchigPacket) error {
	names := strings.Split(packet.Name, "@")
	if len(names) != 2 {
		return fmt.Errorf("unrecognized name: %v", packet.Name)
	}

	targetName := names[0]
	searcherName := names[1]

	var targetPacket, searcherPacket NATPunchigPacket
	ok := false
	if searcherPacket, ok = serv.clients[searcherName]; !ok {
		packet.LocalAddr = EmptyAddress
		packet.GlobalAddr = EmptyAddress

		outputBuffer, err := json.Marshal(packet)
		if err != nil {
			return err
		}

		_, err = serv.conn.WriteTo(outputBuffer, addr)
		if err != nil {
			return err
		}

		return fmt.Errorf("searcher %v unregistered", searcherName)
	}
	if targetPacket, ok = serv.clients[targetName]; !ok {
		packet.LocalAddr = EmptyAddress
		packet.GlobalAddr = EmptyAddress

		outputBuffer, err := json.Marshal(packet)
		if err != nil {
			return err
		}

		_, err = serv.conn.WriteTo(outputBuffer, addr)
		if err != nil {
			return err
		}

		return fmt.Errorf("target %v unregistered", searcherName)
	}

	outputBuffer, err := json.Marshal(targetPacket)
	if err != nil {
		return err
	}

	_, err = serv.conn.WriteTo(outputBuffer, searcherPacket.GlobalAddr)
	if err != nil {
		return err
	}

	outputBuffer, err = json.Marshal(searcherPacket)
	if err != nil {
		return err
	}

	_, err = serv.conn.WriteTo(outputBuffer, targetPacket.GlobalAddr)
	if err != nil {
		return err
	}

	return nil
}
