package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/GoGoGoDoge/qqbot"
)

var store *qqbot.ProactiveStore
var client *qqbot.Client

func main() {
	port := flag.Int("port", 3721, "Server port")
	flag.Parse()

	appID := os.Getenv("QQBOT_APP_ID")
	secret := os.Getenv("QQBOT_CLIENT_SECRET")
	if appID == "" || secret == "" {
		log.Fatal("QQBOT_APP_ID and QQBOT_CLIENT_SECRET required")
	}

	client = qqbot.NewClient(appID, secret)
	store = qqbot.NewProactiveStore(".")

	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/send", handleSend)
	http.HandleFunc("/users", handleUsers)
	http.HandleFunc("/users/stats", handleStats)

	addr := fmt.Sprintf(":%d", *port)
	log.Printf("Server running at http://localhost%s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]interface{}{
		"name":    "QQBot Proactive API",
		"version": "1.0.0",
	})
}

func handleSend(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", 405)
		return
	}

	var req struct {
		To   string `json:"to"`
		Text string `json:"text"`
		Type string `json:"type"`
	}
	if json.NewDecoder(r.Body).Decode(&req) != nil || req.To == "" || req.Text == "" {
		http.Error(w, "Invalid request", 400)
		return
	}

	if req.Type == "" {
		req.Type = "c2c"
	}

	var resp *qqbot.MessageResponse
	var err error

	if req.Type == "c2c" {
		resp, err = client.API.SendC2CMessage(req.To, req.Text, "")
	} else {
		resp, err = client.API.SendGroupMessage(req.To, req.Text, "")
	}

	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "error": err.Error()})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":   true,
		"messageId": resp.ID,
		"timestamp": resp.Timestamp,
	})
}

func handleUsers(w http.ResponseWriter, r *http.Request) {
	typ := r.URL.Query().Get("type")
	users := store.List(typ, "", 0)
	json.NewEncoder(w).Encode(map[string]interface{}{"total": len(users), "users": users})
}

func handleStats(w http.ResponseWriter, r *http.Request) {
	users := store.List("", "", 0)
	c2c, group := 0, 0
	for _, u := range users {
		if u.Type == "c2c" {
			c2c++
		} else if u.Type == "group" {
			group++
		}
	}
	json.NewEncoder(w).Encode(map[string]int{"total": len(users), "c2c": c2c, "group": group})
}
