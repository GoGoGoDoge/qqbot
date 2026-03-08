package qqbot

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"
)

type KnownUser struct {
	Type               string `json:"type"`
	OpenID             string `json:"openid"`
	AccountID          string `json:"accountId"`
	Nickname           string `json:"nickname,omitempty"`
	FirstInteractionAt int64  `json:"firstInteractionAt"`
	LastInteractionAt  int64  `json:"lastInteractionAt"`
}

type ProactiveStore struct {
	mu    sync.RWMutex
	users map[string]*KnownUser
	path  string
}

func NewProactiveStore(dataDir string) *ProactiveStore {
	path := filepath.Join(dataDir, "known-users.json")
	s := &ProactiveStore{users: make(map[string]*KnownUser), path: path}
	s.load()
	return s
}

func (s *ProactiveStore) userKey(typ, openid, accountID string) string {
	return accountID + ":" + typ + ":" + openid
}

func (s *ProactiveStore) Record(typ, openid, accountID, nickname string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := s.userKey(typ, openid, accountID)
	now := time.Now().Unix()

	if u, ok := s.users[key]; ok {
		u.LastInteractionAt = now
		if nickname != "" {
			u.Nickname = nickname
		}
	} else {
		s.users[key] = &KnownUser{
			Type:               typ,
			OpenID:             openid,
			AccountID:          accountID,
			Nickname:           nickname,
			FirstInteractionAt: now,
			LastInteractionAt:  now,
		}
	}
	s.save()
}

func (s *ProactiveStore) List(typ, accountID string, limit int) []*KnownUser {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []*KnownUser
	for _, u := range s.users {
		if (typ == "" || u.Type == typ) && (accountID == "" || u.AccountID == accountID) {
			result = append(result, u)
		}
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].LastInteractionAt > result[j].LastInteractionAt
	})

	if limit > 0 && len(result) > limit {
		result = result[:limit]
	}
	return result
}

func (s *ProactiveStore) load() {
	data, err := os.ReadFile(s.path)
	if err != nil {
		return
	}
	var users []*KnownUser
	if json.Unmarshal(data, &users) == nil {
		for _, u := range users {
			s.users[s.userKey(u.Type, u.OpenID, u.AccountID)] = u
		}
	}
}

func (s *ProactiveStore) save() {
	users := make([]*KnownUser, 0, len(s.users))
	for _, u := range s.users {
		users = append(users, u)
	}
	data, _ := json.MarshalIndent(users, "", "  ")
	os.MkdirAll(filepath.Dir(s.path), 0755)
	os.WriteFile(s.path, data, 0644)
}
