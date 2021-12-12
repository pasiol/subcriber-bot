package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nats-io/nats.go"
	"log"
	"os"
	"strconv"
	"sync"
	"time"
)

type App struct {
	NC     *nats.Conn
	TGBot  *tgbotapi.BotAPI
	ChatId int64
}

func (a *App) Initialize() {
	var err error
	for i := 1; i <= 10; i++ {
		a.NC, err = nats.Connect(os.Getenv("NATS_URL"), nats.Timeout(10*time.Second))
		if err != nil {
			log.Printf("connecting to nats server failed: %s, sleeping some time", err)
			time.Sleep(time.Duration(i*10) * time.Second)
		} else {
			log.Printf("subscriber started %v, %s, %s", a.NC.Opts, os.Getenv("NATS_CHANNEL"), os.Getenv("NATS_GROUP"))
			break
		}
	}

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
