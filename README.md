# Chat Application

A real-time chat application built with Go and WebSocket.

## Features

- Real-time messaging using WebSocket
- Modern, responsive UI with gradient design
- Auto-reconnection on connection loss
- Message timestamps
- System status messages
- Clean and intuitive interface

## Project Structure

```
chat/
├── main.go              # Main application entry point
├── client.go            # WebSocket client implementation
├── room.go              # Chat room management
├── go.mod               # Go module dependencies
├── template/
│   └── chat.html        # HTML template for chat interface
├── static/
│   ├── css/
│   │   └── styles.css   # Modern CSS styling
│   └── js/
│       └── chat.js      # Frontend JavaScript for WebSocket
└── README.md            # This file
```

## Getting Started

### Prerequisites

- Go 1.24.5 or higher

### Installation

1. Clone the repository:
```bash
git clone https://github.com/yourusername/chat.git
cd chat
```

2. Install dependencies:
```bash
go mod tidy
```

3. Run the application:
```bash
go run .
```

4. Open your browser and navigate to:
```
http://localhost:8080
```

## Usage

- The chat application runs on port 8080 by default
- Multiple users can connect simultaneously
- Messages are broadcast to all connected users
- The interface automatically reconnects if the connection is lost

## Development

The application uses:
- **Backend**: Go with Gorilla WebSocket
- **Frontend**: Vanilla JavaScript with modern CSS
- **Architecture**: WebSocket-based real-time communication

## License

MIT License 