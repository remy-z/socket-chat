let socket = new WebSocket("ws://localhost:3000/ws")
socket.onmessage = (event) => {
    console.log("Recieved from the server:", event.data)
}


var button = document.getElementById("sender");
button.addEventListener("click", function (event) {
    socket.send("Client says hello");
});