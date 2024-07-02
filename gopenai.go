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

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func Setup(openAiKey string) *GopenAiInstance {
	return SetupCustom(openAiKey, &http.Client{})
}

func SetupCustom(openAiKey string, client HTTPClient) *GopenAiInstance {
	return &GopenAiInstance{
		Client: client,
		key:    openAiKey,
	}
}

func (c *GopenAiInstance) GenerateChatCompletion(prompt ChatCompletionRequestBody) (*ChatCompletion, error) {
	prompt.Stream = false
	return c.sendChatCompletionRequest(prompt)
}

func (c *GopenAiInstance) GenerateChatCompletionStream(prompt ChatCompletionRequestBody) (chan ChatCompletionChunk, chan error) {
	prompt.Stream = true
	marshalled, err := json.Marshal(prompt)
	if err != nil {
		return nil, sendError(err)
	}

	resultCh := make(chan ChatCompletionChunk)
	errCh := make(chan error, 1)

	go c.streamChatCompletion(marshalled, resultCh, errCh)

	return resultCh, errCh
}

func (c *GopenAiInstance) GenerateImage(prompt string) (ImageGenerationResponse, error) {
	if len(prompt) == 0 {
		return ImageGenerationResponse{}, errors.New("prompt length must be more than 0 chars")
	}

	imageGenerationRequestBody := ImageGenerationRequestBody{
		Model:  "dall-e-3",
		Count:  1,
		Size:   "1024x1024",
		Prompt: prompt,
	}

	return c.sendImageGenerationRequest(imageGenerationRequestBody)
}

func (c *GopenAiInstance) sendChatCompletionRequest(prompt ChatCompletionRequestBody) (*ChatCompletion, error) {
	marshalled, err := json.Marshal(prompt)
	if err != nil {
		return nil, err
	}

	req, err := c.createPostRequest("https://api.openai.com/v1/chat/completions", marshalled)
	if err != nil {
		return nil, err
	}

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

func (c *GopenAiInstance) streamChatCompletion(marshalled []byte, resultCh chan ChatCompletionChunk, errCh chan error) {
	req, err := c.createPostRequest("https://api.openai.com/v1/chat/completions", marshalled)
	if err != nil {
		errCh <- err
		return
	}

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
}

func (c *GopenAiInstance) sendImageGenerationRequest(body ImageGenerationRequestBody) (ImageGenerationResponse, error) {
	marshalled, err := json.Marshal(body)
	if err != nil {
		return ImageGenerationResponse{}, err
	}

	req, err := c.createPostRequest("https://api.openai.com/v1/images/generations", marshalled)
	if err != nil {
		return ImageGenerationResponse{}, err
	}

	res, err := c.Client.Do(req)
	if err != nil {
		return ImageGenerationResponse{}, fmt.Errorf("error sending HTTP request: %v", err)
	}
	defer res.Body.Close()

	var response ImageGenerationResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return ImageGenerationResponse{}, fmt.Errorf("error decoding JSON response: %v", err)
	}

	return response, nil
}

func (c *GopenAiInstance) createPostRequest(url string, body []byte) (*http.Request, error) {
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP request: %v", err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.key))
	req.Header.Add("Content-Type", "application/json")
	return req, nil
}

func sendError(err error) chan error {
	errCh := make(chan error, 1)
	errCh <- err
	close(errCh)
	return errCh
}
