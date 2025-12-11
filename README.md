# goNextCore

`goNextCore` is a robust, modular, and feature-rich Go framework designed for building scalable backend applications. It provides a solid foundation with built-in support for dependency injection, configuration management, async processing, and more.

## Features

-   **Dependency Injection**: Built-in container for managing dependencies and lifecycle.
-   **Modular Architecture**: Component-based design for cleaner code organization.
-   **Async Processing**: Utilities for async/await patterns and promise-like execution.
-   **Job Queues**: Integration with `hibiken/asynq` for background job processing.
-   **Task Scheduling**: Cron-based task scheduling.
-   **Security**: JWT authentication service and guards.
-   **Database & Caching**: Helpers for Postgres (pgx) and Redis.

## Installation

```bash
go get github.com/Alexigbokwe/goNextCore
```

## Usage

### Basic Application Setup

```go
package main

import (
    "github.com/Alexigbokwe/goNextCore/core"
)

func main() {
    // Create a new app instance
    app := core.NewApp()

    // Create a container
    container := core.NewContainer()
    
    // Register your modules and services
    // container.Register(...)

    // Start the application
    if err := app.Start(); err != nil {
        panic(err)
    }
}
```

### Async/Await Example

```go
import (
    "context"
    "fmt"
    "github.com/Alexigbokwe/goNextCore/core/utils"
)

func main() {
    // Create an async task
    future := utils.Async(func(ctx context.Context) (string, error) {
        // Do some work...
        return "Result", nil
    })

    // Await the result
    result, err := future.Await()
    if err != nil {
        panic(err)
    }
    fmt.Println(result)
}
```

## Modules

-   **Core**: The heart of the framework, handling lifecycle and DI.
-   **Security**: JWT utilities and authentication guards.
-   **Utils**: Helper functions for async operations and more.

## License

MIT
