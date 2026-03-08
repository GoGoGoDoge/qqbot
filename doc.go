// Package qqbot provides a Go client for the QQ Bot API.
//
// Basic usage:
//
//	client := qqbot.NewClient(appID, clientSecret)
//	client.Gateway.OnC2CMessage(func(e qqbot.C2CMessageEvent) {
//		// Handle message
//	})
//	client.Connect()
package qqbot
