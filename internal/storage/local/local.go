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
