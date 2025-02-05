package main

import (
	"context"
	"fmt"
	"os"

	"github.com/mateopresacastro/qstnnr"
)

// This function could be replace for the native
// os.Getenv. This is usually done for testing but
// decdided to add it here for ease of use.
func getenv(key string) string {
	switch key {
	case "PORT":
		return "5974"
	case "LOG_LEVEL":
		return "INFO"
	default:
		return ""
	}
}

// Inspired by: https://grafana.com/blog/2024/02/09/how-i-write-http-services-in-go-after-13-years/
func main() {
	ctx := context.Background()
	if err := qstnnr.Run(ctx, getenv, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
