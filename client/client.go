package main

import (
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
)

type client struct {
}

func main() {
	var options struct {
		Args struct {
			Name       string
			ServerAddr string
		}
	}

	parser := flags.NewParser(&options, flags.Default&(^flags.PrintErrors))
	_, err := parser.Parse()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	//clients[]

}
