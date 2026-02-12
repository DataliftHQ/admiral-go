package client

import (
	agentv1 "go.admiral.io/sdk/proto/agent/v1"
	clusterv1 "go.admiral.io/sdk/proto/cluster/v1"
	healthcheckv1 "go.admiral.io/sdk/proto/healthcheck/v1"
	runnerv1 "go.admiral.io/sdk/proto/runner/v1"
	serviceaccountv1 "go.admiral.io/sdk/proto/serviceaccount/v1"
	userv1 "go.admiral.io/sdk/proto/user/v1"
)

// AdmiralClient provides access to Admiral service clients.
type AdmiralClient interface {
	// Agent returns the AgentAPI client.
	Agent() agentv1.AgentAPIClient
	// Cluster returns the ClusterAPI client.
	Cluster() clusterv1.ClusterAPIClient
	// Healthcheck returns the HealthcheckAPI client.
	Healthcheck() healthcheckv1.HealthcheckAPIClient
	// Runner returns the RunnerAPI client.
	Runner() runnerv1.RunnerAPIClient
	// ServiceAccount returns the ServiceAccountAPI client.
	ServiceAccount() serviceaccountv1.ServiceAccountAPIClient
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
