package main

import (
	"log"
	"os"
	"os/signal"

	qqbot "github.com/GoGoGoDoge/qqbot"
)

func main() {
	client := qqbot.NewClient(os.Getenv("QQBOT_APPID"), os.Getenv("QQBOT_SECRET"))
	client.API.MarkdownMode = true
	defer client.API.StopTokenRefresh()

	client.Gateway.OnC2CMessage(func(e qqbot.C2CMessageEvent) {
		log.Printf("C2C: %s", e.Content)
		client.API.SendC2CInputNotify(e.Author.UserOpenID, e.ID, 5)
		client.API.SendC2CMessage(e.Author.UserOpenID, "Echo: "+e.Content, e.ID)
	})

	client.Gateway.OnGroupMessage(func(e qqbot.GroupMessageEvent) {
		log.Printf("Group: %s", e.Content)
		client.API.SendGroupMessage(e.GroupOpenID, "Echo: "+e.Content, e.ID)
	})

	client.Gateway.OnChannelMessage(func(e qqbot.ChannelMessageEvent) {
		log.Printf("Channel: %s", e.Content)
		client.API.SendChannelMessage(e.ChannelID, "Echo: "+e.Content, e.ID)
	})

	client.Connect()
	log.Println("Bot started")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	client.Gateway.Close()
}
