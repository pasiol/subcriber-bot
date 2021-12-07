package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nats-io/nats.go"
	"log"
	"time"
)

func (a *App) broadcast2TG(m *nats.Msg) {
	log.Printf("nats message: %s", m.Data)
	var msg = tgbotapi.NewMessage(a.ChatId, string(m.Data))

	for i := 1; i <= 10; i++ {
		send, err := a.TGBot.Send(msg)
		if err != nil {
			log.Printf("sending tg message failed: %s", err)
			time.Sleep(time.Duration(i*10) * time.Second)
		} else {
			log.Printf("sending tg message succeed: %v", send)
			break
		}
	}
}
