package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nats-io/nats.go"
	"log"
	"os"
	"strconv"
	"sync"
)

type App struct {
	NC     *nats.Conn
	TGBot  *tgbotapi.BotAPI
	ChatId int64
}

func (a *App) Initialize() {
	var err error
	a.NC, err = nats.Connect(os.Getenv("NATS_URL"))
	if err != nil {
		log.Fatalf("connecting to nats server failed: %s", err)
	}
	log.Printf("subscriber started %v, %s, %s", a.NC.Opts, os.Getenv("NATS_CHANNEL"), os.Getenv("NATS_GROUP"))

	a.TGBot, err = tgbotapi.NewBotAPI(os.Getenv("API_KEY"))
	if err != nil {
		log.Fatalf("initializing TGBot failed: %s", err)
	}
	a.TGBot.Debug = true
	log.Printf("Authorized on account %s", a.TGBot.Self.UserName)

	a.ChatId, err = strconv.ParseInt(os.Getenv("CHAT_ID"), 10, 64)
	if err != nil {
		log.Fatalf("parsing chat id failed: %s", err)
	}
	log.Printf("broadcaster is alive")
}

func (a *App) Run() {

	defer a.NC.Close()

	a.NC.QueueSubscribe(os.Getenv("NATS_CHANNEL"), os.Getenv("NATS_GROUP"), a.broadcast2TG)
	wg := sync.WaitGroup{}

	wg.Add(1)
	wg.Wait()
}
