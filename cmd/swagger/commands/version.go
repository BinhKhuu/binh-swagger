package commands

import (
	"binh-swagger/cmd/swagger/commands/internal/meta"
	"fmt"
	"io"
)

type VersionCommand struct {
	Out io.Writer
}

func (v *VersionCommand) Execute(_ []string) error {
	_, err := fmt.Fprintf(v.Out, "%s\n", meta.Version)
	if err != nil {
		fmt.Fprintf(v.Out, "Error Printing version: %v", err)
	}
	return nil
}
