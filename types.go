package gopenai

const (
	ModelGPT4o       = "gpt-4o"
	ModelGPT4Turbo   = "gpt-4-turbo"
	ModelGPT4        = "gpt-4" // continuous model upgrades, which points to other models
	ModelGPT4__32k   = "gpt-4-32k"
	ModelGPT3_5Turbo = "gpt-3.5-turbo" // points to gpt-3.5-turbo-0125
)

type ChatCompletion struct {
	ID                string   `json:"id"`
	Object            string   `json:"object"`
	Created           int64    `json:"created"`
	Model             string   `json:"model"`
	SystemFingerprint string   `json:"system_fingerprint"`
	Choices           []Choice `json:"choices"`
	Usage             Usage    `json:"usage"`
}

type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type ChatCompletionChunk struct {
	ID                string          `json:"id"`
	Object            string          `json:"object"`
	Created           int64           `json:"created"`
	Model             string          `json:"model"`
	SystemFingerprint string          `json:"system_fingerprint"`
	Choices           []ChunkedChoice `json:"choices"`
}

type ChunkedChoice struct {
	Index        int    `json:"index"`
	Delta        Delta  `json:"delta"`
	FinishReason string `json:"finish_reason"`
}

type Delta struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatCompletionRequestBody struct {
	Messages []Message `json:"messages"`
	Model    string    `json:"model"`
	Stream   bool      `json:"stream"`
}

type GopenAiInstance struct {
	Client HTTPClient
	key    string
}

type ImageGenerationRequestBody struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Count  int    `json:"n"`
	Size   string `json:"size"`
}

type ImageGenerationResponse struct {
	Created int64                         `json:"created"`
	Data    []ImageGenerationResponseData `json:"data"`
}

type ImageGenerationResponseData struct {
	RevisedPrompt string `json:"revised_prompt"`
	URL           string `json:"url"`
}
