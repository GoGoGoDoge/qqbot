package qqbot

import (
	"encoding/json"
	"log"
	"time"
)

func (g *Gateway) read() {
	for {
		var p WSPayload
		if err := g.conn.ReadJSON(&p); err != nil {
			return
		}
		if p.S > 0 {
			g.seq = p.S
		}
		switch p.Op {
		case 10:
			g.identify(p)
		case 0:
			g.dispatch(p)
		}
	}
}

func (g *Gateway) identify(p WSPayload) {
	data := p.D.(map[string]interface{})
	interval := int(data["heartbeat_interval"].(float64))
	token, _ := g.api.GetAccessToken()
	g.conn.WriteJSON(WSPayload{Op: 2, D: map[string]interface{}{
		"token": token, "intents": 0 | (1 << 30) | (1 << 25) | (1 << 12), "shard": []int{0, 1},
	}})
	go g.heartbeat(interval)
}

func (g *Gateway) heartbeat(ms int) {
	t := time.NewTicker(time.Duration(ms) * time.Millisecond)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			g.conn.WriteJSON(WSPayload{Op: 1, D: g.seq})
		case <-g.done:
			return
		}
	}
}

func (g *Gateway) dispatch(p WSPayload) {
	switch p.T {
	case "READY":
		log.Println("Gateway ready")
	case "C2C_MESSAGE_CREATE":
		var e C2CMessageEvent
		b, _ := json.Marshal(p.D)
		json.Unmarshal(b, &e)
		if g.client != nil && g.client.Proactive != nil {
			g.client.Proactive.Record("c2c", e.Author.UserOpenID, "default", "")
		}
		if g.onC2C != nil {
			g.onC2C(e)
		}
	case "GROUP_AT_MESSAGE_CREATE":
		var e GroupMessageEvent
		b, _ := json.Marshal(p.D)
		json.Unmarshal(b, &e)
		if g.client != nil && g.client.Proactive != nil {
			g.client.Proactive.Record("group", e.GroupOpenID, "default", "")
		}
		if g.onGroup != nil {
			g.onGroup(e)
		}
	case "AT_MESSAGE_CREATE":
		var e ChannelMessageEvent
		b, _ := json.Marshal(p.D)
		json.Unmarshal(b, &e)
		if g.onChannel != nil {
			g.onChannel(e)
		}
	}
}

func (g *Gateway) Close() error {
	close(g.done)
	if g.conn != nil {
		return g.conn.Close()
	}
	return nil
}
