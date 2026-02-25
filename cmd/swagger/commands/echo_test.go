package commands

import (
	"bytes"
	"strings"
	"testing"
)

var text = "ugh"

func Test_EchoCommand_UpperTrue(t *testing.T) {
	cmd, buff := setupEchoTests()
	cmd.Upper = true

	err := cmd.Execute(nil)
	if err != nil {
		t.Fatalf("Error executing echo command: %v", err)
	}

	output := strings.TrimSpace(buff.String())

	expectedOutput := strings.TrimSpace(strings.ToUpper(text))
	if output != expectedOutput {
		t.Fatalf("Error exspect %s but got %s", strings.ToUpper(text), output)
	}
}

func Test_EchoCommand_UpperFalse(t *testing.T) {
	cmd, buff := setupEchoTests()
	cmd.Upper = false

	err := cmd.Execute(nil)
	if err != nil {
		t.Fatalf("Error executing echo command: %v", err)
	}

	output := strings.TrimSpace(buff.String())

	expectedOutput := strings.TrimSpace(text)
	if output != expectedOutput {
		t.Fatalf("Error exspect %s but got %s", text, output)
	}
}

func setupEchoTests() (*EchoCommand, *bytes.Buffer) {
	var buff bytes.Buffer

	args := EchoArgs{Text: text}
	baseCommand := SetupBaseCommand(&buff)
	cmd := &EchoCommand{
		BaseCommand: *baseCommand,
		Upper:       true,
		Args:        args,
	}

	return cmd, &buff
}
