# SafeComms Go SDK

Official Go client for the SafeComms API.

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
