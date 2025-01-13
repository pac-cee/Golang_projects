let currentUser = null;
let currentRoom = null;
let socket = null;

// DOM Elements
const loginSection = document.getElementById('login-section');
const chatSection = document.getElementById('chat-section');
const usernameInput = document.getElementById('username-input');
const loginBtn = document.getElementById('login-btn');
const roomsList = document.getElementById('rooms-list');
const newRoomInput = document.getElementById('new-room-input');
const createRoomBtn = document.getElementById('create-room-btn');
const roomHeader = document.getElementById('room-header');
const roomUsers = document.getElementById('room-users');
const messages = document.getElementById('messages');
const messageInput = document.getElementById('message-input');
const sendBtn = document.getElementById('send-btn');

// Event Listeners
loginBtn.addEventListener('click', handleLogin);
createRoomBtn.addEventListener('click', handleCreateRoom);
sendBtn.addEventListener('click', handleSendMessage);
messageInput.addEventListener('keypress', (e) => {
    if (e.key === 'Enter') handleSendMessage();
});

// Login Handler
function handleLogin() {
    const username = usernameInput.value.trim();
    if (username) {
        currentUser = {
            id: generateUserId(),
            username: username
        };
        loginSection.classList.add('hidden');
        chatSection.classList.remove('hidden');
        fetchRooms();
    }
}

// Room Handlers
async function fetchRooms() {
    const response = await fetch('http://localhost:8080/api/rooms');
    const data = await response.json();
    displayRooms(data.rooms);
}

async function handleCreateRoom() {
    const roomName = newRoomInput.value.trim();
    if (roomName) {
        const response = await fetch('http://localhost:8080/api/rooms', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ name: roomName })
        });
        const data = await response.json();
        newRoomInput.value = '';
        fetchRooms();
    }
}

function displayRooms(rooms) {
    roomsList.innerHTML = rooms.map(room => `
        <div class="room-item ${currentRoom === room ? 'active' : ''}" 
             onclick="joinRoom('${room}')">
            ${room}
        </div>
    `).join('');
}

// WebSocket Handlers
function joinRoom(roomId) {
    if (socket) {
        socket.close();
    }

    currentRoom = roomId;
    displayRooms([...document.querySelectorAll('.room-item')].map(el => el.textContent.trim()));

    const wsProtocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const wsUrl = `${wsProtocol}//${window.location.host}/ws/${roomId}?user_id=${currentUser.id}&username=${currentUser.username}`;
    
    socket = new WebSocket(wsUrl);
    
    socket.onopen = () => {
        console.log('Connected to room:', roomId);
        messages.innerHTML = '';
        roomHeader.querySelector('h2').textContent = `Room: ${roomId}`;
    };

    socket.onmessage = (event) => {
        const message = JSON.parse(event.data);
        displayMessage(message);
    };

    socket.onerror = (error) => {
        console.error('WebSocket error:', error);
    };

    socket.onclose = () => {
        console.log('Disconnected from room:', roomId);
    };
}

function handleSendMessage() {
    const content = messageInput.value.trim();
    if (content && socket && socket.readyState === WebSocket.OPEN) {
        const message = {
            content: content
        };
        socket.send(JSON.stringify(message));
        messageInput.value = '';
    }
}

// Message Display
function displayMessage(message) {
    const messageDiv = document.createElement('div');
    messageDiv.className = `message ${getMessageClass(message)}`;
    
    const timestamp = new Date(message.timestamp).toLocaleTimeString();
    
    if (message.type === 'message') {
        messageDiv.innerHTML = `
            <div class="username">${message.username}</div>
            <div class="content">${escapeHtml(message.content)}</div>
            <div class="timestamp">${timestamp}</div>
        `;
    } else {
        messageDiv.innerHTML = `
            <div class="content">${message.content}</div>
            <div class="timestamp">${timestamp}</div>
        `;
    }

    messages.appendChild(messageDiv);
    messages.scrollTop = messages.scrollHeight;
}

// Utility Functions
function getMessageClass(message) {
    if (message.type === 'join' || message.type === 'leave') {
        return message.type;
    }
    return message.user_id === currentUser.id ? 'user' : 'other';
}

function generateUserId() {
    return 'user-' + Math.random().toString(36).substr(2, 9);
}

function escapeHtml(unsafe) {
    return unsafe
        .replace(/&/g, "&amp;")
        .replace(/</g, "&lt;")
        .replace(/>/g, "&gt;")
        .replace(/"/g, "&quot;")
        .replace(/'/g, "&#039;");
}
