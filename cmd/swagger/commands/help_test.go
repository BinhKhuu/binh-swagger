package commands

import (
	"bytes"
	"strings"
	"testing"
)

func Test_HelpPrintsMessage(t *testing.T) {
	var buf bytes.Buffer
	cmd := &HelpCommand{Out: &buf}
	cmd.Execute(nil)

	if !strings.Contains(buf.String(), "Available commands") {
		t.Fatalf("expected help output")
	}
}
