package main

import (
	"log"
	"net/http"
	"time"
)

var mongoCtx MongoContext = MongoContext{}

func main() {
	mongoCtx.InitMongo()
	// mongoCtx.AddPrompt(Prompt{Query: "123"})

	go SpawnServer()

	client := &http.Client{}
	active_queries := make(map[string]void)
	active_queries["%23elections"] = nullptr

	for {
		prompts := mongoCtx.GetAllPrompts()
		for _, prompt := range prompts {
			if prompt.IsActive {
				log.Printf("Running query '%s' with last read '%f'", prompt.Query, prompt.LastIdRead)
				lastReadId := RunQuery(prompt.Query, prompt.LastIdRead, client)
				mongoCtx.UpdatePromptLastRead(prompt.Query, lastReadId)
			}
		}
		time.Sleep(10 * time.Second)
	}

}
