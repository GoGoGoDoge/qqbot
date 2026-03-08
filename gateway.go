package qqbot

import (
	"github.com/gorilla/websocket"
)

type Gateway struct {
	api           *API
	client        *Client
	conn          *websocket.Conn
	seq           int
	done          chan struct{}
	onC2C         func(C2CMessageEvent)
	onGroup       func(GroupMessageEvent)
	onChannel     func(ChannelMessageEvent)
	MarkdownMode  bool
}

func NewGateway(api *API) *Gateway {
	return &Gateway{api: api, done: make(chan struct{})}
}

func (g *Gateway) OnC2CMessage(h func(C2CMessageEvent))         { g.onC2C = h }
func (g *Gateway) OnGroupMessage(h func(GroupMessageEvent))     { g.onGroup = h }
func (g *Gateway) OnChannelMessage(h func(ChannelMessageEvent)) { g.onChannel = h }

func (g *Gateway) Connect() error {
	url, err := g.api.GetGatewayURL()
	if err != nil {
		return err
	}
	g.conn, _, err = websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return err
	}
	go g.read()
	return nil
}
