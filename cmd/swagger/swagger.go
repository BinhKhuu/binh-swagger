package main

import (
	"binh-swagger/cmd/swagger/commands"
	"log"
	"os"

	"github.com/jessevdk/go-flags"
)

func main() {
	parser := flags.NewParser(nil, flags.Default)

	_, err := parser.AddCommand(
		"help",
		"Show help",
		"Displays help information",
		&commands.HelpCommand{Out: os.Stdout},
	)
	if err != nil {
		log.Fatalf("%v", err)
	}
	_, err = parser.AddCommand(
		"version",
		"Shows version",
		"Prints Version information",
		&commands.VersionCommand{Out: os.Stdout},
	)
	if err != nil {
		log.Fatalf("%v", err)
	}
	_, err = parser.Parse()
	if err != nil {
		os.Exit(1)
	}
}
