package main

import (
	"context"
	"fmt"
	"os"

	"github.com/mateopresacastro/qstnnr"
)

// Inspired by: https://grafana.com/blog/2024/02/09/how-i-write-http-services-in-go-after-13-years/
func main() {
	ctx := context.Background()
	if err := qstnnr.Run(ctx, os.Getenv, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}