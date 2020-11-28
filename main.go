package main

import (
	"fmt"
	"github.com/line/line-bot-sdk-go/linebot"
	"log"
	"net/http"
	"os"
	"strconv"

	_ "github.com/joho/godotenv/autoload"
)

var (
	client *linebot.Client
	err    error
)

func main() {
	client, err = linebot.New(os.Getenv("CHANNEL_SECRET"), os.Getenv("CHANNEL_ACCESS_TOKEN"))

	if err != nil {
		log.Println(err.Error())
	}

	http.HandleFunc("/callback", callbackHandler)

	log.Fatal(http.ListenAndServe(":84", nil))
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	events, err := client.ParseRequest(r)

	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}
		return
	}

	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				quota, err := client.GetMessageQuota().Do()
				if err != nil {
					log.Println(err.Error())
					return
				}

				text := fmt.Sprintf("Received Message: %s (%s)", message.Text, strconv.FormatInt(quota.Value, 10))

				if _, err = client.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(text)).Do(); err != nil {
					log.Println(err.Error())
				}
			}
		}
	}
}