package main

import (
  "fmt"
  "net"
)

func main() {
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

  fmt.Printf("Listneng on %s:%s\n", host, port)

  for {
    conn, err := l.Accept()
    if err != nil {
      panic(err)
    }

    go onConnection(conn)
  }
}

func onConnection(conn net.Conn) {
  buf := make([]byte, 1024)
  lenght, err := conn.Read(buf)
  if err != nil {
    fmt.Printf("Error reading: %#v\n", err)
    return
  }
  fmt.Printf("Message received: %s\n", string(buf[:lenght]))

  conn.Write([]byte("Message received\n"))
  conn.Close()
}
