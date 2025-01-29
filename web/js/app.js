document.addEventListener('DOMContentLoaded', async function () {
  const chat = document.getElementById('chat');
  const host = window.location.hostname;
  const port = window.location.port;
  wsHost = `${window.location.protocol}//${host}${port}`

  const manager = new io.Manager(wsHost, {
    transports: ["polling", "websocket" /*, "webtransport"*/],
    reconnectionDelayMax: 10000,
    reconnectionAttempts: 5,
    reconnectionDelay: 5000,
    timeout: 10000,
  });

  const chatSocket = manager.socket("chat", {});

  chatSocket.on('connect', function () {
    displayMessage('System: Connected to the chat server.', 'system');
  });

  chatSocket.on('message', function (event) {
    displayMessage(`Received Socket.IO message: ${JSON.parse(event)}`);
  })

  chatSocket.on("reconnect", (attempt) => {
    console.log(`Reconnected after ${attempt} attempts`);
  });

  chatSocket.on("connect_error", (error) => {
    if (chatSocket.active) {
      // temporary failure, the socket will automatically try to reconnect
    } else {
      // the connection was denied by the server
      // in that case, `socket.connect()` must be manually called in order to reconnect
      displayMessage(`System: ${error.message}`, 'error');
    }
  });

  chatSocket.on('disconnect', function () {
    displayMessage('System: Disconnected from the chat server.', 'error');
  });

  chatSocket.on('error', function (error) {
    displayMessage('System: Socket.IO encountered an error.', 'error');
  });

  function displayMessage(text, type) {
    const timestamp = new Date().toLocaleTimeString();
    const messageElement = document.createElement('div');
    messageElement.className = type; // Apply different styling based on the message type
    messageElement.textContent = `[${timestamp}] ${text}`;
    chat.appendChild(messageElement);
    chat.scrollTop = chat.scrollHeight; // Automatically scroll to the latest message
  }
})