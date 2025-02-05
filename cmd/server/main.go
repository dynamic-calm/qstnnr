package main

import (
	"context"
	"fmt"
	"os"

	"github.com/mateopresacastro/qstnnr"
)

func getenv(key string) string {
	switch key {
	case "PORT":
		return "5974"
	default:
		return ""
	}
}

func main() {
	ctx := context.Background()
	if err := qstnnr.Run(ctx, getenv); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
