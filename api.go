package qqbot

import (
	"bytes"
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

const (
	APIBase  = "https://api.sgroup.qq.com"
	TokenURL = "https://bots.qq.com/app/getAppAccessToken"
)

type API struct {
	appID        string
	clientSecret string
	token        string
	tokenExpiry  time.Time
	mu           sync.RWMutex
	MarkdownMode bool
	stopRefresh  chan struct{}
}

func NewAPI(appID, clientSecret string) *API {
	return &API{appID: appID, clientSecret: clientSecret, stopRefresh: make(chan struct{})}
}

func (a *API) StartTokenRefresh() {
	go func() {
		ticker := time.NewTicker(10 * time.Minute)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				a.GetAccessToken()
			case <-a.stopRefresh:
				return
			}
		}
	}()
}

func (a *API) StopTokenRefresh() {
	close(a.stopRefresh)
}

func (a *API) GetAccessToken() (string, error) {
	a.mu.RLock()
	if time.Now().Before(a.tokenExpiry.Add(-5 * time.Minute)) {
		token := a.token
		a.mu.RUnlock()
		return token, nil
	}
	a.mu.RUnlock()

	a.mu.Lock()
	defer a.mu.Unlock()

	body, _ := json.Marshal(map[string]string{"appId": a.appID, "clientSecret": a.clientSecret})
	resp, err := http.Post(TokenURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	json.NewDecoder(resp.Body).Decode(&result)
	a.token = result.AccessToken
	a.tokenExpiry = time.Now().Add(time.Duration(result.ExpiresIn) * time.Second)
	return a.token, nil
}
