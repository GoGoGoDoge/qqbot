package session

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

type Store struct {
	mu       sync.RWMutex
	sessions map[string]map[string]interface{}
	path     string
}

func NewStore(dataDir string) *Store {
	path := filepath.Join(dataDir, "sessions.json")
	s := &Store{sessions: make(map[string]map[string]interface{}), path: path}
	s.load()
	return s
}

func (s *Store) Get(key string) (map[string]interface{}, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.sessions[key]
	return v, ok
}

func (s *Store) Set(key string, value map[string]interface{}) {
	s.mu.Lock()
	s.sessions[key] = value
	s.mu.Unlock()
	s.save()
}

func (s *Store) Delete(key string) {
	s.mu.Lock()
	delete(s.sessions, key)
	s.mu.Unlock()
	s.save()
}

func (s *Store) load() {
	data, err := os.ReadFile(s.path)
	if err != nil {
		return
	}
	json.Unmarshal(data, &s.sessions)
}

func (s *Store) save() {
	data, _ := json.MarshalIndent(s.sessions, "", "  ")
	os.MkdirAll(filepath.Dir(s.path), 0755)
	os.WriteFile(s.path, data, 0644)
}
