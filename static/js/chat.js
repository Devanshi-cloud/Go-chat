const params = new URLSearchParams(window.location.search);
const room = params.get("room");

if (!room) {
  alert("No room specified. Redirecting to homepage...");
  window.location.href = "/";
}

// Display room name in header
document.getElementById("roomName").textContent = room;

// Get username from localStorage or prompt user
let username = localStorage.getItem("chat_username");
if (!username) {
  username = prompt("Enter your username:");
  if (username) {
    localStorage.setItem("chat_username", username);
  } else {
    username = "Anonymous";
  }
}

const socket = new WebSocket(`ws://${location.host}/room?room=${room}`);

socket.onopen = function(event) {
  console.log('Connected to chat room:', room);
  addSystemMessage('Connected to chat room: ' + room);
};

socket.onclose = function(event) {
  console.log('Disconnected from chat room');
  addSystemMessage('Disconnected from chat room');
};

socket.onerror = function(error) {
  console.error('WebSocket error:', error);
  addSystemMessage('Connection error');
};

socket.onmessage = (event) => {
  try {
    const data = JSON.parse(event.data);

    // Create the container div
    const msgContainer = document.createElement("div");
    msgContainer.classList.add("message-container");

    // Create the username div
    const usernameDiv = document.createElement("div");
    usernameDiv.classList.add("username");
    usernameDiv.textContent = data.name;

    // Create the message div
    const messageDiv = document.createElement("div");
    messageDiv.classList.add("message");
    messageDiv.textContent = data.message;

    // Append username and message in correct order
    msgContainer.appendChild(usernameDiv);
    msgContainer.appendChild(messageDiv);

    // Append the whole message container to the messages div
    document.getElementById("messages").appendChild(msgContainer);

    // Auto-scroll
    const messagesDiv = document.getElementById("messages");
    messagesDiv.scrollTop = messagesDiv.scrollHeight;

  } catch (err) {
    console.error("Invalid JSON received:", event.data);
    // Try to display as plain text
    addMessage("System", event.data);
  }
};

function addSystemMessage(message) {
  addMessage("System", message);
}

function addMessage(name, message) {
  // Create the container div
  const msgContainer = document.createElement("div");
  msgContainer.classList.add("message-container");

  // Create the username div
  const usernameDiv = document.createElement("div");
  usernameDiv.classList.add("username");
  usernameDiv.textContent = name;

  // Create the message div
  const messageDiv = document.createElement("div");
  messageDiv.classList.add("message");
  messageDiv.textContent = message;

  // Append username and message in correct order
  msgContainer.appendChild(usernameDiv);
  msgContainer.appendChild(messageDiv);

  // Append the whole message container to the messages div
  document.getElementById("messages").appendChild(msgContainer);

  // Auto-scroll
  const messagesDiv = document.getElementById("messages");
  messagesDiv.scrollTop = messagesDiv.scrollHeight;
}

function sendMessage() {
  const input = document.getElementById("msg");
  if (input.value.trim() !== "") {
    const messageData = {
      name: username,
      message: input.value.trim(),
      room: room
    };
    socket.send(JSON.stringify(messageData));
    input.value = "";
  }
}

document.getElementById("sendBtn").addEventListener("click", sendMessage);

document.getElementById("msg").addEventListener("keyup", function (event) {
  if (event.key === "Enter") {
    sendMessage();
  }
});