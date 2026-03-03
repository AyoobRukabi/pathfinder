package local

import (
	"bufio"
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
// The repository's job is to read the file and return a valid adjeisency list.
// If the file is broken, the repository should fail to build the object and return an error.
// (e.g., fmt.Errorf("invalid map: duplicate coordinates at line 10")).
// The service layer should never have to worry if the map structurally makes sense.

// Errors to check:
// - [x] It displays "Error" on stderr when the map does not contain a "stations:" section.
// - [x] It displays "Error" on stderr when the map does not contain a "connections:" section.
// - [x] It displays "Error" on stderr when any of the coordinates are not valid positive integers.
// - [x] It displays "Error" on stderr when two stations exist at the same coordinates.
// - [x] It displays "Error" on stderr when station names are duplicated.
// - [x] It displays "Error" on stderr when station names are invalid.
// - [x] It displays "Error" on stderr when a map contains more than 10000 stations.

// - [x] It displays "Error" on stderr when a connection is made with a station which does not exist.
// - [x] It displays "Error" on stderr when duplicate routes exist between two stations, including in reverse.

const (
	stationsSection    = "stations:"
	connectionsSection = "connections:"
)

type parserState int

const (
	stateInit parserState = iota
	stateStations
	stateConnections
)

type coord struct{ x, y int }
type edge struct{ u, v int }

type Storage struct {
	log     *slog.Logger
	MapPath string
}

func New(logger *slog.Logger, filePath string) *Storage {
	return &Storage{
		log:     logger,
		MapPath: filePath,
	}
}

func (s *Storage) BuildMap() (domain.MapData, error) {
	const op = "storage.local.BuildMap"

	log := s.log.With(
		slog.String("op", op),
	)

	file, err := os.Open(s.MapPath)
	if err != nil {
		log.Error("can't open map file", slog.Any("error", err))
		return domain.MapData{}, e.Wrap("can't open map file", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCount := 0

	var adjList [][]int
	var stations []domain.Station          // statsions[0] -> "waterloo"
	nodeNameToID := make(map[string]int)   // map of IDs "waterloo" -> [0]
	seenCoords := make(map[coord]string)   // to check stations with the same coordinates
	seenConnections := make(map[edge]bool) // to check duplicate edges

	currentState := stateInit

	for scanner.Scan() {
		lineCount++
		line := scanner.Text()

		// Skip empty lines and lines starting with comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Remove comments
		if strings.Contains(line, "#") {
			cutLine, _, _ := strings.Cut(line, "#")
			line = cutLine
		}

		line = strings.TrimSpace(line)

		if strings.EqualFold(line, stationsSection) {
			currentState = stateStations
			continue
		}

		if strings.EqualFold(line, connectionsSection) {
			if currentState != stateStations {
				log.Error("connections section found before stations section", slog.Int("map line:", lineCount))
				return domain.MapData{}, fmt.Errorf("connections section found before stations section, line: %d", lineCount)
			}
			currentState = stateConnections
			continue
		}

		if currentState == stateStations {

			if len(stations) >= 10000 {
				log.Error("map contains more than 10000 stations", slog.Int("map line", lineCount))
				return domain.MapData{}, fmt.Errorf("map contains more than 10000 stations, map line: %d", lineCount)
			}
			// We are strictly parsing stations here.
			station, err := buildStation(log, line, lineCount)
			if err != nil {
				log.Error("invalid station line", slog.Any("error", err))
				return domain.MapData{}, err
			}

			if _, ok := nodeNameToID[station.Name]; ok {
				log.Error("station names are duplicated",
					slog.String("station", station.Name),
					slog.Int("map line", lineCount))
				return domain.MapData{}, fmt.Errorf("station names are duplicated, station name: %s, map line: %d",
					station.Name, lineCount)
			}

			c := coord{x: station.X, y: station.Y}
			if existingName, exists := seenCoords[c]; exists {
				log.Error("two stations exist at the same coordinates",
					slog.String("station one", station.Name),
					slog.String("station two", existingName),
					slog.Int("map line", lineCount),
				)
				return domain.MapData{}, fmt.Errorf("two stations exist at the same coordinates, stations %s and %s, map line: %d",
					existingName, station.Name, lineCount)
			}
			seenCoords[c] = station.Name

			stations = append(stations, station)
			nodeNameToID[station.Name] = len(stations) - 1
			adjList = append(adjList, []int{}) // add an empty row so adjList[id] exists
		} else if currentState == stateConnections {
			// Validate Connections here
			nodeIDs, err := validateConnection(log, line, lineCount, nodeNameToID)
			if err != nil {
				log.Error("invalid connection line", slog.Any("error", err))
				return domain.MapData{}, err
			}

			from := nodeIDs[0]
			to := nodeIDs[1]

			e := makeEdge(from, to)
			if seenConnections[e] {
				log.Error("duplicate routes exist between two stations",
					slog.String("station", stations[from].Name),
					slog.String("station", stations[to].Name),
					slog.Int("map line", lineCount))
				return domain.MapData{}, fmt.Errorf("duplicate routes exist between two stations, %s and %s, map line: %d",
					stations[from].Name, stations[to].Name, lineCount)
			}
			seenConnections[e] = true

			adjList[from] = append(adjList[from], to)
			adjList[to] = append(adjList[to], from)
		} else {
			log.Error("invalid data found before any section headers", slog.Int("map line:", lineCount))
			return domain.MapData{}, fmt.Errorf("invalid data found before any section headers, map line: %d", lineCount)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Error("unexpected scanner error", slog.Int("map line:", lineCount))
		return domain.MapData{}, fmt.Errorf("unexpected scanner error, line: %d, error: %w", lineCount, err)
	}

	if currentState != stateConnections {
		return domain.MapData{}, fmt.Errorf("the map is missing the connections section")
	}

	mapData := domain.MapData{
		Stations:         stations,
		StationsNameToID: nodeNameToID,
		AdjList:          adjList,
	}

	return mapData, nil
}

func buildStation(log *slog.Logger, line string, lineCount int) (domain.Station, error) {
	const op = "storage.local.buildStation"

	log = log.With(
		slog.String("op", op),
	)

	var stationName string
	stationParams := strings.Split(line, ",")

	if len(stationParams) != 3 {
		log.Error("station data is malformed",
			slog.Int("map line", lineCount),
		)
		return domain.Station{}, fmt.Errorf("station data is malformed, map line: %d", lineCount)
	}

	paramOne := strings.TrimSpace(stationParams[0])

	if paramOne == "" {
		log.Error("station name is invalid",
			slog.Int("map line", lineCount),
		)
		return domain.Station{}, fmt.Errorf("station name is invalid, map line: %d", lineCount)
	} else {
		stationName = paramOne
	}

	x, err := strconv.Atoi(strings.TrimSpace(stationParams[1]))
	if err != nil {
		log.Error("can't parse int for x coordinate",
			slog.String("station", stationName),
			slog.Int("map line", lineCount),
			slog.Any("error", err),
		)
		return domain.Station{}, fmt.Errorf("can't parse int for x coordinate, station name: %s, map line: %d, error: %w", stationName, lineCount, err)
	}

	if x < 0 {
		log.Error("x coordinate is negative",
			slog.String("station", stationName),
			slog.Int("map line", lineCount),
		)
		return domain.Station{}, fmt.Errorf("the \"x\" coordinate is not valid positive integer, station name: %s, map line: %d", stationName, lineCount)
	}

	y, err := strconv.Atoi(strings.TrimSpace(stationParams[2]))
	if err != nil {
		log.Error("can't parse float for y coordinate",
			slog.String("station", stationName),
			slog.Int("map line", lineCount),
			slog.Any("error", err),
		)
		return domain.Station{}, fmt.Errorf("can't parse int for y coordinate, station name: %s, map line: %d, error: %w", stationName, lineCount, err)
	}

	if y < 0 {
		log.Error("y coordinate is negative",
			slog.String("station", stationName),
			slog.Int("map line", lineCount),
		)
		return domain.Station{}, fmt.Errorf("the \"y\" coordinate is not valid positive integer, station name: %s, map line: %d", stationName, lineCount)
	}

	tmpStation := domain.Station{
		Name: stationName,
		X:    x,
		Y:    y,
	}

	return tmpStation, nil
}

func validateConnection(log *slog.Logger, line string, lineCount int, nodeToIdMap map[string]int) ([]int, error) {
	const op = "storage.local.validateConnection"

	log = log.With(
		slog.String("op", op),
	)

	nodes := make([]int, 0, 2)

	rawNodes := strings.Split(line, "-")

	if len(rawNodes) != 2 {
		log.Error("connection data is not complete",
			slog.Int("map line", lineCount),
		)
		return nil, fmt.Errorf("connection data is not complete, map line: %d", lineCount)
	}

	nodeOne := strings.TrimSpace(rawNodes[0])
	nodeTwo := strings.TrimSpace(rawNodes[1])

	if strings.EqualFold(nodeOne, nodeTwo) {
		log.Error("loop - connection to the same station",
			slog.String("station", nodeOne),
			slog.Int("map line", lineCount),
		)
		return nil, fmt.Errorf("loop - connection to the same station, station name: %s, map line: %d", nodeOne, lineCount)
	}

	trimmedNodes := make([]string, 0, 2)
	trimmedNodes = append(trimmedNodes, nodeOne, nodeTwo)

	for i := range trimmedNodes {
		if id, ok := nodeToIdMap[trimmedNodes[i]]; !ok {
			log.Error("connection is made with a station which does not exist",
				slog.String("station", trimmedNodes[i]),
				slog.Int("map line", lineCount),
			)
			return nil, fmt.Errorf("connection is made with a station which does not exist, station name: %s, map line: %d", trimmedNodes[i], lineCount)
		} else {
			nodes = append(nodes, id)
		}
	}

	return nodes, nil
}

// Helper func to revert edge IDs
func makeEdge(u, v int) edge {
	if u < v {
		return edge{u, v}
	}
	return edge{v, u}
}
