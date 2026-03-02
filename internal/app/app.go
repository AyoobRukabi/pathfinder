package app

import (
	"fmt"
	"log/slog"

	"gitea.kood.tech/ivanandreev/pathfinder/internal/config"
	"gitea.kood.tech/ivanandreev/pathfinder/internal/service/sur"
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
	)

	//test
	if err := testFunc(storage); err != nil { // TODO: Remove
		return err
	}

	// service

	// graceful shutdown

	return nil
}

func testFunc(s *local.Storage) error {
	networkMap, err := s.BuildMap()
	if err != nil {
		return err
	}
	Service := sur.New(
		s.StartStation,
		s.EndStation,
		s.Trains,
		networkMap,
	)
	paths := Service.FindOptimalPaths()
	fmt.Println(paths)

	for i := range networkMap.AdjList {
		fmt.Printf("Node: %d connects to nodes: %v\n", i, networkMap.AdjList[i])
	}

	fmt.Printf("-------------------------------------\n")

	for i := range networkMap.AdjList {
		var tmpStations []string
		for _, s := range networkMap.AdjList[i] {
			tmpStations = append(tmpStations, networkMap.Stations[s].Name)
		}
		fmt.Printf("Station: %s has connection with: %v\n", networkMap.Stations[i].Name, tmpStations)
	}
	return nil
}
