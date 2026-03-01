package sur

// SUR (Shortest Unconstrained Route)

// Responsibility: Business rules and logical operations.

// Validations:
// - It displays "Error" on stderr when the start station does not exist.
// - It displays "Error" on stderr when the end station does not exist.
// Is there an actual path between them?

// Why here:
// These rules define how your application behaves regardless of where the data comes from.
// If you later add an HTTP API or a gRPC transport, your core business logic still enforces that a trip must have distinct start and end points.
