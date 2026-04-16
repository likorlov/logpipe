# logpipe

Lightweight structured log aggregator with pluggable output sinks for Go services.

## Installation

```bash
go get github.com/yourorg/logpipe
```

## Usage

```go
package main

import (
    "github.com/yourorg/logpipe"
)

func main() {
    // Create a new logger with a stdout sink
    logger := logpipe.New(
        logpipe.WithSink(logpipe.StdoutSink()),
        logpipe.WithSink(logpipe.FileSink("app.log")),
    )
    defer logger.Close()

    // Log structured messages
    logger.Info("server started", logpipe.Fields{
        "port": 8080,
        "env":  "production",
    })

    logger.Error("request failed", logpipe.Fields{
        "status": 500,
        "path":   "/api/users",
    })
}
```

### Custom Sinks

Implement the `Sink` interface to pipe logs to any destination:

```go
type Sink interface {
    Write(entry logpipe.Entry) error
    Close() error
}
```

## Features

- Structured JSON logging out of the box
- Pluggable output sinks (stdout, file, HTTP, custom)
- Minimal allocations and low overhead
- Concurrent-safe log aggregation

## License

MIT — see [LICENSE](LICENSE) for details.