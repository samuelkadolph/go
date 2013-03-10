# httpstream

httpstream is a library for connecting to an HTTP server that streams back information.

## Description

httpstream uses the net.http package to connect to a server and waits for data to be received and pushes it into a channel.

## Installation

    $ go get github.com/samuelkadolph/go/httpstream

## Usage

```go
package main

import (
  "fmt"
  "github.com/samuelkadolph/go/httpstream"
  "log"
  "net/url"
)

func main() {
  u := &url.URL{}
  u.Scheme = "https"
  u.Host = "streaming.campfirenow.com"
  u.Path = "/room/111111/live.json"

  s, err := httpstream.New(httpstream.Get{u}, httpstream.BasicAuth{"TTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTT", "x"}, httpstream.CarriageReturn)
  if err != nil {
    log.Fatalf("%s", err.Error())
  }

  for d := range s.Data {
    fmt.Printf("%s\n", string(d))
  }
}
```

## Contributing

Fork, branch & pull request.
