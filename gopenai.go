package gopenai

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
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

func (c *GopenAiInstance) GenerateChatCompletion(prompt ChatCompletionRequestBody) (*ChatCompletion, error) {
	prompt.Stream = false
	marshalled, err := json.Marshal(prompt)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewReader(marshalled))
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP request: %v", err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.key))
	req.Header.Add("Content-Type", "application/json")

	res, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending HTTP request: %v", err)
	}
	defer res.Body.Close()

	var completeResult ChatCompletion
	if err := json.NewDecoder(res.Body).Decode(&completeResult); err != nil {
		return nil, fmt.Errorf("error decoding JSON response: %v", err)
	}

	return &completeResult, nil
}

func (c *GopenAiInstance) GenerateChatCompletionStream(prompt ChatCompletionRequestBody) (chan ChatCompletionChunk, chan error) {
	prompt.Stream = true
	marshalled, err := json.Marshal(prompt)
	if err != nil {
		errCh := make(chan error, 1)
		errCh <- err
		close(errCh)
		return nil, errCh
	}

	resultCh := make(chan ChatCompletionChunk)
	errCh := make(chan error, 1)

	go func() {
		defer close(resultCh)
		defer close(errCh)

		req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewReader(marshalled))
		if err != nil {
			errCh <- fmt.Errorf("error creating HTTP request: %v", err)
			return
		}
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.key))
		req.Header.Add("Content-Type", "application/json")

		res, err := c.Client.Do(req)
		if err != nil {
			errCh <- fmt.Errorf("error sending HTTP request: %v", err)
			return
		}
		defer res.Body.Close()

		scanner := bufio.NewScanner(res.Body)
		for scanner.Scan() {
			line := scanner.Text()
			if line == "" {
				continue
			}

			if line == "data: [DONE]" {
				break
			}

			cleanLine := regexp.MustCompile(`^data: `).ReplaceAllString(line, "")
			var chunk ChatCompletionChunk
			err := json.Unmarshal([]byte(cleanLine), &chunk)
			if err != nil {
				errCh <- fmt.Errorf("error unmarshalling JSON: %v", err)
				return
			}

			resultCh <- chunk

			if chunk.Choices[0].FinishReason == "stop" {
				break
			}
		}

		if err := scanner.Err(); err != nil {
			errCh <- fmt.Errorf("error reading response body: %v", err)
		}
	}()

	return resultCh, errCh
}

func (c *GopenAiInstance) GenerateImage(prompt string) (response ImageGenerationResponse, err error) {
	if len(prompt) < 1 {
		return response, errors.New("prompt length must be more than 0 chars")
	}
	imageGenerationRequestBody := ImageGenerationRequestBody{
		Model:  "dall-e-3",
		Count:  1,
		Size:   "1024x1024",
		Prompt: prompt,
	}
	marshalled, err := json.Marshal(imageGenerationRequestBody)
	if err != nil {
		return response, err
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/images/generations", bytes.NewReader(marshalled))
	if err != nil {
		return response, fmt.Errorf("error creating HTTP request: %v", err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.key))
	req.Header.Add("Content-Type", "application/json")

	res, err := c.Client.Do(req)
	if err != nil {
		return response, fmt.Errorf("error sending HTTP request: %v", err)
	}
	defer res.Body.Close()

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return response, fmt.Errorf("error decoding JSON response: %v", err)
	}

	return response, nil
}
