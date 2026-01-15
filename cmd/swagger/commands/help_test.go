package commands

import (
	"bytes"
	"strings"
	"testing"
)

func Test_HelpPrintsMessage(t *testing.T) {
	var buf bytes.Buffer
	baseCommand := SetupBaseCommand(&buf)
	cmd := &HelpCommand{BaseCommand: *baseCommand}
	err := cmd.Execute(nil)
	if err != nil {
		t.Fatalf("error executing command: %v", err)
	}
	if !strings.Contains(buf.String(), "Available commands") {
		t.Fatalf("expected help output \"Available commands\" but got: %s", buf.String())
	}
}
