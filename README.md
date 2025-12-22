# SafeComms Go SDK

Official Go client for the SafeComms API.

SafeComms is a powerful content moderation platform designed to keep your digital communities safe. It provides real-time analysis of text to detect and filter harmful content, including hate speech, harassment, and spam.

**Get Started for Free:**
We offer a generous **Free Tier** for all users, with **no credit card required**. Sign up today and start protecting your community immediately.

## Documentation

For full API documentation and integration guides, visit [https://safecomms.dev/docs](https://safecomms.dev/docs).

## Installation

```bash
go get github.com/safecomms/safecomms-go
```

## Usage

```go
package main

import (
    "fmt"
    "log"
    "github.com/safecomms/safecomms-go"
)

func main() {
    client := safecomms.NewClient("your-api-key", "")

    // Moderate text
    result, err := client.ModerateText(safecomms.ModerateTextRequest{
        Content:  "Some text to check",
        Language: "en",
        Replace:  true,
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(result)

    // Get usage
    usage, err := client.GetUsage()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(usage)
}
```
