# GopenAI - golang open ai client

basic usage:

```go
package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/blackestwhite/gopenai"
	"golang.org/x/net/proxy"
)

func main() {
	key := "YOUR-OPEN-AI-KEY"

	instance := gopenai.Setup(key)

	p := gopenai.ChatCompletionRequestBody{
		Model: "gpt-3.5-turbo",
		Messages: []gopenai.Message{
			{Role: "user", Content: "hi"},
		},
		Stream: true,
	}

	resultCh, err := instance.GenerateChatCompletion(p)
	if err != nil {
		log.Fatal(err)
	}

	for chunk := range resultCh {
		log.Println(chunk)
	}
}
```

with custom http client

```go
package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/blackestwhite/gopenai"
	"golang.org/x/net/proxy"
)

func main() {
	key := "YOUR-OPEN-AI-KEY"

    // open ai is blocked in my country so i use socks5 proxy to consume it
	transport := &http.Transport{}
	dialer, err := proxy.SOCKS5("tcp", "localhost:8586", nil, proxy.Direct)
	if err != nil {
		log.Fatal(err)
	}
	transport.DialContext = func(ctx context.Context, network, address string) (net.Conn, error) {
		return dialer.Dial(network, address)
	}

	client := &http.Client{Transport: transport}

	instance := gopenai.SetupCustom(key, client)

	p := gopenai.ChatCompletionRequestBody{
		Model: "gpt-3.5-turbo",
		Messages: []gopenai.Message{
			{Role: "user", Content: "hi"},
		},
		Stream: true,
	}

	resultCh, err := instance.GenerateChatCompletion(p)
	if err != nil {
		log.Fatal(err)
	}

	for chunk := range resultCh {
		log.Println(chunk)
	}
}
```