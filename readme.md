# GopenAI - golang open ai client

## Features
- Supports streaming/SSE for chat completions
- Image generation using OpenAI's DALL-E

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

## image generation
```go
package main

import (
    "log"
    "github.com/blackestwhite/gopenai"
)

func main() {
    key := "YOUR-OPEN-AI-KEY"

    instance := gopenai.Setup(key)

    prompt := "a cute persian cat"

    res, _ := instance.GenerateImage(prompt)
    println(res.Data[0].URL) // prints the url of the generated image
}
```

## Donations

ETH: `blackestwhite.eth`

TON: `blackestwhite.ton`