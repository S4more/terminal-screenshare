package main

import (
	"fmt"
	"net"
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

  _, err = conn.Write([]byte("Hello, World!"))
  if err != nil {
    println("Couldn't send message.")
    panic(err)
  }

  reply := make([]byte, 1024)
  _, err = conn.Read(reply)
  if err != nil {
    println("Read failed")
    panic(err)
  }

  fmt.Printf("Server reply: %s\n", string(reply))

  conn.Close()
}
