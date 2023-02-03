package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
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

  sendMessage(conn)
  // go listenToMessages(conn)

  for {}

  // conn.Close()
}

func sendMessage(conn net.Conn) {
  for {
    reader := bufio.NewReader(os.Stdin)
    fmt.Print("-> ")
    text, _ := reader.ReadString('\n')
    text = strings.Replace(text, "\n", "", -1)
    conn.Write([]byte(text))
  }
}

func listenToMessages(conn net.Conn) {
  for {
    reply := make([]byte, 1024)
    _, err := conn.Read(reply)
    if err != nil {
      println("Read failed")
      panic(err)
    }
    fmt.Printf("Server echo: %s\n", string(reply))
  }
}
