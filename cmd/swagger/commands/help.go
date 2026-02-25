package commands

import (
	"fmt"
	"io"
)

type HelpCommand struct {
	BaseCommand
}

func writeLines(w io.Writer, lines []string) error {
	for _, l := range lines {
		if _, err := fmt.Fprintln(w, l); err != nil {
			return err
		}
	}

	return nil
}

func (h *HelpCommand) Execute(_ []string) error {
	lines := []string{
		"Available commands:",
		"  help    Show this help message",
	}

	return writeLines(h.Out, lines)
}
