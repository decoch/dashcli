package main

import (
	"context"
	"io"
	"os"

	"github.com/decoch/dashcli/internal/app"
)

var version = "dev"

func runMain(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	app.Version = version
	return app.Run(ctx, args, stdout, stderr)
}

func main() {
	os.Exit(runMain(context.Background(), os.Args[1:], os.Stdout, os.Stderr))
}

