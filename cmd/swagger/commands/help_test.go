package commands

import (
	"bytes"
	"strings"
	"testing"
)

func Test_HelpPrintsMessage(t *testing.T) {
	var buf bytes.Buffer
	cmd := &HelpCommand{Out: &buf}
	err := cmd.Execute(nil)
	if err != nil {
		t.Fatalf("error executing command: %v", err)
	}
	if !strings.Contains(buf.String(), "Available commands") {
		t.Fatalf("expected help output")
	}
}
