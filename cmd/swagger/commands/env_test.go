package commands

import (
	"bytes"
	"strings"
	"testing"
)

func Test_GoCommand(t *testing.T) {
	var buf bytes.Buffer

	baseCommand := SetupBaseCommand(&buf)
	cmd := &EnvCommand{BaseCommand: *baseCommand}

	err := cmd.Execute(nil)
	if err != nil {
		t.Fatalf("error executing command: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "OS") {
		t.Fatalf("expected OS name in output, got %q", out)
	}

	if !strings.Contains(out, "Architecture") {
		t.Fatalf("expected Architecture in output, got %q", out)
	}

	if !strings.Contains(out, "Compiler") {
		t.Fatalf("expected Compiler in output, got %q", out)
	}
}
