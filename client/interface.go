package client

import (
	healthcheckv1 "go.admiral.io/admiral-go/proto/healthcheck/v1"
	userv1 "go.admiral.io/admiral-go/proto/user/v1"
)

// AdmiralClient provides access to Admiral service clients.
type AdmiralClient interface {
	// Healthcheck returns the HealthcheckAPI client.
	Healthcheck() healthcheckv1.HealthcheckAPIClient
	// User returns the UserAPI client.
	User() userv1.UserAPIClient

	// ValidateToken validates the client's auth token.
	ValidateToken() error

	// GetTokenInfo returns information about the client's auth token.
	GetTokenInfo() (*TokenInfo, error)

	// Version returns the client library version string.
	Version() string

	// Close closes the underlying connection.
	Close() error
}
