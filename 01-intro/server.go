package main

import (
	"io"
	"log"
	"net"
	"net/http"
)

func echoServer(rw http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	io.Copy(rw, r.Body)
}

func main() {
	l, err := net.Listen("tcp", ":8082")
	if err != nil {
		log.Fatal("listen error:", err)
	}

	http.Serve(l, http.HandlerFunc(echoServer))
}
