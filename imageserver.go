package qqbot

import (
	"fmt"
	"net/http"
	"path/filepath"
	"sync"
)

type ImageServer struct {
	mu      sync.RWMutex
	files   map[string]string
	baseURL string
	port    int
}

func NewImageServer(port int, baseURL string) *ImageServer {
	return &ImageServer{files: make(map[string]string), port: port, baseURL: baseURL}
}

func (s *ImageServer) Start() error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		s.mu.RLock()
		path, ok := s.files[r.URL.Path]
		s.mu.RUnlock()
		if !ok {
			http.NotFound(w, r)
			return
		}
		http.ServeFile(w, r, path)
	})
	go http.ListenAndServe(fmt.Sprintf(":%d", s.port), nil)
	return nil
}

func (s *ImageServer) AddFile(localPath string) string {
	s.mu.Lock()
	defer s.mu.Unlock()
	name := filepath.Base(localPath)
	urlPath := "/" + name
	s.files[urlPath] = localPath
	return s.baseURL + urlPath
}
