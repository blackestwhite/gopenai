package gopenai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func Setup(openAiKey string) *GopenAiInstance {
	instance := &GopenAiInstance{
		Client: &http.Client{},
		key:    openAiKey,
	}
	return instance
}

func (h *GopenAiInstance) GenerateChatCompletion(prompt ChatCompletionRequestBody) (<-chan ChatCompletionChunk, error) {
	marshalled, err := json.Marshal(prompt)
	if err != nil {
		return nil, err
	}

	resultCh := make(chan ChatCompletionChunk, 1)

	go func() {
		req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewReader(marshalled))
		if err != nil {
			close(resultCh)
			log.Printf("Error creating HTTP request: %v", err)
			return
		}
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", h.key))
		req.Header.Add("Content-Type", "application/json")
		req.Header.Set("Cache-Control", "no-cache")
		req.Header.Set("Accept", "text/event-stream")
		req.Header.Set("Connection", "keep-alive")
		res, err := h.Client.Do(req)
		if err != nil {
			close(resultCh)
			log.Printf("Error sending HTTP request: %v", err)
			return
		}
		defer res.Body.Close()

		for {
			bufferSize := 1024
			data := make([]byte, bufferSize)
			n, err := res.Body.Read(data)
			if err != nil {
				if n == 0 {
					// End of stream
					break
				}
				close(resultCh)
				log.Printf("Error reading from response body: %v", err)
				return
			}

			var chunk ChatCompletionChunk
			err = json.Unmarshal(data[:n], &chunk)
			if err != nil {
				close(resultCh)
				log.Printf("Error unmarshalling JSON: %v", err)
				return
			}

			resultCh <- chunk

			if chunk.Choices[0].FinishReason == "stop" {
				break
			}
		}
	}()

	return resultCh, nil
}
