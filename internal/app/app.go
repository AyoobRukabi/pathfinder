package app

import (
	"log/slog"

	"gitea.kood.tech/ivanandreev/pathfinder/internal/config"
	"gitea.kood.tech/ivanandreev/pathfinder/pkg/logger"
)

type App struct {
	cfg *config.Config
	log *slog.Logger
	// some other fields
}

func New(cfg *config.Config) *App {
	return &App{
		cfg: cfg,
		log: logger.New(cfg.Env),
	}
}

func (app *App) Run() error {
	// main program wiering file

	// storage

	// service

	return nil
}
