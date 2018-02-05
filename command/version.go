package command

import (
	"strings"

	"github.com/mitchellh/cli"
)

type VersionCommand struct {
	Ui cli.Ui
}

func (this *VersionCommand) Run(args []string) int {
	this.Ui.Info("0.0.1")
	return 0
}

func (this *VersionCommand) Help() string {
	helpText := `
Usage: cli version

	Prints the current version
`
	return strings.TrimSpace(helpText)
}

func (this *VersionCommand) Synopsis() string {
	return "Prints the current version"
}
