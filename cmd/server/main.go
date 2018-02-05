package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/kr/pty"
	"golang.org/x/crypto/ssh/terminal"
	"golang.org/x/net/websocket"
)

func EchoServer(ws *websocket.Conn) {
	// fmt.Fprint(os.Stdout, "Client ===>")
	// multiWriter := io.MultiWriter(os.Stdout, ws)
	// io.Copy(multiWriter, ws)

	// ==========================

	// Start the command with a pty.
	c := exec.Command("bash")

	// Start the command with a pty.
	ptmx, err := pty.Start(c)
	if err != nil {
		log.Fatal(err.Error())
	}
	// Make sure to close the pty at the end.
	defer func() { _ = ptmx.Close() }() // Best effort.

	// Handle pty size.
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGWINCH)
	go func() {
		for range ch {
			if err := pty.InheritSize(os.Stdin, ptmx); err != nil {
				log.Printf("error resizing pty: %s", err)
			}
		}
	}()
	ch <- syscall.SIGWINCH // Initial resize.

	// Set stdin in raw mode.
	oldState, err := terminal.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer func() { _ = terminal.Restore(int(os.Stdin.Fd()), oldState) }() // Best effort.

	// teeReader := io.TeeReader(ws, os.Stdout)
	multiWriter := io.MultiWriter(ws, os.Stdout)

	go func() { _, _ = io.Copy(ptmx, ws) }() //reading from Stdin, writing to tty
	_, _ = io.Copy(multiWriter, ptmx)        //reading from tty, writing to Stdout

	fmt.Println("=====> Exiting EchoHandler")
}

func main() {
	http.Handle("/", websocket.Handler(EchoServer))
	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
