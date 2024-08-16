package main

import (
	"log"
	"os"

	"github.com/blackestwhite/gopenai"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("./example/.env")
	if err != nil {
		log.Fatal(err)
	}

	openAIKey := os.Getenv("OPEN_AI_KEY")

	instance := gopenai.Setup(openAIKey)

	p := gopenai.ChatCompletionRequestBody{
		Model: gopenai.ModelGPT4oMini,
		Messages: []gopenai.Message{
			{Role: "user", Content: "hi"},
		},
		Stream: true,
	}

	resultCh, errCh := instance.GenerateChatCompletionStream(p)
	shouldBreak := false
	for !shouldBreak {
		select {
		case chunk, ok := <-resultCh:
			if !ok {
				resultCh = nil
			} else {
				if len(chunk.Choices) > 0 && chunk.Choices[0].FinishReason == "stop" {
					shouldBreak = true
				}
				log.Println(chunk)
			}
		case err, ok := <-errCh:
			if ok {
				log.Printf("Error: %v", err)
			}
			errCh = nil
		}

		if resultCh == nil && errCh == nil {
			break
		}
	}
}
