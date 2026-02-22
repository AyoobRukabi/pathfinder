package service

// Responsibility: Business rules and logical operations.

// Validations:
// - start == end (Error: Start and destination cannot be the same).
// Do start and end actually exist in the parsed map?
// Is there an actual path between them?

// Why here:
// These rules define how your application behaves regardless of where the data comes from.
// If you later add an HTTP API or a gRPC transport, your core business logic still enforces that a trip must have distinct start and end points.
