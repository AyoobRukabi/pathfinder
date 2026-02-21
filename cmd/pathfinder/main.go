package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"gitea.kood.tech/ivanandreev/pathfinder/internal/app"
	"gitea.kood.tech/ivanandreev/pathfinder/internal/config"
)

var (
	helpFlag bool
	usageStr = "\x1b[4mUsage:\x1b[24m go run . [path to file containing network map] [start station] [end station] [number of trains]"
)

func init() {
	flag.BoolVar(&helpFlag, "help", false, "Display this help message or help for a specific command")
	flag.BoolVar(&helpFlag, "h", false, "Display this help message (shorthand)")
}

func main() {
	flag.Parse()

	// Parse and validate
	networkMapPath, startStation, endStation, numTrains, err := parseArgs()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		fmt.Fprintln(os.Stderr, usageStr)
		os.Exit(1)
	}

	cfg := config.MustLoad(networkMapPath, startStation, endStation, numTrains)

	app := app.New(cfg)

	// run the app, exit on error
	if err := app.Run(); err != nil {
		os.Exit(1)
	}
}

func parseArgs() (string, string, string, int, error) {
	if helpFlag {
		fmt.Println(usageStr)
		fmt.Println("\nFlags:")
		flag.PrintDefaults()
		os.Exit(0)
	}

	if flag.NArg() != 4 {
		return "", "", "", 0, fmt.Errorf("incorrect number of arguments. Expected 4 arguments, got %d", flag.NArg())
	}

	args := flag.Args()
	numTrains, err := strconv.Atoi(args[3])
	if err != nil || numTrains <= 0 {
		return "", "", "", 0, fmt.Errorf("number of trains must be a positive integer, got '%s'", args[3])
	}

	return args[0], args[1], args[2], numTrains, nil
}
