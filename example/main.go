package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	qqbot "github.com/GoGoGoDoge/qqbot"
)

func main() {
	fmt.Println("os get env: ", os.Getenv("QQBOT_APP_ID"), "|", os.Getenv("QQBOT_APP_SECRET"))
	client := qqbot.NewClient(os.Getenv("QQBOT_APP_ID"), os.Getenv("QQBOT_APP_SECRET"))
	client.API.MarkdownMode = true
	defer client.API.StopTokenRefresh()

	client.Gateway.OnC2CMessage(func(e qqbot.C2CMessageEvent) {
		log.Printf("收到C2C消息: %s (from: %s)", e.Content, e.Author.ID)
		if err := client.API.SendC2CInputNotify(e.Author.ID, e.ID, 5); err != nil {
			log.Printf("❌ 发送C2C输入通知失败: %v", err)
		}
		resp, err := client.API.SendC2CMessage(e.Author.ID, "Echo: "+e.Content, e.ID)
		if err != nil {
			log.Printf("❌ 发送C2C消息失败: %v", err)
		} else {
			log.Printf("✅ 消息发送成功: id=%s", resp.ID)
		}
	})

	client.Gateway.OnGroupMessage(func(e qqbot.GroupMessageEvent) {
		log.Printf("收到群消息: %s (group: %s)", e.Content, e.GroupOpenID)
		resp, err := client.API.SendGroupMessage(e.GroupOpenID, "Echo: "+e.Content, e.ID)
		if err != nil {
			log.Printf("❌ 发送群消息失败: %v", err)
		} else {
			log.Printf("✅ 消息发送成功: id=%s", resp.ID)
		}
	})

	client.Gateway.OnChannelMessage(func(e qqbot.ChannelMessageEvent) {
		log.Printf("Channel: %s", e.Content)
		if _, err := client.API.SendChannelMessage(e.ChannelID, "Echo: "+e.Content, e.ID); err != nil {
			log.Printf("发送频道消息失败: %v", err)
		}
	})

	if err := client.Connect(); err != nil {
		log.Fatalf("连接失败: %v", err)
	}
	log.Println("Bot started")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	client.Gateway.Close()
}
