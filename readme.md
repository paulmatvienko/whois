# Go WHOIS Client Library

[![Go Reference](https://pkg.go.dev/badge/github.com/paulmatvienko/whois.svg)](https://pkg.go.dev/github.com/paulmatvienko/whois)
[![GitHub release](https://img.shields.io/github/v/release/paulmatvienko/whois)](https://github.com/paulmatvienko/whois/releases)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Coverage Status](https://coveralls.io/repos/github/paulmatvienko/whois/badge.svg?branch=main)](https://coveralls.io/github/paulmatvienko/whois?branch=main)

Go library for WHOIS queries with referral support.

## Features

- 500+ TLD support (automatic WHOIS server selection)
- Recursive queries with referral following
- Timeouts and cancellation via context.Context

## Installation

```bash
go get github.com/paulmatvienko/whois
```

## Usage

### Basic Example

```go
package main

import (
    "context"
    "fmt"
    "time"
    "github.com/paulmatvienko/whois"
)

func main() {
    client, _ := whois.New(whois.Options{
        Timeout: 5 * time.Second,
    })
    
    result, _ := client.Lookup(context.Background(), "example.com")
    fmt.Println(result.RawData)
}
```

## License

MIT - See LICENSE for details.