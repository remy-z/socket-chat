package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
	"log"

	"golang.org/x/net/websocket"
)

type Server struct {
	conns map[*websocket.Conn]bool
}

func NewServer() *Server {
	return &Server{
		conns: make(map[*websocket.Conn]bool),
	}
}

func (s *Server) exampleFeed(ws *websocket.Conn) {
	fmt.Println("new connection from client to feed:", ws.RemoteAddr())

	for {
		payload := fmt.Sprintf("Just saying hi to all my clients <3 at -> %d\n", time.Now().UnixNano())
		ws.Write([]byte(payload))
		time.Sleep(time.Second * 2)
	}
}

func (s *Server) handleWS(ws *websocket.Conn) {
	fmt.Println("new connection from client:", ws.RemoteAddr())

	s.conns[ws] = true // maps in golang not concurrent safe should use mutex

	s.readLoop(ws)
}

func (s *Server) readLoop(ws *websocket.Conn) {
	buf := make([]byte, 1024)
	for {
		n, err := ws.Read(buf)
		if err != nil {
			if err == io.EOF {
				//EOF means connection on other side has closed
				break
			}
			fmt.Println("read error:", err)
			continue // returning would drop the connection
		}
		msg := buf[:n]
		s.broadcast(msg)
	}

}

func (s *Server) broadcast(b []byte) {
	for ws := range s.conns {
		go func(ws *websocket.Conn) {
			if _, err := ws.Write(b); err != nil {
				fmt.Println("error", err)
			}
		}(ws)
	}
}

func allowCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		h.ServeHTTP(w, r)
	})
}

func main() {
	server := NewServer()
	http.Handle("/ws", allowCORS(websocket.Handler(server.handleWS)))
	http.Handle("/examplefeed", allowCORS(websocket.Handler(server.exampleFeed)))
	log.Fatal(http.ListenAndServe(":3000", nil))
}
