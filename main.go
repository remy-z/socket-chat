package main

import (
	"context"
	"log"
	"net/http"
	"os"
)

func main() {
	setupServer()

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	log.Fatal(http.ListenAndServeTLS(":"+port, "server.crt", "server.key", nil))
	//log.Fatal(http.ListenAndServe(":8080", nil))
}

func setupServer() {

	ctx := context.Background()

	manager := NewManager(ctx)
	http.Handle("/", http.FileServer(http.Dir("./frontend")))
	http.HandleFunc("/ws", manager.serveWS)
	http.HandleFunc("/login", manager.loginHandler)
}
