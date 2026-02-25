package commands

import (
	"binh-swagger/cmd/swagger/commands/internal/meta"
	"bytes"
	"strings"
	"testing"
)

func Test_VersionCommand(t *testing.T) {
	var outBuf bytes.Buffer

	baseCommand := SetupBaseCommand(&outBuf)
	cmd := &VersionCommand{BaseCommand: *baseCommand}

	err := cmd.Execute(nil)
	if err != nil {
		t.Fatalf("error executing command: %v", err)
	}

	if strings.TrimSpace(outBuf.String()) != meta.Version {
		t.Fatalf("expected version output: %s", meta.Version)
	}
}
