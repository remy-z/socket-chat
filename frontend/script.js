var selectedchat = "general";
var username = "guest";

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
    var date = new Date().toLocaleString()

    messageContainer = document.getElementById("message-container");
    message = messageDivBuilder(messageEvent.from, messageEvent.message, date);
    messageContainer.appendChild(message);
}

// returns a div with class message to be appended to the message container
function messageDivBuilder(name, message, date) {
    var div = document.createElement('div');
    div.className = 'message';

    var p1 = document.createElement('p');
    p1.className = 'msg message-body';

    var span = document.createElement('span');
    span.className = 'msg nameplate';
    span.textContent = `${name}: `;

    var messageText = document.createTextNode(message)

    p1.appendChild(span);
    p1.appendChild(messageText);

    var p2 = document.createElement('p');
    p2.className = 'msg timestamp';
    p2.textContent = date;


    div.appendChild(p1);
    div.appendChild(p2);

    return div
}

function sendEvent(eventName, payload) {
    const event = new Event(eventName, payload)

    socket.send(JSON.stringify(event))
}

function sendMessage() {
    var message = document.getElementById("message-input");
    if (message != null) {
        //TODO allow username login 
        let outgoingEvent = new SendMessageEvent(message.value, username);
        sendEvent("send_message", outgoingEvent)
    }
    return false;
}

function login() {
    let formData = {
        "username": document.getElementById("username").value,
        "password": document.getElementById("password").value
    }
    username = formData["username"]

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
        loadChatPage()
    }).catch((e) => { alert(e) });


    return false;
}

function connectWebSocket(otp) {
    if (window["WebSocket"]) {
        console.log("supports websockets");

        socket = new WebSocket("wss://" + document.location.host + "/ws?otp=" + otp)

        socket.onopen = function (event) {
            console.log("socket-opened")
        }

        socket.onclose = function (event) {
            console.log("socket-closed")
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
    //document.getElementById("chatroom-selection").onsubmit = changeChatRoom;
    document.getElementById("login-form").onsubmit = login;
    document.getElementById("send-message").onsubmit = sendMessage;

}

function loadChatPage() {
    // hide login and display chat items
    var elements = document.querySelectorAll(".login, .chat");
    for (var i = 0; i < elements.length; i++) {
        elements[i].classList.toggle("hide");
    }
    //load messages 
    //focus on input
}
