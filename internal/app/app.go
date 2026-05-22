package app

import (
	"context"
	"log"

	"github.com/PaingPhyoAungKhant/url-shortener/internal/config"
	"github.com/PaingPhyoAungKhant/url-shortener/internal/server"
)

type App struct {
	cfg *config.Config
}

func New(cfg *config.Config) *App {
	return &App{cfg}
}

func (a *App) Run(ctx context.Context) error {
	log.Println("Starting the application...")
	srv := server.New(a.cfg)
	return srv.Start(ctx)
}
