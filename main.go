package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/techies/orders-api/application"

	// Load .env if present (for local dev)
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	cfg := application.DefaultConfig()
	app := application.NewApp(cfg)
	ctx, can := signal.NotifyContext(context.Background(), os.Interrupt)
	defer can()

	if err := app.Start(ctx); err != nil {
		_, err := fmt.Fprintf(os.Stderr, "failed to start the app: %v\n", err)
		if err != nil {
			return
		}
		os.Exit(1)
	}
}
