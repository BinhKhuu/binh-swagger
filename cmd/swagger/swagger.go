package main

import (
	"binh-swagger/cmd/swagger/commands"
	"log"
	"os"

	"github.com/jessevdk/go-flags"
)

func main() {
	parser, err := BuildParser()
	if err != nil {
		log.Fatalf("%v", err)
	}

	_, err = parser.Parse()
	if err != nil {
		os.Exit(1)
	}
}

func BuildParser() (*flags.Parser, error) {
	parser := flags.NewParser(nil, flags.Default)
	baseCommand := &commands.BaseCommand{Out: os.Stdout}

	_, err := parser.AddCommand(
		"help",
		"Show help",
		"Displays help information",
		&commands.HelpCommand{BaseCommand: *baseCommand},
	)
	if err != nil {
		return nil, err
	}
	_, err = parser.AddCommand(
		"version",
		"Shows version",
		"Prints Version information",
		&commands.VersionCommand{BaseCommand: *baseCommand},
	)
	if err != nil {
		return nil, err
	}
	_, err = parser.AddCommand(
		"env",
		"Shows version",
		"Prints Version information",
		&commands.EnvCommand{BaseCommand: *baseCommand},
	)
	if err != nil {
		return nil, err
	}
	_, err = parser.AddCommand(
		"echo",
		"repeats arguments",
		"Prints arguments",
		&commands.EchoCommand{BaseCommand: *baseCommand},
	)
	if err != nil {
		return nil, err
	}
	_, err = parser.AddCommand(
		"config",
		"loads config",
		"loads config and reads it",
		&commands.ConfigCommand{BaseCommand: *baseCommand},
	)
	if err != nil {
		return nil, err
	}

	return parser, nil
}
