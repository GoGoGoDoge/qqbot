package qqbot

import "github.com/GoGoGoDoge/qqbot/internal/session"

type Client struct {
	API       *API
	Gateway   *Gateway
	Session   *session.Store
	Proactive *ProactiveStore
}

func NewClient(appID, clientSecret string) *Client {
	api := NewAPI(appID, clientSecret)
	c := &Client{
		API:       api,
		Gateway:   NewGateway(api),
		Session:   session.NewStore("."),
		Proactive: NewProactiveStore("."),
	}
	c.Gateway.client = c
	return c
}

func (c *Client) Connect() error {
	return c.Gateway.Connect()
}
