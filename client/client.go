package client

import (
	"context"
	"crypto/tls"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	healthcheckv1 "go.admiral.io/sdk/proto/healthcheck/v1"
	userv1 "go.admiral.io/sdk/proto/user/v1"
)

// Compile-time check that Client implements AdmiralClient
var _ AdmiralClient = (*Client)(nil)

// Client is the Admiral API client.
type Client struct {
	conn      *grpc.ClientConn
	logger    Logger
	authToken string
	healthcheck healthcheckv1.HealthcheckAPIClient
	user userv1.UserAPIClient
}

// New creates a new Admiral client with the given configuration.
func New(ctx context.Context, cfg Config) (*Client, error) {
	if err := cfg.CheckAndSetDefaults(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	dialOpts := cfg.ConnectionOptions.DialOptions

	// Configure transport credentials
	if cfg.ConnectionOptions.Insecure {
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		tlsConfig := cfg.ConnectionOptions.TLSConfig
		if tlsConfig == nil {
			tlsConfig = &tls.Config{MinVersion: tls.VersionTLS12}
		}
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))
	}

	// Add user agent
	dialOpts = append(dialOpts, grpc.WithUserAgent(ClientUserAgent()))

	// Dial with timeout
	dialCtx, cancel := context.WithTimeout(ctx, cfg.ConnectionOptions.DialTimeout)
	defer cancel()

	conn, err := grpc.DialContext(dialCtx, cfg.HostPort, dialOpts...)
	if err != nil {
		return nil, fmt.Errorf("failed to dial %s: %w", cfg.HostPort, err)
	}

	cfg.Logger.Info("connected to Admiral API", "host", cfg.HostPort)

	return &Client{
		conn:      conn,
		logger:    cfg.Logger,
		authToken: cfg.AuthToken,
		healthcheck: healthcheckv1.NewHealthcheckAPIClient(conn),
		user: userv1.NewUserAPIClient(conn),
	}, nil
}

// Healthcheck returns the HealthcheckAPI client.
func (c *Client) Healthcheck() healthcheckv1.HealthcheckAPIClient {
	return c.healthcheck
}

// User returns the UserAPI client.
func (c *Client) User() userv1.UserAPIClient {
	return c.user
}

// ValidateToken validates the client's auth token format and expiration.
func (c *Client) ValidateToken() error {
	return ValidateAuthToken(c.authToken)
}

// GetTokenInfo returns information about the client's auth token.
func (c *Client) GetTokenInfo() (*TokenInfo, error) {
	claims, err := ParseJWTToken(c.authToken)
	if err != nil {
		return nil, err
	}
	return &TokenInfo{JWTClaims: claims}, nil
}

// Version returns the client library version string.
func (c *Client) Version() string {
	return Version()
}

// Close closes the underlying gRPC connection.
func (c *Client) Close() error {
	if c.conn != nil {
		c.logger.Debug("closing connection")
		return c.conn.Close()
	}
	return nil
}
