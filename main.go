package main

import (
	nats "github.com/nats-io/nats.go"
	"log"
	"os"
	"sync"
)

func broadcast2TG(m *nats.Msg) {
	log.Printf("nats message: %s", m.Data)
}

func main() {

	nc, err := nats.Connect(os.Getenv("NATS_URL"))
	defer nc.Close()
	if err != nil {
		log.Fatalf("connecting to nats server failed: %s", err)
	}
	nc.QueueSubscribe(os.Getenv("NATS_CHANNEL"), os.Getenv("NATS_GROUP"), broadcast2TG)
	wg := sync.WaitGroup{}

	wg.Add(1)
	wg.Wait()

}
