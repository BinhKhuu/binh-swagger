package commands

import (
	"fmt"
	"strings"
)

type EchoCommand struct {
	BaseCommand

	Upper bool     `description:"Uppercase output" long:"upper" short:"u"`
	Args  EchoArgs `positional-args:"yes"`
}

type EchoArgs struct {
	Text string `positional-arg-name:"text"`
}

func (e *EchoCommand) Execute(_ []string) error {
	text := e.Args.Text
	if e.Upper {
		text = strings.ToUpper(text)
	}

	fmt.Fprintf(e.Out, "%s\n", text)
	return nil
}
