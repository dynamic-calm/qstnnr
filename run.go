package qstnnr

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"sync"

	"google.golang.org/grpc"
)

func Run(ctx context.Context, getenv func(string) string) error {
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

	logger := slog.Default()
	cfg := &ServerConfig{
		Logger:  logger,
		Service: service,
	}

	server, err := NewServer(cfg)
	ln, err := net.Listen("tcp", ":"+getenv("PORT"))

	go func() {
		log.Printf("listening on %s\n", getenv("PORT"))
		if err := server.Serve(ln); err != nil && err != grpc.ErrServerStopped {
			fmt.Fprintf(os.Stderr, "error listening and serving: %s\n", err)
		}
	}()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		server.GracefulStop()
	}()

	wg.Wait()
	return nil
}
