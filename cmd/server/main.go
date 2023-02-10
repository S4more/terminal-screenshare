package main

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"

	"github.com/creack/pty"
	"golang.org/x/term"
)

type Server struct {
  connections []net.Conn
  terminalFile os.File
}

func newServer() *Server {
  p := Server {connections: make([]net.Conn, 0)}
  return &p
}

func main() {
  s := newServer()
  s.listen()
}

func (s *Server) listen() {
  addr := "localhost:8080"
  l, err := net.Listen("tcp", addr)
  if err != nil {
    panic(err)
  }
  defer l.Close()

  host, port, err := net.SplitHostPort(l.Addr().String())
  if err != nil {
    panic(err)
  }

  fmt.Printf("Listening on %s:%s\n", host, port)

  ch := make(chan []byte)
  go s.startTerminal(ch)

  go func() {
    for {
      msg := <- ch
      // if len(msg) == 0 {
      //   return
      // }
      for _, conn:= range s.connections {
        sendMessage(conn, msg)
      }
    }
  }()

  for {
    conn, err := l.Accept()
    if err != nil {
      panic(err)
    }

    go s.onConnection(conn)
  }
}


func (s *Server) onConnection(conn net.Conn) {
  s.connections = append(s.connections, conn)
  for {
    buf := make([]byte, 1)
    lenght, err := conn.Read(buf)
    if err != nil {
      fmt.Printf("Error reading: %#v\n", err)
      return
    }
    s.terminalFile.Write(buf[:lenght])
  }
}

func (s *Server) startTerminal(ch chan []byte) {
	c := exec.Command("bash")
	terminalFile, err := pty.Start(c)
	if err != nil {
		panic(err)
	}

  s.terminalFile = *terminalFile

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
