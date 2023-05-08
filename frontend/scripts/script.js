let socket = new WebSocket("ws://localhost:3000/ws")
socket.onmessage = (event) => {
    const el = document.createElement('li');
    el.innerHTML = event.data;
    document.querySelector('ul').appendChild(el)
}

document.querySelector('button').onclick = () => {

    const text = document.querySelector('input').value;
    socket.send(text)

}

