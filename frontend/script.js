
var selectedchat = "general";

// what we will be sending and recieveing trhough websocket
class Event {
    constructor(type, payload) {
        this.type = type;
        this.payload = payload;
    }
}

class SendMessageEvent {
    constructor(message, from) {
        this.message = message;
        this.from = from;
    }
}

class RecieveMessageEvent {
    constructor(message, from, sent) {
        this.message = message;
        this.from = from;
        this.sent = sent;
    }
}

class ChangeRoomEvent {
    constructor(name) {
        this.name = name;
    }
}

function changeChatRoom() {
    // Change Header to reflect the Changed chatroom
    var newchat = document.getElementById("chatroom");
    if (newchat != null && newchat.value != selectedchat) {
        selectedchat = newchat.value;
        header = document.getElementById("chat-header").innerHTML = "Currently in chat: " + selectedchat;

        let changeEvent = new ChangeRoomEvent(selectedchat);
        sendEvent("change_room", changeEvent);
        textarea = document.getElementById("chatmessages");
        textarea.innerHTML = `You changed room into: ${selectedchat}`;
    }
    return false;
}

function routeEvent(event) {
    if (event.type === undefined) {
        alert('no type field in the event');
    }

    switch (event.type) {
        case "recieve_message":
            const messageEvent = Object.assign(new RecieveMessageEvent, event.payload);
            appendChatMessage(messageEvent)
            break;
        default:
            alert("unsupported message type");
            break;

    }
}

function appendChatMessage(messageEvent) {
    var date = new Date()
    console.log(date)
    const formattedMsg = `${date.toLocaleString()}: ${messageEvent.message}`

    textarea = document.getElementById('chatmessages')
    textarea.innerHTML = textarea.innerHTML + "\n" + formattedMsg;
    textarea.scrollTop = textarea.scrollHeight;
}

function sendEvent(eventName, payload) {
    const event = new Event(eventName, payload)

    socket.send(JSON.stringify(event))
}

function sendMessage() {
    var message = document.getElementById("message");
    if (message != null) {
        //TODO allow username login 
        let outgoingEvent = new SendMessageEvent(message.value, "test");
        sendEvent("send_message", outgoingEvent)
    }
    return false;
}

function login() {
    let formData = {
        "username": document.getElementById("username").value,
        "password": document.getElementById("password").value
    }

    fetch("login", {
        method: 'post',
        body: JSON.stringify(formData),
        mode: 'cors'
    }).then((response) => {
        if (response.ok) {
            return response.json();
        } else {
            throw "unauthorized";
        }
    }).then((data) => {
        //we are authenticated
        connectWebSocket(data.otp)
    }).catch((e) => { alert(e) });

    return false;
}

function connectWebSocket(otp) {
    if (window["WebSocket"]) {
        console.log("supports websockets");

        socket = new WebSocket("wss://" + document.location.host + "/ws?otp=" + otp)

        socket.onopen = function (event) {
            document.getElementById("connection-header").innerHTML = "Connected to WebSocket: true";
        }

        socket.onclose = function (event) {
            document.getElementById("connection-header").innerHTML = "Connected to WebSocket: false";
            // add automatic recconection unless it was manually closed by server
        }

        socket.onmessage = function (event) {

            const eventData = JSON.parse(event.data);

            const eventObject = Object.assign(new Event, eventData);

            routeEvent(eventObject);
        }
    } else {
        alert("browser does not support websockets")
    }
}

window.onload = function () {
    document.getElementById("chatroom-selection").onsubmit = changeChatRoom;
    document.getElementById("chatroom-message").onsubmit = sendMessage;
    document.getElementById("login-form").onsubmit = login;

}
