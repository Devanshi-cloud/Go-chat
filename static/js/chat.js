(function() {
  let socket = null;
  let msgInput = null;
  let sendBtn = null;
  let messagesDiv = null;

  function init() {
    msgInput = document.getElementById('msg');
    sendBtn = document.getElementById('sendBtn');
    messagesDiv = document.getElementById('messages');

    // Connect to WebSocket
    connectWebSocket();

    // Event listeners
    sendBtn.addEventListener('click', sendMessage);
    msgInput.addEventListener('keypress', function(e) {
      if (e.key === 'Enter') {
        sendMessage();
      }
    });

    // Focus on input when page loads
    msgInput.focus();
  }

  function connectWebSocket() {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const wsUrl = `${protocol}//${window.location.host}/room`;
    
    socket = new WebSocket(wsUrl);

    socket.onopen = function(event) {
      console.log('Connected to chat room');
      addSystemMessage('Connected to chat room');
    };

    socket.onmessage = function(event) {
      const message = event.data;
      addMessage(message);
    };

    socket.onclose = function(event) {
      console.log('Disconnected from chat room');
      addSystemMessage('Disconnected from chat room');
      
      // Try to reconnect after 3 seconds
      setTimeout(connectWebSocket, 3000);
    };

    socket.onerror = function(error) {
      console.error('WebSocket error:', error);
      addSystemMessage('Connection error. Trying to reconnect...');
    };
  }

  function sendMessage() {
    const message = msgInput.value.trim();
    
    if (message && socket && socket.readyState === WebSocket.OPEN) {
      socket.send(message);
      msgInput.value = '';
      msgInput.focus();
    }
  }

  function addMessage(message) {
    const messageDiv = document.createElement('div');
    messageDiv.className = 'message';
    
    const contentDiv = document.createElement('div');
    contentDiv.className = 'message-content';
    contentDiv.textContent = message;
    
    const timeDiv = document.createElement('div');
    timeDiv.className = 'message-time';
    timeDiv.textContent = new Date().toLocaleTimeString();
    
    messageDiv.appendChild(contentDiv);
    messageDiv.appendChild(timeDiv);
    
    messagesDiv.appendChild(messageDiv);
    scrollToBottom();
  }

  function addSystemMessage(message) {
    const messageDiv = document.createElement('div');
    messageDiv.className = 'message';
    messageDiv.style.borderLeftColor = '#28a745';
    
    const contentDiv = document.createElement('div');
    contentDiv.className = 'message-content';
    contentDiv.style.fontStyle = 'italic';
    contentDiv.style.color = '#666';
    contentDiv.textContent = message;
    
    const timeDiv = document.createElement('div');
    timeDiv.className = 'message-time';
    timeDiv.textContent = new Date().toLocaleTimeString();
    
    messageDiv.appendChild(contentDiv);
    messageDiv.appendChild(timeDiv);
    
    messagesDiv.appendChild(messageDiv);
    scrollToBottom();
  }

  function scrollToBottom() {
    messagesDiv.scrollTop = messagesDiv.scrollHeight;
  }

  // Initialize when DOM is loaded
  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', init);
  } else {
    init();
  }
})(); 