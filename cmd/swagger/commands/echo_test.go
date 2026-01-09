package commands

import (
	"bytes"
	"testing"
)

var text = "ugh"

func Test_EchoCommand_UpperTrue(t *testing.T) {
	cmd, buff := setupTests()
	cmd.Upper = true
	err := cmd.Execute(nil)
	if err != nil {
		t.Fatalf("Error executing echo command: %v", err)
	}
	output := buff.String()
	if output != "UGH" {
		t.Fatalf("Error exspect UGH but got %s", output)
	}
}

func Test_EchoCommand_UpperFalse(t *testing.T) {
	cmd, buff := setupTests()
	cmd.Upper = false
	err := cmd.Execute(nil)
	if err != nil {
		t.Fatalf("Error executing echo command: %v", err)
	}
	output := buff.String()
	if output != text {
		t.Fatalf("Error exspect %s but got %s", text, output)
	}
}

func setupTests() (*EchoCommand, *bytes.Buffer) {
	var buff bytes.Buffer
	args := EchoArgs{Text: text}
	cmd := &EchoCommand{
		Out:   &buff,
		Upper: true,
		Args:  args,
	}

	return cmd, &buff
}
