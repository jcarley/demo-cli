package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/mitchellh/panicwrap"
	"golang.org/x/crypto/ssh/terminal"
	"golang.org/x/net/websocket"
)

func main() {
	os.Exit(realMain())
}

func realMain() int {
	exitStatus, err := panicwrap.BasicWrap(panicHandler)
	if err != nil {
		// Something went wrong setting up the panic wrapper. Unlikely,
		// but possible.
		panic(err)
	}

	// If exitStatus >= 0, then we're the parent process and the panicwrap
	// re-executed ourselves and completed. Just exit with the proper status.
	if exitStatus >= 0 {
		os.Exit(exitStatus)
	}

	// Otherwise, exitStatus < 0 means we're the child. Continue executing as
	// normal...

	// Create arbitrary command.
	// c := exec.Command("bash")

	// Start the command with a pty.
	// ptmx, err := pty.Start(c)
	// if err != nil {
	//   return err
	// }
	// Make sure to close the pty at the end.
	// defer func() { _ = ptmx.Close() }() // Best effort.

	// Handle pty size.
	// ch := make(chan os.Signal, 1)
	// signal.Notify(ch, syscall.SIGWINCH)
	// go func() {
	//   for range ch {
	//     if err := pty.InheritSize(os.Stdin, ptmx); err != nil {
	//       log.Printf("error resizing pty: %s", err)
	//     }
	//   }
	// }()
	// ch <- syscall.SIGWINCH // Initial resize.

	// Set stdin in raw mode.
	oldState, err := terminal.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer func() { _ = terminal.Restore(int(os.Stdin.Fd()), oldState) }() // Best effort.

	ws := dial()

	// Copy stdin to the pty and the pty to stdout.
	// logFile, err := os.OpenFile("access.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	// if err != nil {
	//   fmt.Printf("Error: %s", err.Error())
	//   os.Exit(1)
	// }
	// defer logFile.Close()

	// multiWriter := io.MultiWriter(os.Stdout, ws)
	// go func() { _, _ = io.Copy(ptmx, os.Stdin) }() //reading from Stdin, writing to tty
	// _, _ = io.Copy(multiWriter, ptmx)              //reading from tty, writing to Stdout

	// multiWriter := io.MultiWriter(ptmx, ws)
	// multiReader := io.MultiReader(ptmx, ws)

	go func() { _, _ = io.Copy(ws, os.Stdin) }() //reading from Stdin, writing to tty
	_, _ = io.Copy(os.Stdout, ws)                //reading from tty, writing to Stdout

	return nil
}

func panicHandler(output string) {
	// output contains the full output (including stack traces) of the
	// panic. Put it in a file or something.
	fmt.Printf("The child panicked:\n\n%s\n", output)
	os.Exit(1)
}

func dial() *websocket.Conn {
	origin := "http://localhost/"
	url := "ws://localhost:12345/ws"
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Fatal(err)
	}

	return ws
}
