package app

import (
	"context"
	"log"

	"github.com/PaingPhyoAungKhant/url-shortener/internal/server"
)

type App struct{}

func New() *App {
	return &App{}
}

func (a *App) Run(ctx context.Context) error {
	log.Println("Starting the application...")
	srv := server.New()
	return srv.Start(ctx)
}
