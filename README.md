# socket-chat
Testing out websockets in go with a simple chat application

test in devtools console 

let socket = new WebSocket("ws://localhost:3000/ws")

socket.onmessage = (event) => {console.log("Recieved from the server:", event.data)}

socket.send("Client says hello")
