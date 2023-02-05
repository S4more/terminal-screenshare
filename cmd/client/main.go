package main

import (
	"bytes"
	"io"
	"net"
	"os"
	"os/exec"

	"github.com/creack/pty"
	"golang.org/x/term"
)

func main() {
  addr := "localhost:8080"
  tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
  if err != nil {
    panic(err)
  }

  conn, err := net.DialTCP("tcp", nil, tcpAddr)
  if err != nil {
    panic(err)
  }

  ch := make(chan []byte)
  go listenToMessages(conn)
  go startTerminal(ch)

  for {
    msg := <- ch
    if len(msg) == 0 {
      return
    }
    sendMessage(conn, msg)
  }
}

func startTerminal(ch chan []byte) {
	c := exec.Command("zsh")
	terminalFile, err := pty.Start(c)
	if err != nil {
		panic(err)
	}

  defer func() { _ = terminalFile.Close() }() // Best effort.

  // Set stdin in raw mode.
  oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
  if err != nil {
          panic(err)
  }
  defer func() { _ = term.Restore(int(os.Stdin.Fd()), oldState) }() // Best effort.
  go func() { _, _ = io.Copy(terminalFile, os.Stdin) }()

  for {
    var buf bytes.Buffer
    writer := io.MultiWriter(os.Stdout, &buf)
    io.CopyN(writer, terminalFile, 1)
    ch <- buf.Bytes()
  }
}

func sendMessage(conn net.Conn, data []byte) {
  conn.Write(data)
}

func listenToMessages(conn net.Conn) {
  for {
    reply := make([]byte, 1024)
    _, err := conn.Read(reply)
    if err != nil {
      println("Read failed")
      panic(err)
    }
    // fmt.Printf("Server echo: %s\n", string(reply))
  }
}
