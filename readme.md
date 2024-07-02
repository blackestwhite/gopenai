# GopenAI - golang open ai client

## Features
- supports stream/sse

## How to use

basic usage(with SSE):

```go
package main

import (
    "log"
    "github.com/blackestwhite/gopenai"
)

func main() {
    key := "YOUR-OPEN-AI-KEY"

    instance := gopenai.Setup(key)

    p := gopenai.ChatCompletionRequestBody{
        Model: gopenai.ModelGPT3_5Turbo,
        Messages: []gopenai.Message{
            {Role: "user", Content: "hi"},
        },
    }

    resultCh, errCh := instance.GenerateChatCompletionStream(p)

    for {
        select {
        case chunk, ok := <-resultCh:
            if !ok {
                resultCh = nil
            } else {
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
```

## available models
- gopenai.ModelGPT4o
- gopenai.ModelGPT4Turbo
- gopenai.ModelGPT4
- gopenai.ModelGPT4__32k
- gopenai.ModelGPT3_5Turbo

## Donations

ETH: `blackestwhite.eth`

TON: `blackestwhite.ton`