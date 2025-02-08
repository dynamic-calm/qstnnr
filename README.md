# qstnnr - A Go Quiz Application

A command-line quiz application for a take home assignment built with `Go`, featuring a `gRPC` API server and CLI client. Users can take quizzes, get feedback, and compare their performance against other participants.

## Features

- `gRPC` API Server
- Command-line interface for taking the quiz
- In-Memory Storage
- Performance comparison with other participants

## Technical Stack

- Backend: `Go`
- `API`: `gRPC` with `Protocol Buffers`
- Storage: In-memory
- CLI Framework: `Cobra`
- Error Handling: Domain-specific error types with stack traces
- Testing: Unit and integration tests

## Design Decisions

- `gRPC`: I prefer `gRPC` over `REST` for a couple reasons:
  - Development experience
  - Schema first approach. The `.proto` file are the source of truth. In a sense, similar to `GraphQL`.
  - Language agnostic
  - Performance
  - Growing ecosystem
  - Used by `Kubernetes`, `Etcd`, `Cockroach DB`, etc.
- Design:
  - The `API` or server layer (`pkg/server` and `pkg/api` [proto]) processes the requests data to and from the service layer.
  - Service layer (`pkg/qservice`) has te business logic and queries the `store`.
  - Store layer interacts with the in-memory database
  - `run.go` starts off the `server`.
  - Minimal `main` functions.
- Testing:
  - Tests per package
  - Integration test
  - E2E test
- Errors:

  - For the error design I opted to follow the opinionated approach lied out in the book [Concurrency in Go by Katherine Cox-Buday](https://www.oreilly.com/library/view/concurrency-in-go/9781491941294/).
  - If an error is a known edge case, you return the error wrapped on a custom error type, if not, you return the error as is. E.g:

    ```go
        solutions, err := qs.store.Solutions()
        if err != nil {
            if _, ok := err.(store.StoreError); !ok {
                // Bug
                return nil, err
            }
        // Known edge case. Wrap it and return.
        return nil, ServiceError{qerr.Wrap(err, qerr.Internal, "failed to get solutions")}
    }
    ```

- CLI:
  - `Cobra` as specified
  - `promptui` for interactivity

## Getting Started

### Prerequisites

- `Go`
- Optional: `Protocol Buffers` compiler (`protoc`) and `Make`.

### Installation

Clone the repository:

```bash
git clone https://github.com/mateopresacastro/qstnnr.git && cd qstnnr
```

Build the application:

```bash
make build
```

This will create two binaries in the `bin` directory:

`server`: The `gRPC` server

`qstnnr`: The CLI client

## Running the Application

To run the quiz run `qstnnr server start` and then `qstnnr take`:

```bash
➜ bin/qstnnr server start
Server started on port 4000 (PID: 12152)
```

```console
➜ bin/qstnnr take
Question 1 of 10
Use the arrow keys to navigate: ↓ ↑ → ←
What is the purpose of the blank identifier (_) in Go?
  ➜ To discard an unwanted value
    To declare a private variable
    To create an anonymous function
    To mark a variable as nullable
```

The CLI has two main commands: `server` and `take`.

```bash
➜ bin/qstnnr help
A CLI application to check you Go knowledge.

Usage:
  qstnnr [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  server      Manage the qstnnr server
  take        Take the quiz

Flags:
  -h, --help   help for qstnnr

Use "qstnnr [command] --help" for more information about a command.
```

### `server` command

To see the help for the server command run:

```bash
➜ bin/qstnnr server help
Commands to start, stop, restart and check status of the qstnnr server

Usage:
  qstnnr server [command]

Available Commands:
  restart     Restart the qstnnr server
  start       Start the qstnnr server
  status      Check the status of the qstnnr server
  stop        Stop the qstnnr server

Flags:
  -h, --help   help for server

Use "qstnnr server [command] --help" for more information about a command.
```

You can set `PORT` in your `env` vars before starting the server. By default it uses `5974`.

```bash
➜ bin/qstnnr server start
Server started on port 5974 (PID: 95853)
```

```bash
➜ export PORT=4000
➜ bin/qstnnr server start
Server started on port 4000 (PID: 96592)
```

## `take` command

The `take` command starts the quiz. At the end you can see your results.

```bash
➜ bin/qstnnr take
Question 1 of 10
Use the arrow keys to navigate: ↓ ↑
What is the zero value for a pointer in Go?
  ➜ 0
    undefined
    void
    nil
```

## Project Structure

```bash
├── cmd/ # Application entrypoints
│ ├── cli/ # CLI implementation
│ └── server/ # Server implementation
├── pkg/
│ ├── api/ # gRPC protocol definitions
│ ├── qerr/ # Error handling
│ ├── qservice/ # Business logic
│ ├── server/ # gRPC server implementation
│ └── store/ # Data storage
├── Makefile # Build and development commands
├── questions.go # Quiz content and initial data
└── run.go # Main application setup and server initialization
```

## Development

## Running

```bash
make start-server
```

```bash
make start-cli
```

### Building

```bash
make build # Build all binaries
```

```bash
make build-server # Build server only
```

```bash
make build-cli # Build CLI only
```

### Testing

```bash
make test # Run all tests with coverage and race detection
```

### Generating Protocol Buffers

```bash
make proto # Generate Go code from .proto files
```
