package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
)

func redServer(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(rw, "red")
}
func blueServer(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusAccepted)
	fmt.Fprintln(rw, "blue")
}

func main() {
	l, err := net.Listen("tcp", ":8082")
	if err != nil {
		log.Fatal("listen error:", err)
	}
	l2, err := net.Listen("tcp", ":8083")
	if err != nil {
		log.Fatal("listen error:", err)
	}
	go http.Serve(l, http.HandlerFunc(redServer))
	go http.Serve(l2, http.HandlerFunc(blueServer))

	select {}
}
