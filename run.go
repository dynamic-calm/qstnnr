package qstnnr

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"sync"

	"google.golang.org/grpc"
)

const defaultPort = "5974"

func Run(
	ctx context.Context,
	getenv func(string) string,
	output io.Writer,
) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	questions := map[QuestionID]Question{
		1: {
			ID:   1,
			Text: "What is the capital of France?",
			Options: map[OptionID]Option{
				1: {ID: 1, Text: "London"},
				2: {ID: 2, Text: "Paris"},
				3: {ID: 3, Text: "Berlin"},
				4: {ID: 4, Text: "Madrid"},
			},
		},
		2: {
			ID:   2,
			Text: "Which planet is known as the Red Planet?",
			Options: map[OptionID]Option{
				1: {ID: 1, Text: "Venus"},
				2: {ID: 2, Text: "Mars"},
				3: {ID: 3, Text: "Jupiter"},
				4: {ID: 4, Text: "Saturn"},
			},
		},
		3: {
			ID:   3,
			Text: "What is 2 + 2?",
			Options: map[OptionID]Option{
				1: {ID: 1, Text: "3"},
				2: {ID: 2, Text: "4"},
				3: {ID: 3, Text: "5"},
				4: {ID: 4, Text: "6"},
			},
		},
	}

	solutions := map[QuestionID]OptionID{
		1: 2, // Paris
		2: 2, // Mars
		3: 2, // 4
	}

	data := InitialData{
		Questions: questions,
		Solutions: solutions,
	}

	store, err := NewMemoryStore(data)
	if err != nil {
		return err
	}

	service := NewQstnnrService(store)

	logger := slog.New(slog.NewJSONHandler(output, &slog.HandlerOptions{
		Level: parseLogLevel(getenv("LOG_LEVEL")),
	}))

	cfg := &ServerConfig{
		Logger:  logger,
		Service: service,
	}

	server, err := NewServer(cfg)

	port := getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	ln, err := net.Listen("tcp", ":"+port)

	go func() {
		logger.Info("listening", "port", getenv("PORT"))
		if err := server.Serve(ln); err != nil && err != grpc.ErrServerStopped {
			fmt.Fprintf(os.Stderr, "error listening and serving: %s\n", err)
		}
	}()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		logger.Info("shutting down server")
		server.GracefulStop()
	}()

	wg.Wait()
	return nil
}

func parseLogLevel(level string) slog.Level {
	switch level {
	case "DEBUG":
		return slog.LevelDebug
	case "INFO":
		return slog.LevelInfo
	case "WARN":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
