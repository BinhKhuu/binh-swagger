package main

import (
	"binh-swagger/cmd/swagger/commands"
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
		os.Exit(1)
	}

	_, err = parser.Parse()
	if err != nil {
		os.Exit(1)
	}
}
