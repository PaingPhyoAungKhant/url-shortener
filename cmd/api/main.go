package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/PaingPhyoAungKhant/url-shortener/internal/app"
	"github.com/PaingPhyoAungKhant/url-shortener/internal/config"
	"github.com/PaingPhyoAungKhant/url-shortener/internal/logger"
)

func main() {
	os.Exit(run())
}

func run() int {
	cfg := config.Load()

	err := cfg.Validate()
	if err != nil {
		log.Fatalf("config validation failed: %v", err)
	}

	logger, err := logger.New(cfg.App.LogLevel, cfg.App.LogFormat)
	if err != nil {
		log.Fatalf("fail to load logger: %v", err)
	}
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	app := app.New(cfg, logger)
	if err := app.Run(ctx); err != nil {
		log.Fatalf("could not start the application: %v", err)
	}
	log.Println("Gateway stopped")
	return 0
}
