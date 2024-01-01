package gopenai

import (
	"bufio"
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

func SetupCustom(openAiKey string, client *http.Client) *GopenAiInstance {
	instance := &GopenAiInstance{
		Client: client,
		key:    openAiKey,
	}
	return instance
}

func (h *GopenAiInstance) GenerateChatCompletion(prompt ChatCompletionRequestBody) (chan ChatCompletionChunk, error) {
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

		res, err := h.Client.Do(req)
		if err != nil {
			close(resultCh)
			log.Printf("Error sending HTTP request: %v", err)
			return
		}
		defer res.Body.Close()

		scanner := bufio.NewScanner(res.Body)
		for scanner.Scan() {
			line := scanner.Text()
			if line == "" {
				continue
			}

			var chunk ChatCompletionChunk
			err = json.Unmarshal([]byte(line)[6:], &chunk)
			if err != nil {
				log.Printf("Error unmarshalling JSON: %v", err)
				break
			}

			resultCh <- chunk

			if chunk.Choices[0].FinishReason == "stop" {
				res.Body.Close()
				break
			}
		}

		if err := scanner.Err(); err != nil {
			log.Fatalf("Error reading response body: %v", err)
		}

		close(resultCh)
	}()

	return resultCh, nil
}
