package app

import (
	"context"

	"github.com/PaingPhyoAungKhant/url-shortener/internal/config"
	"github.com/PaingPhyoAungKhant/url-shortener/internal/logger"
	"github.com/PaingPhyoAungKhant/url-shortener/internal/server"
)

type App struct {
	cfg *config.Config
	log *logger.Logger
}

func New(cfg *config.Config, log *logger.Logger) *App {
	return &App{cfg: cfg, log: log}
}

func (a *App) Run(ctx context.Context) error {
	a.log.Info(ctx, "Starting the application...")
	srv := server.New(a.cfg, a.log)
	return srv.Start(ctx)
}
