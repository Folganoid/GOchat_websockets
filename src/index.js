	window.socket = new WebSocket("ws://" + location.host + "/ws");
	
function sendMessage(msg) {
	socket.send(msg)
}

function handleSubmit() {
	var el = document.getElementById("chat-msg")
	sendMessage(el.value);
	el.value = "";
	return false;
}

function setUpSocket(onmessage) {
	socket.onopen = function() {
		console.log("Connected");
	};

	socket.onclose = function(event) {
		if (event.wasClean) {
			console.log('Conection closed');
		} else {
    		console.log('Connection interrupt');
		}
	console.log('Code: ' + event.code + ' reason: ' + event.reason);
	};

	socket.onmessage = function(event) {
		console.log("Recieved data " + event.data);
	};

	socket.onmessage = onmessage;

	socket.onerror = function(error) {
		console.log("Error " + error.message);
	};
}

function displayMessage(msg) {
	var container = document.getElementById("container")

	var div = document.createElement("div")
	div.className = 'message'

	var textNode = document.createTextNode(msg.data);

	div.appendChild(textNode)
	container.appendChild(div)
}