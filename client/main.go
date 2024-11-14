package main

import (
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/karetskiiVO/NATpunching/natpunch"
)

func main() {
	var options struct {
		Args struct {
			Name       string
			ServerAddr string
		} `positional-args:"yes" required:"1"`
		Port string `short:"p" description:"custom port" default:""`
	}

	parser := flags.NewParser(&options, flags.Default&(^flags.PrintErrors))
	_, err := parser.Parse()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	client, err := natpunch.NewClient(options.Args.Name, options.Port, options.Args.ServerAddr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer client.Finalize()

	_, err = client.Registrate()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = client.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
