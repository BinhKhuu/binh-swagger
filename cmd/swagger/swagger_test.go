package main

import (
	"reflect"
	"sort"
	"testing"

	"github.com/jessevdk/go-flags"
)

var ExpectedCommandNames = []string{
	"help",
	"version",
	"env",
	"echo",
	"config",
}

func Test_SwaggerCommandsRegistered_HasName(t *testing.T) {
	parser, err := makeParser()
	if err != nil {
		t.Fatalf("unexpected error when creating the parser")
	}

	commands := parser.Commands()

	existingNames := make(map[string]struct{})
	for _, cmd := range commands {
		existingNames[cmd.Name] = struct{}{}
	}

	if len(existingNames) != len(ExpectedCommandNames) {
		t.Fatalf("found %d commands expected %d commands", len(existingNames), len(ExpectedCommandNames))
	}

	for _, name := range ExpectedCommandNames {
		if _, hasCommand := existingNames[name]; !hasCommand {
			t.Fatalf("expected command %q to be registered", name)
		}
	}
}

// Duplicated test but i want to document how to compare and test something by name.
func Test_SwaggerCommandAllRegistered_ByName(t *testing.T) {
	parser, err := makeParser()
	if err != nil {
		t.Fatalf("unexpected error when creating the parser")
	}

	commands := parser.Commands()

	got := make([]string, 0, len(commands))
	for _, cmd := range commands {
		got = append(got, cmd.Name)
	}

	sort.Strings(got)
	sort.Strings(ExpectedCommandNames)

	if !reflect.DeepEqual(got, ExpectedCommandNames) {
		t.Fatalf("commands mismatch: got %v, want %v", got, ExpectedCommandNames)
	}
}

func makeParser() (*flags.Parser, error) {
	parser, err := BuildParser()
	return parser, err
}
