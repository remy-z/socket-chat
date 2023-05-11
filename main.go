package main

import (
	"context"
	"log"
	"net/http"
)

func main() {
	setupServer()

	log.Fatal(http.ListenAndServeTLS(":8080", "server.crt", "server.key", nil))
	//log.Fatal(http.ListenAndServe(":8080", nil))
}

func setupServer() {

	ctx := context.Background()

	manager := NewManager(ctx)
	http.Handle("/", http.FileServer(http.Dir("./frontend")))
	http.HandleFunc("/ws", manager.serveWS)
	http.HandleFunc("/login", manager.loginHandler)
}
