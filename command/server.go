package command

import (
	"strings"

	"github.com/mitchellh/cli"
)

type ServerCommand struct {
	Ui cli.Ui
}

func (this *ServerCommand) Run(args []string) int {
	return 0
}

func (this *ServerCommand) Help() string {
	helpText := `
Usage: cli server

	Starts up the cli server process.
`
	return strings.TrimSpace(helpText)
}

func (this *ServerCommand) Synopsis() string {
	return "Starts up the cli server process."
}
