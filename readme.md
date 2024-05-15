# GopenAI - golang open ai client

## Features
- supports stream/sse

## How to use

basic usage:

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
        Model: "gpt-3.5-turbo",
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

## Donations

ETH: blackestwhite.eth

TON: blackestwhite.ton