package main

import (
	"log"

	"github.com/karetskiiVO/NATpunching/natpunch"
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

	server, err := natpunch.NewMeetupServer(options.Args.Port)
	if err != nil {
		log.Fatal(err)
	}
	defer server.Finalize()

	server.Run()
}
