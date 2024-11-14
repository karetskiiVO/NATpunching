package natpunch

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"sync"
	"time"
)

type Client struct {
	name   string
	conn   net.PacketConn
	server net.Addr

	clientsMu  sync.RWMutex
	clients    map[string]*net.UDPAddr
	selfPacket NATPunchigPacket
}

func NewClient(name, port, serverAddr string) (*Client, error) {
	resport := ":0"
	if port != "" {
		resport = ":"+port 
	}

	ipres, err := net.Dial("udp", serverAddr)
	if err != nil {
		return nil, err
	}
	ipstr := ipres.LocalAddr().(*net.UDPAddr).IP.String()
	ipres.Close()

	conn, err := net.ListenPacket("udp", ipstr+resport)

	fmt.Println(conn.LocalAddr().String())

	if err != nil {
		return nil, err
	}

	server, err := net.ResolveUDPAddr("udp", serverAddr)
	if err != nil {
		return nil, err
	}

	return &Client{
		name:    name,
		conn:    conn,
		server:  server,
		clients: make(map[string]*net.UDPAddr),
	}, nil
}

func (cl *Client) Registrate() (bool, error) {
	packet := NATPunchigPacket{
		Name:       cl.name,
		GlobalAddr: EmptyAddress,
	}
	packet.LocalAddr = cl.conn.LocalAddr().(*net.UDPAddr)

	msg, err := json.Marshal(packet)
	if err != nil {
		return false, err
	}
	_, err = cl.conn.WriteTo(msg, cl.server)
	if err != nil {
		return false, err
	}

	reply := make([]byte, 1024)
	_, addr, err := cl.conn.ReadFrom(reply)
	if err != nil {
		return false, err
	}
	if addr.String() != cl.conn.LocalAddr().String() {
		return false, nil
	}

	err = json.Unmarshal(reply, &packet)
	if err != nil {
		return false, err
	}

	cl.selfPacket = packet

	return true, nil
}

func (cl *Client) Finalize() error {
	return cl.conn.Close()
}

func (cl *Client) Run() error {
	scanner := bufio.NewScanner(os.Stdin)

	go cl.demonListner()

	for scanner.Scan() {
		buffer := scanner.Bytes()

		delim := bytes.Index(buffer, []byte{':'})
		if delim < 0 {
			fmt.Println("the user's name is not recognized")
			continue
		}

		destname := string(buffer[:delim])

		var destAddr *net.UDPAddr
		cl.clientsMu.RLock()
		destAddr, ok := cl.clients[destname]
		cl.clientsMu.Unlock()

		if !ok {
			var err error
			destAddr, err = cl.punch(destname)
			if err != nil {
				fmt.Println(err)
				continue
			}
		}

		_, err := cl.conn.WriteTo(append([]byte(cl.name), buffer[delim+1:]...), destAddr)
		if err != nil {
			fmt.Println(err)
			continue
		}
	}

	return scanner.Err()
}

func (cl *Client) demonListner() {
	serverAddrString := cl.server.String()

	buf := make([]byte, 1024)

	for {
		n, addr, err := cl.conn.ReadFrom(buf)
		if err != nil {
			fmt.Println(err)
			continue
		}

		var packet NATPunchigPacket
		err = json.Unmarshal(buf[:n], &packet)
		switch {
		case addr.String() == serverAddrString:
			if err != nil {
				fmt.Println(err)
				continue
			}

			reply, _ := json.Marshal(cl.selfPacket)
			cl.conn.WriteTo(reply, packet.LocalAddr)
			cl.conn.WriteTo(reply, packet.GlobalAddr)
		case err == nil:
			cl.clientsMu.RLock()
			_, ok := cl.clients[packet.Name]
			cl.clientsMu.RUnlock()

			if ok {
				break
			}

			cl.clientsMu.Lock()
			cl.clients[packet.Name] = addr.(*net.UDPAddr)
			cl.clientsMu.Unlock()

			reply, _ := json.Marshal(cl.selfPacket)
			cl.conn.WriteTo(reply, addr)
		default:
			fmt.Printf("[%v] as %v", addr.String(), string(buf[:n]))
		}
	}
}

func (cl *Client) punch(username string) (*net.UDPAddr, error) {
	packet := NATPunchigPacket{
		Name:       username + "@" + cl.name,
		GlobalAddr: EmptyAddress,
		LocalAddr:  EmptyAddress,
	}

	punchStarter, _ := json.Marshal(packet)
	for i := 0; i < 10; i++ {
		cl.conn.WriteTo(punchStarter, cl.server)

		time.Sleep(20 * time.Millisecond)
		cl.clientsMu.RLock()
		addr, ok := cl.clients[username]
		cl.clientsMu.RUnlock()

		if ok {
			return addr, nil
		}
	}

	return nil, fmt.Errorf("can't find %v", username)
}
