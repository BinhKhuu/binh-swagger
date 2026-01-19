package commands

import (
	"binh-swagger/cmd/swagger/commands/internal/meta"
	"fmt"
)

type VersionCommand struct {
	BaseCommand
}

func (v *VersionCommand) Execute(_ []string) error {
	_, err := fmt.Fprintf(v.Out, "%s\n", meta.Version)
	if err != nil {
		fmt.Fprintf(v.Out, "Error Printing version: %v", err)
	}
	return nil
}
