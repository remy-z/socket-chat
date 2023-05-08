# socket-chat
Testing out websockets in go with a simple chat application

In one terminal host go ws server: 
    go run main.go

In another terminal host simple http with python:
    navigate to /frontend
    python -m http.server

ngrok http 8000