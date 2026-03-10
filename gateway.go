package qqbot

import (
	"log"
	"sync"

	"github.com/tencent-connect/botgo"
	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/event"
)

type Gateway struct {
	api            *API
	client         *Client
	sessionManager botgo.SessionManager
	onC2C          func(C2CMessageEvent)
	onGroup        func(GroupMessageEvent)
	onChannel      func(ChannelMessageEvent)
	mu             sync.RWMutex
	running        bool
}

func NewGateway(api *API) *Gateway {
	return &Gateway{
		api:            api,
		sessionManager: botgo.NewSessionManager(),
	}
}

func (g *Gateway) OnC2CMessage(h func(C2CMessageEvent))         { g.onC2C = h }
func (g *Gateway) OnGroupMessage(h func(GroupMessageEvent))     { g.onGroup = h }
func (g *Gateway) OnChannelMessage(h func(ChannelMessageEvent)) { g.onChannel = h }

func (g *Gateway) Connect() error {
	// 启动 token 自动刷新
	if err := g.api.StartTokenRefresh(); err != nil {
		return err
	}

	// 注册事件处理器
	intent := event.RegisterHandlers(
		g.handleC2CMessage(),
		g.handleGroupATMessage(),
		g.handleChannelMessage(),
	)

	// 获取 WebSocket 信息
	wsInfo, err := g.api.GetClient().WS(g.api.GetContext(), nil, "")
	if err != nil {
		return err
	}

	log.Printf("WebSocket info: shards=%d", wsInfo.Shards)

	// 启动 WebSocket 连接
	g.setRunning(true)
	go func() {
		if err := g.sessionManager.Start(wsInfo, g.api.GetTokenSource(), &intent); err != nil {
			log.Printf("WebSocket session error: %v", err)
			g.setRunning(false)
		}
	}()

	return nil
}

func (g *Gateway) IsRunning() bool {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.running
}

func (g *Gateway) setRunning(running bool) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.running = running
	if running {
		log.Println("Gateway 已启动")
	} else {
		log.Println("Gateway 已停止")
	}
}

func (g *Gateway) Close() error {
	g.setRunning(false)
	g.api.StopTokenRefresh()
	return nil
}

// handleC2CMessage 处理私聊消息
func (g *Gateway) handleC2CMessage() event.C2CMessageEventHandler {
	return func(event *dto.WSPayload, data *dto.WSC2CMessageData) error {
		if data.Author == nil || data.Author.ID == "" {
			return nil
		}

		log.Printf("收到C2C消息: %s (from: %s)", data.Content, data.Author.ID)

		if g.onC2C != nil {
			// 转换 Attachments
			attachments := make([]Attachment, 0, len(data.Attachments))
			for _, att := range data.Attachments {
				attachments = append(attachments, Attachment{
					ContentType: att.ContentType,
					Filename:    att.FileName,
					Height:      att.Height,
					Width:       att.Width,
					Size:        att.Size,
					URL:         att.URL,
				})
			}

			g.onC2C(C2CMessageEvent{
				ID:          data.ID,
				Content:     data.Content,
				Timestamp:   string(data.Timestamp),
				Attachments: attachments,
				Author: Author{
					ID:         data.Author.ID,
					UserOpenID: data.Author.UnionOpenID,
				},
			})
		}

		return nil
	}
}

// handleGroupATMessage 处理群聊@消息
func (g *Gateway) handleGroupATMessage() event.GroupATMessageEventHandler {
	return func(event *dto.WSPayload, data *dto.WSGroupATMessageData) error {
		if data.Author == nil || data.Author.ID == "" {
			return nil
		}

		log.Printf("收到群消息: %s (group: %s)", data.Content, data.GroupID)

		if g.onGroup != nil {
			// 转换 Attachments
			attachments := make([]Attachment, 0, len(data.Attachments))
			for _, att := range data.Attachments {
				attachments = append(attachments, Attachment{
					ContentType: att.ContentType,
					Filename:    att.FileName,
					Height:      att.Height,
					Width:       att.Width,
					Size:        att.Size,
					URL:         att.URL,
				})
			}

			g.onGroup(GroupMessageEvent{
				ID:          data.ID,
				Content:     data.Content,
				Timestamp:   string(data.Timestamp),
				GroupOpenID: data.GroupID,
				Attachments: attachments,
				Author: Author{
					ID:         data.Author.ID,
					UserOpenID: data.Author.UnionOpenID,
				},
			})
		}

		return nil
	}
}

// handleChannelMessage 处理频道消息
func (g *Gateway) handleChannelMessage() event.ATMessageEventHandler {
	return func(event *dto.WSPayload, data *dto.WSATMessageData) error {
		if data.Author == nil || data.Author.ID == "" {
			return nil
		}

		log.Printf("收到频道消息: %s (channel: %s)", data.Content, data.ChannelID)

		if g.onChannel != nil {
			g.onChannel(ChannelMessageEvent{
				ID:        data.ID,
				Content:   data.Content,
				ChannelID: data.ChannelID,
				Author: Author{
					ID:       data.Author.ID,
					Username: data.Author.Username,
				},
			})
		}

		return nil
	}
}
