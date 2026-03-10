package qqbot

import (
	"context"
	"time"

	"github.com/tencent-connect/botgo"
	"github.com/tencent-connect/botgo/openapi"
	"github.com/tencent-connect/botgo/token"
	"golang.org/x/oauth2"
)

type API struct {
	appID       string
	appSecret   string
	tokenSource oauth2.TokenSource
	client      openapi.OpenAPI
	ctx         context.Context
	cancel      context.CancelFunc
	MarkdownMode bool
}

func NewAPI(appID, appSecret string) *API {
	ctx, cancel := context.WithCancel(context.Background())

	credentials := &token.QQBotCredentials{
		AppID:     appID,
		AppSecret: appSecret,
	}
	tokenSource := token.NewQQBotTokenSource(credentials)

	return &API{
		appID:       appID,
		appSecret:   appSecret,
		tokenSource: tokenSource,
		client:      botgo.NewOpenAPI(appID, tokenSource).WithTimeout(5 * time.Second),
		ctx:         ctx,
		cancel:      cancel,
	}
}

func (a *API) StartTokenRefresh() error {
	return token.StartRefreshAccessToken(a.ctx, a.tokenSource)
}

func (a *API) StopTokenRefresh() {
	if a.cancel != nil {
		a.cancel()
	}
}

func (a *API) GetAppID() string {
	return a.appID
}

func (a *API) GetTokenSource() oauth2.TokenSource {
	return a.tokenSource
}

func (a *API) GetClient() openapi.OpenAPI {
	return a.client
}

func (a *API) GetContext() context.Context {
	return a.ctx
}
