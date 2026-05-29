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
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("fail to load the configuration: %w", err)
	}

	err = cfg.Validate()
	if err != nil {
		log.Fatal("config validation failed: %w", err)
	}

	logger, err := logger.New(cfg.App.LogLevel, cfg.App.LogFormat)
	if err != nil {
		log.Fatal("fail to load logger: %w", err)
	}
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	app := app.New(cfg, logger)
	if err := app.Run(ctx); err != nil {
		log.Fatal("could not start the application %w", err)
		return 1
	}
	log.Println("Gateway stopped")
	return 0
}
