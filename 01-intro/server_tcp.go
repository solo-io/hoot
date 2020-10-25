package main

import (
	"io"
	"log"
	"net"
)

func echoServer(conn net.Conn) {
	defer conn.Close()
	io.Copy(conn, conn)
}

func main() {
	l, err := net.Listen("tcp", ":8082")
	if err != nil {
		log.Fatal("listen error:", err)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Panicln(err)
		}

		go echoServer(conn)
	}

}
