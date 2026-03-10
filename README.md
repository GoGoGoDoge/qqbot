<div align="center">

# QQBot Go SDK

Go client library for QQ Bot API - supports C2C private chats, group messages, and channel messages.

Built on top of [tencent-connect/botgo](https://github.com/tencent-connect/botgo) official SDK.

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
        os.Getenv("QQBOT_APP_SECRET"),
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
export QQBOT_APP_SECRET=your_secret
```

#### Send Text Message

```bash
go run ./cmd/send-proactive --to <openid> --text "Hello" --type c2c
```

#### Send Media from URL

```bash
# Send image (two-step method, doesn't consume proactive message quota)
go run ./cmd/send-proactive --to <openid> --media "https://example.com/image.jpg" --media-type image

# Send video
go run ./cmd/send-proactive --to <openid> --media "https://example.com/video.mp4" --media-type video

# Send voice
go run ./cmd/send-proactive --to <openid> --media "https://example.com/audio.mp3" --media-type voice

# Send file
go run ./cmd/send-proactive --to <openid> --media "https://example.com/doc.pdf" --media-type file

# Direct send (consumes proactive message quota)
go run ./cmd/send-proactive --to <openid> --media "https://example.com/image.jpg" --direct
```

#### Send Local File

```bash
# Send local image (max 5MB due to base64 encoding)
go run ./cmd/send-proactive --to <openid> --file ./image.jpg --media-type image

# Send local video
go run ./cmd/send-proactive --to <openid> --file ./video.mp4 --media-type video

# Send local file
go run ./cmd/send-proactive --to <openid> --file ./document.pdf --media-type file
```

#### List Known Users

```bash
# List all users
go run ./cmd/send-proactive --list

# List C2C users only
go run ./cmd/send-proactive --list --type c2c

# List group users only
go run ./cmd/send-proactive --list --type group
```

#### Command Options

- `--to` - Target user OpenID (required for sending)
- `--text` - Message text content
- `--type` - Message type: `c2c` (default) or `group`
- `--media` - Media URL for remote files
- `--media-type` - Media type: `image`, `video`, `voice`, or `file` (default: `image`)
- `--direct` - Send media directly (consumes proactive message quota)
- `--file` - Local file path (max 5MB)
- `--list` - List known users

### Proactive API Server

Start the HTTP API server:

```bash
go run ./cmd/proactive-server --port 3721
```

#### API Endpoints

**GET /** - Server info
```bash
curl http://localhost:3721/
```

**POST /send** - Send proactive message
```bash
# Send text
curl -X POST http://localhost:3721/send \
  -H "Content-Type: application/json" \
  -d '{"to":"openid","text":"Hello!","type":"c2c"}'

# Send to group
curl -X POST http://localhost:3721/send \
  -H "Content-Type: application/json" \
  -d '{"to":"group_openid","text":"Hello group!","type":"group"}'
```

**GET /users** - List known users
```bash
# List all users
curl http://localhost:3721/users

# Filter by type
curl http://localhost:3721/users?type=c2c
```

**GET /users/stats** - User statistics
```bash
curl http://localhost:3721/users/stats
```

#### Server Options

- `--port` - Server port (default: `3721`)

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
