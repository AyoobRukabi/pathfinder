package local

// Responsibility: Data integrity and structural validation of the map file.

// Validations:
// - Malformed lines,
// - negative coordinates,
// - two stations sharing the exact same coordinates,
// - a connection pointing to a station that doesn't exist.

// The repository's job is to read the file and return a valid domain.Graph or domain.Map object.
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
