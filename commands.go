package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/jcarley/cli/command"
	"github.com/mitchellh/cli"
)

// Commands is the mapping of all the available Serf commands.
var Commands map[string]cli.CommandFactory

func init() {

	ui := &cli.ColoredUi{
		InfoColor:  cli.UiColorGreen,
		ErrorColor: cli.UiColorRed,
		Ui:         &cli.BasicUi{Writer: os.Stdout, Reader: os.Stdin, ErrorWriter: os.Stderr},
	}

	prefixedUi := &cli.PrefixedUi{
		InfoPrefix: "==> ",
		Ui:         ui,
	}

	Commands = map[string]cli.CommandFactory{

		"server": func() (cli.Command, error) {
			return &command.ServerCommand{
				Ui: prefixedUi,
			}, nil
		},

		"version": func() (cli.Command, error) {
			return &command.VersionCommand{
				Ui: ui,
			}, nil
		},
	}
}

// makeShutdownCh returns a channel that can be used for shutdown
// notifications for commands. This channel will send a message for every
// interrupt received.
func makeShutdownCh() <-chan struct{} {
	resultCh := make(chan struct{})

	signalCh := make(chan os.Signal, 4)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		for {
			<-signalCh
			resultCh <- struct{}{}
		}
	}()

	return resultCh
}
