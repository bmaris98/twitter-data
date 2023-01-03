package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var mongoCtx MongoContext = MongoContext{}
var kafkaCtx KafkaContext = KafkaContext{}

func main() {
	mongoCtx.InitMongo()
	kafkaCtx.Init(mongoCtx)
	go kafkaCtx.StartReading()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		fmt.Println(sig)
		fmt.Println("Shutting down")
		kafkaCtx.CloseReader()
		kafkaCtx.CloseWriter()
	}()

	go SpawnServer()
	client := &http.Client{}

	for {
		prompts := mongoCtx.GetAllPrompts()
		for _, prompt := range prompts {
			if prompt.IsActive {
				log.Printf("Running query '%s' with last read '%f'", prompt.Query, prompt.LastIdRead)
				lastReadId := RunQuery(prompt.Query, prompt.LastIdRead, client, kafkaCtx)
				mongoCtx.UpdatePromptLastRead(prompt.Query, lastReadId)
			}
		}
		time.Sleep(10 * time.Second)
	}

}
