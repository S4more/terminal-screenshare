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

// 3 -> Ctrl C
// 4 -> Ctrl D

func main() {
  // Investigate : https://stackoverflow.com/questions/72765557/using-a-pty-without-a-command
  addr := os.Args[1]
  tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	c := exec.Command("bash")
  _, err = pty.Start(c)

  if err != nil {
    panic(err)
  }

	if err != nil {
		panic(err)
	}

  conn, err := net.DialTCP("tcp", nil, tcpAddr)
  if err != nil {
    panic(err)
  }

  go listenToMessages(conn)

  oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
  if err != nil {
    panic(err)
  }
  defer func() { _ = term.Restore(int(os.Stdin.Fd()), oldState) }() // Best effort.
  // go func() { _, _ = io.Copy(terminalFile, os.Stdin) }()

  for {
    var buf bytes.Buffer
    io.CopyN(&buf, os.Stdin, 1)
    // fmt.Println(buf)
    sendMessage(conn, buf.Bytes())
  }
}

func sendMessage(conn net.Conn, data []byte) {
  conn.Write(data)
}

func listenToMessages(conn net.Conn) {
  for {
    reply := make([]byte, 1)
    _, err := conn.Read(reply)
    if err != nil {
      panic(err)
    }
    // fmt.Println(string(reply[:length]))
    os.Stdout.Write(reply)
  }
}
