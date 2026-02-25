package local

import (
	"bufio"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"

	"gitea.kood.tech/ivanandreev/pathfinder/internal/domain"
	"gitea.kood.tech/ivanandreev/pathfinder/internal/lib/e"
)

// Responsibility: Data integrity and structural validation of the map file.

// Validations:
// - Malformed lines,
// - negative coordinates,
// - two stations sharing the exact same coordinates,
// - a connection pointing to a station that doesn't exist.

// The repository's job is to read the file and return a valid adjeisency list.
// If the file is broken, the repository should fail to build the object and return an error
// (e.g., fmt.Errorf("invalid map: duplicate coordinates at line 10")).
// The service layer should never have to worry if the map structurally makes sense.

// Errors to check:
// - It displays "Error" on stderr when the map does not contain a "stations:" section.
// - It displays "Error" on stderr when the map does not contain a "connections:" section.
// - It displays "Error" on stderr when the start station does not exist.
// - It displays "Error" on stderr when the end station does not exist.
// - It displays "Error" on stderr when no path exists between the start and end stations.
// - It displays "Error" on stderr when duplicate routes exist between two stations, including in reverse.
// - It displays "Error" on stderr when any of the coordinates are not valid positive integers.
// - It displays "Error" on stderr when two stations exist at the same coordinates.
// - It displays "Error" on stderr when a connection is made with a station which does not exist.
// - It displays "Error" on stderr when station names are duplicated.
// - It displays "Error" on stderr when station names are invalid.
// - It displays "Error" on stderr when a map contains more than 10000 stations.

// - Disconnected components (SCCs) ???
// - We don't have it in the errors list, But: if after the connection buildup we have some stations in the list that doesn't have route to it. should we also report an error???

type Storage struct {
	log          *slog.Logger
	MapPath      string
	StartStation string
	EndStation   string
}

func New(logger *slog.Logger, filePath, start, end string) *Storage {
	return &Storage{
		log:          logger,
		MapPath:      filePath,
		StartStation: start,
		EndStation:   end,
	}
}

/*
# London Network Map

stations:
waterloo,3,1
victoria,6,7
euston,11,23
st_pancras,5,15

connections:
waterloo-victoria
waterloo-euston
st_pancras-euston
victoria-st_pancras
*/

func (s *Storage) ParseMap() (domain.Graph, error) {
	const op = "storage.local.ParseMap"

	log := s.log.With(
		slog.String("op", op),
	)

	file, err := os.Open(s.MapPath)
	if err != nil {
		log.Error("can't open map file", slog.Any("error", err))
		return domain.Graph{}, e.Wrap("can't open map file", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCount := 0

	stationsSectionExists := false
	connectionsSectionExists := false
	const stationsSection = "stations:"
	const connectionsSection = "connections:"

	var graph domain.Graph
	var stations []domain.Station
	var edges []domain.Edge
	nameToID := make(map[string]int)
	var edgesMap [][]domain.Edge

	for scanner.Scan() {
		lineCount++
		line := scanner.Text()

		// Skip empty lines and lines starting with comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.Contains(line, "#") {
			splittedLine := strings.SplitN(line, "#", 2)
			line = splittedLine[0]
		}

		line = strings.TrimSpace(line)

		if strings.EqualFold(line, stationsSection) {
			stationsSectionExists = true
		}

		if strings.Contains(line, ",") && stationsSectionExists && !connectionsSectionExists {
			stationParams := strings.Split(line, ",")

			stationName := stationParams[0]
			x, err := strconv.ParseFloat(stationParams[1], 64)
			if err != nil {
				log.Error("can't parse float for Y coordinate",
					slog.String("station", stationName),
					slog.Int("map line", lineCount),
					slog.Any("error", err),
				)
				return domain.Graph{}, fmt.Errorf("can't parse float for y coordinate, station name: %s, map line: %d, error: %v", stationName, lineCount, err)
			}

			y, err := strconv.ParseFloat(stationParams[2], 64)
			if err != nil {
				log.Error("can't parse float for y coordinate",
					slog.String("station", stationName),
					slog.Int("map line", lineCount),
					slog.Any("error", err),
				)
				return domain.Graph{}, fmt.Errorf("can't parse float for y coordinate, station name: %s, map line: %d, error: %v", stationName, lineCount, err)
			}

			tmpStation := domain.Station{
				Name: stationName,
				X:    x,
				Y:    y,
			}

			stations = append(stations, tmpStation)
		}

		if strings.EqualFold(line, connectionsSection) && !stationsSectionExists {
			err := errors.New("the connections section is found before the stations section")
			log.Error("the map does not contain a stations section", slog.Any("error", err))
			return domain.Graph{}, e.Wrap("the map does not contain a stations section", err)
		} else {
			connectionsSectionExists = true
		}

	}

	if !stationsSectionExists || !connectionsSectionExists {
		err := errors.New("the stations or connections section is missing")
		log.Error("the map does not contain stations or connections section",
			slog.Bool("stations", stationsSectionExists),
			slog.Bool("connections", connectionsSectionExists),
			slog.Any("error", err),
		)
		return domain.Graph{}, err
	}

	return graph, nil
}
