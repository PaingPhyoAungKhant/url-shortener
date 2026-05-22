package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/PaingPhyoAungKhant/url-shortener/internal/app"
)

func main() {
	os.Exit(run())
}

func run() int {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	app := app.New()
	if err := app.Run(ctx); err != nil {
		log.Fatal("could not start the application %w", err)
		return 1
	}
	log.Println("Gateway stopped")
	return 0
}
