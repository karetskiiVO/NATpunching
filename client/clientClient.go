package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func ClientMode(address string) error {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}
	defer conn.Close()

	go func() {
		netScanner := bufio.NewScanner(conn)
		for netScanner.Scan() {
			message := netScanner.Text()
			fmt.Println(message)
		}
	}()

	inputScanner := bufio.NewScanner(os.Stdin)
	for inputScanner.Scan() {
		text := inputScanner.Text()
		conn.Write([]byte(text))
	}

	return inputScanner.Err()
}
