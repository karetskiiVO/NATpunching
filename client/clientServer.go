package main

import (
	"bufio"
	"log"
	"net"
)

func ServerMode(port string) error {
	tcpListner, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}
	defer tcpListner.Close()

	for {
		conn, err := tcpListner.Accept()
		if err != nil {
			log.Print(err)
			continue
		}

		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			message := scanner.Text()
			log.Println(message)
			conn.Write([]byte(message))
		}

		err = scanner.Err()
		if err != nil {
			log.Print(err)
			continue
		}
	}
}
