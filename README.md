<div align="center">

# QQBot Go SDK

Go client library for QQ Bot API - supports C2C private chats, group messages, and channel messages.

[![Go Version](https://img.shields.io/badge/Go-1.18+-00ADD8)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-green)](./LICENSE)
[![QQ Bot](https://img.shields.io/badge/QQ_Bot-API_v2-red)](https://bot.q.qq.com/wiki/)

</div>

---

## ✨ Features

- 🔒 **Multi-scenario Support** - C2C private chats, group @messages, channel messages
- 🖼️ **Rich Media** - Image and file upload support
- ⏰ **Proactive Messaging** - Send messages with user tracking
- 📝 **Markdown** - Full Markdown format support
- 🔄 **Auto Token Refresh** - Automatic token management
- 🌐 **HTTP API Server** - Built-in server for proactive messages

---

## 📦 Installation

```bash
go get github.com/GoGoGoDoge/qqbot
```

---

## 🚀 Quick Start

```go
package main

import (
    "log"
    "os"
    "os/signal"
    "github.com/GoGoGoDoge/qqbot"
)

func main() {
    client := qqbot.NewClient(
        os.Getenv("QQBOT_APP_ID"),
        os.Getenv("QQBOT_CLIENT_SECRET"),
    )

    client.Gateway.OnC2CMessage(func(e qqbot.C2CMessageEvent) {
        log.Printf("C2C: %s", e.Content)
        client.API.SendC2CMessage(e.Author.UserOpenID, "Echo: "+e.Content, e.ID)
    })

    client.Connect()
    log.Println("Bot started")

    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt)
    <-c
}
```

---

## 🛠️ CLI Tools

### Send Proactive Message

```bash
export QQBOT_APP_ID=your_app_id
export QQBOT_CLIENT_SECRET=your_secret

# Send message
go run ./cmd/send-proactive --to <openid> --text "Hello" --type c2c

# List known users
go run ./cmd/send-proactive --list
```

### Proactive API Server

```bash
go run ./cmd/proactive-server --port 3721
```

API endpoints:
- `POST /send` - Send proactive message
- `GET /users` - List known users
- `GET /users/stats` - User statistics

Example:
```bash
curl -X POST http://localhost:3721/send \
  -H "Content-Type: application/json" \
  -d '{"to":"openid","text":"Hello!","type":"c2c"}'
```

---

## 📁 Project Structure

```
qqbot/
├── client.go          # Main client entry point
├── api.go            # API client with token management
├── api_send.go       # Message sending methods
├── gateway.go        # WebSocket gateway
├── gateway_handlers.go # Event handlers
├── proactive.go      # Proactive messaging & user tracking
├── types.go          # Public types
├── users.go          # User operations
├── imageserver.go    # Image handling
├── internal/         # Internal packages
│   └── session/
│       └── store.go  # Session storage
└── cmd/              # CLI tools
    ├── send-proactive/
    └── proactive-server/
```

---

## 📖 Documentation

For detailed API documentation and usage examples, visit the [QQ Bot Official Documentation](https://bot.q.qq.com/wiki/).

---

## 📄 License

MIT License - see [LICENSE](./LICENSE) for details.
