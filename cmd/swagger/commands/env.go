package commands

import (
	"fmt"
	"runtime"
)

type EnvCommand struct {
	BaseCommand
}

func (e *EnvCommand) Execute(_ []string) error {
	goos := runtime.GOOS
	compiler := runtime.Compiler
	goarch := runtime.GOARCH
	fmt.Fprintf(e.Out, "OS: %s\n", goos)
	fmt.Fprintf(e.Out, "Architecture: %s\n", goarch)
	fmt.Fprintf(e.Out, "Compiler: %s\n", compiler)
	return nil
}
