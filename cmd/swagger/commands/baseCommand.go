package commands

import (
	"bytes"
	"io"
)

type BaseCommand struct {
	Out io.Writer
}

// SetupBaseCommand initializes and returns a new BaseCommand with the provided buffer.
func SetupBaseCommand(buff *bytes.Buffer) *BaseCommand {
	return &BaseCommand{
		Out: buff,
	}
}
