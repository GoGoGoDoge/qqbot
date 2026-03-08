package qqbot

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type UserRecord struct {
	OpenID    string    `json:"openid"`
	FirstSeen time.Time `json:"first_seen"`
	LastSeen  time.Time `json:"last_seen"`
}

type UserStore struct {
	mu    sync.RWMutex
	users map[string]*UserRecord
	path  string
}

func NewUserStore(dataDir string) *UserStore {
	path := filepath.Join(dataDir, "users.json")
	s := &UserStore{users: make(map[string]*UserRecord), path: path}
	s.load()
	return s
}

func (s *UserStore) Record(openID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if u, ok := s.users[openID]; ok {
		u.LastSeen = time.Now()
	} else {
		s.users[openID] = &UserRecord{OpenID: openID, FirstSeen: time.Now(), LastSeen: time.Now()}
	}
}

func (s *UserStore) Flush() {
	s.mu.RLock()
	defer s.mu.RUnlock()
	data, _ := json.MarshalIndent(s.users, "", "  ")
	os.WriteFile(s.path, data, 0644)
}

func (s *UserStore) load() {
	data, err := os.ReadFile(s.path)
	if err != nil {
		return
	}
	var users []*UserRecord
	json.Unmarshal(data, &users)
	for _, u := range users {
		s.users[u.OpenID] = u
	}
}
