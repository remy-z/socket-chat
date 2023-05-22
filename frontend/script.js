var selectedchat = "general";
var username = "guest";
var messageContainer = document.getElementById("message-container")

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


class SendDirCommand {
    //user and chat are bool values 
    constructor(chat, user) {
        this.chat = chat;
        this.user = user;
    }
}

class RecieveDirCommand {
    constructor(lines) {
        this.lines = lines
    }
}

function changeChatRoom(newchat) {
    if (newchat != null && newchat.value != selectedchat) {
        selectedchat = newchat;
        document.getElementById("chatroom").innerHTML = selectedchat;
        let changeEvent = new ChangeRoomEvent(selectedchat);
        sendEvent("change_room", changeEvent);
        while (messageContainer.firstChild) {
            messageContainer.removeChild(messageContainer.firstChild)
        }
        displaySystemMessage(`Successfully changed rooms. Welcome to the ${selectedchat} chat`);
    }
}

function routeEvent(event) {
    if (event.type === undefined) {
        alert('no type field in the event');
    }

    switch (event.type) {
        case "recieve_message":
            const messageEvent = Object.assign(new RecieveMessageEvent, event.payload);
            displayChatMessage(messageEvent)
            break;
        default:
            alert("unsupported message type");
            break;

    }
}

function displayChatMessage(messageEvent) {
    var date = new Date().toLocaleString()

    message = messageBuilder(messageEvent.from, messageEvent.message, date);
    messageContainer.appendChild(message);
    messageContainer.scrollTop = messageContainer.scrollHeight;
}

// returns a div with class message to be appended to the message container
function messageBuilder(from, message, date) {
    var div = document.createElement('div');
    div.className = 'message';

    var p1 = document.createElement('p');
    p1.className = 'msg message-body';

    var span = document.createElement('span');
    span.className = 'msg nameplate';
    span.textContent = `${from}: `;

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

function displaySystemMessage(message) {
    var systemMessage = systemMessageBuilder(message)
    messageContainer.appendChild(systemMessage)
}

function systemMessageBuilder(message) {
    var div = document.createElement('div');
    div.className = 'message';

    var p1 = document.createElement('p');
    p1.className = 'msg message-body';

    var span = document.createElement('span');
    span.className = 'msg nameplate';

    var messageText = document.createTextNode(message)

    p1.appendChild(span);
    p1.appendChild(messageText);

    div.appendChild(p1);

    return div
}

function sendEvent(eventName, payload) {
    const event = new Event(eventName, payload)

    socket.send(JSON.stringify(event))
}

function parseMessage() {
    input = document.getElementById("message-input");
    var message = input.value
    input.value = ""
    if (message === "") {
        return false
    }
    if (message.startsWith("/")) {
        console.log("looks like a command")
        args = message.substring(1).split(" ")
        console.log(args)
        switch (args[0]) {
            case "cd":
                console.log("cool directory")
                return cdCommand(args)
                break;
            case "dir":
                console.log("directing")
                return dirCommand(args)
            case "help":
                console.log("helping")
                return helpCommand(args)
                break;
            default:
                return commandError(args)
        }
    } else {
        sendMessage(message)
    }

    return false
}


function commandError(message) {
    console.log(`${message}: command not recognized. /help for a list of available commands`)
}

function helpCommand(args) {
    if (args.length > 1) {
        switch (args[i]) {
            case "cd":
                //give info on cd command
                break;
            case "dir":
                //give info on dir command
                break;
            case "help":
                //give info on help command
                break;
            default:
                commandError(args[i])
        }
    }
}

function cdCommand(args) {
    if (args.length > 2) {
        displaySystemMessage("Invalid number of arguments. For info on a command, try /help [command-name]")
        return
    }
    chatroom = args[1].toLowerCase()
    changeChatRoom(chatroom)
}

function dirCommand(args) {

}

function sendMessage(message) {
    if (message != null) {
        let outgoingEvent = new SendMessageEvent(message, username);
        sendEvent("send_message", outgoingEvent)
    }
}

function login() {
    let formData = {
        "username": document.getElementById("username").value,
        "password": document.getElementById("password").value
    }


    //add check to see if username already in use
    username = formData["username"] == "" ? "guest" : formData["username"]

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
    document.getElementById("login-form").onsubmit = login;
    document.getElementById("username").focus()
}

function loadChatPage() {
    // hide login and display chat items
    var elements = document.querySelectorAll(".login, .chat");
    for (var i = 0; i < elements.length; i++) {
        elements[i].classList.toggle("hide");
    }
    //load messages 
    displaySystemMessage("Welcome to Socket-Chat!")

    //focus on input
    input = document.getElementById("message-input")
    form = document.getElementById("send-message")


    // send message on enter keypress
    input.addEventListener('keypress', function (event) {
        if (event.key === 'Enter') {
            event.preventDefault();
            parseMessage()
        }
    });

    form.addEventListener('submit', function (event) {
        event.preventDefault(); // Prevent form submission
        parseMessage()
        input.focus()
    });

    input.focus()
}

