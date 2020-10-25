package main

import (
	"io"
	"log"
	"math/rand"
	"net"
	"net/http"
	_ "net/http/pprof"
)

func echoServer(rw http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if rand.Int()%5 == 0 {
		http.Error(rw, "error", http.StatusBadRequest)
		return
	}
	io.Copy(rw, r.Body)
}

func main() {
	l, err := net.Listen("tcp", ":8082")
	if err != nil {
		log.Fatal("listen error:", err)
	}

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	http.Serve(l, http.HandlerFunc(echoServer))
}
