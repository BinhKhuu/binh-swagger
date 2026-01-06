package commands

import (
	"fmt"
	"io"
)

type HelpCommand struct {
	Out io.Writer
}

func (h *HelpCommand) Execute(_ []string) error {
	fmt.Fprintln(h.Out, "Available commands:")
	fmt.Fprintln(h.Out, "  help    Show this help message")
	return nil
}
