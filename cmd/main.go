package main

import (
	"log"
	"os"

	"github.com/koenno/aiapi/aiservice"
	"github.com/koenno/aiapi/server"

	"github.com/gin-gonic/gin"
	"github.com/koenno/aidevs2/ai"
)

func main() {
	openaiKey := os.Getenv("OPENAI_KEY")
	if openaiKey == "" {
		log.Fatalf("openai key not found")
	}

	router := gin.Default()

	AIEngine := ai.NewChat(openaiKey)
	AIModerator := ai.NewModerator(AIEngine)
	service := aiservice.New(AIEngine, AIModerator)
	s := server.New(service)
	s.Routes(router)
	err := router.Run()
	if err != nil {
		log.Fatalf("ERR: router failure: %v", err)
	}
}
