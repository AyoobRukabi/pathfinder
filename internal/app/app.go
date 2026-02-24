package app

import (
	"log/slog"

	"gitea.kood.tech/ivanandreev/pathfinder/internal/config"
	"gitea.kood.tech/ivanandreev/pathfinder/internal/storage/local"
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

	app.log.Info("starting pathfinder app", slog.String("env", app.cfg.Env))
	app.log.Debug("debug messages are enabled")

	// storage
	storage := local.New(
		app.log,
		app.cfg.NetworkMapPath,
		app.cfg.StartStation,
		app.cfg.EndStation,
	)

	// service

	// graceful shutdown

	return nil
}
