package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"time"
)

func echoServer(rw http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	fmt.Fprintln(rw, "received request from", r.Header.Get("x-forwarded-for"))
}

func main() {
	l, err := net.Listen("tcp", ":8082")
	if err != nil {
		log.Fatal("listen error:", err)
	}

	go http.Serve(l, http.HandlerFunc(echoServer))

	srv := &http.Server{
		ReadTimeout:  time.Minute,
		WriteTimeout: time.Minute,
		Addr:         ":8083",
		Handler:      http.HandlerFunc(echoServer),
	}
	log.Fatal(srv.ListenAndServeTLS("example_com_cert.pem", "example_com_key.pem"))

}
