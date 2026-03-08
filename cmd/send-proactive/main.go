package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/GoGoGoDoge/qqbot"
)

func main() {
	to := flag.String("to", "", "Target openid")
	text := flag.String("text", "", "Message text")
	typ := flag.String("type", "c2c", "Message type (c2c/group)")
	list := flag.Bool("list", false, "List known users")
	flag.Parse()

	appID := os.Getenv("QQBOT_APP_ID")
	secret := os.Getenv("QQBOT_CLIENT_SECRET")
	if appID == "" || secret == "" {
		log.Fatal("QQBOT_APP_ID and QQBOT_CLIENT_SECRET required")
	}

	client := qqbot.NewClient(appID, secret)
	store := qqbot.NewProactiveStore(".")

	if *list {
		users := store.List(*typ, "", 0)
		fmt.Printf("Known users (%d):\n", len(users))
		for _, u := range users {
			fmt.Printf("%s\t%s\t%s\n", u.Type, u.OpenID, u.Nickname)
		}
		return
	}

	if *to == "" || *text == "" {
		flag.Usage()
		os.Exit(1)
	}

	if _, err := client.API.GetAccessToken(); err != nil {
		log.Fatal(err)
	}

	var err error
	if *typ == "c2c" {
		_, err = client.API.SendC2CMessage(*to, *text, "")
	} else {
		_, err = client.API.SendGroupMessage(*to, *text, "")
	}

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Message sent successfully")
}
