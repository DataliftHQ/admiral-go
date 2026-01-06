package client

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

type Config struct {
	HostPort          string
	AuthToken         string
	ConnectionOptions ConnectionOptions
	Logger            Logger
}

type ConnectionOptions struct {
	TLSConfig                    *tls.Config
	Insecure                     bool
	DialOptions                  []grpc.DialOption
	DialTimeout                  time.Duration
	EnableKeepAliveCheck         bool
	KeepAliveTime                time.Duration
	KeepAliveTimeout             time.Duration
	KeepAlivePermitWithoutStream bool
}

func (c *Config) CheckAndSetDefaults() error {
	if c.Logger == nil {
		c.Logger = NewDefaultLogger()
	}

	if c.HostPort == "" {
		c.HostPort = DefaultHostPort
	}

	_, _, err := net.SplitHostPort(c.HostPort)
	if err != nil {
		return fmt.Errorf("invalid host:port format %q (both host and port required): %w", c.HostPort, err)
	}

	if !c.ConnectionOptions.Insecure && c.ConnectionOptions.TLSConfig == nil {
		c.ConnectionOptions.TLSConfig = &tls.Config{
			MinVersion: tls.VersionTLS12, // Modern TLS minimum
		}
	}
	if c.ConnectionOptions.Insecure && c.ConnectionOptions.TLSConfig != nil {
		c.Logger.Warn("TLSConfig is set but ignored because Insecure is true")
		c.ConnectionOptions.TLSConfig = nil
	}

	if c.ConnectionOptions.DialTimeout == 0 {
		c.ConnectionOptions.DialTimeout = DefaultDialTimeout
	}
	if c.ConnectionOptions.KeepAliveTime == 0 {
		c.ConnectionOptions.KeepAliveTime = DefaultKeepAliveTime
	}
	if c.ConnectionOptions.KeepAliveTimeout == 0 {
		c.ConnectionOptions.KeepAliveTimeout = DefaultKeepAliveTimeout
	}

	if len(c.AuthToken) == 0 {
		return errors.New("auth token is required")
	}

	// Validate token format and expiration
	if err := ValidateAuthToken(c.AuthToken); err != nil {
		return fmt.Errorf("auth token validation failed: %w", err)
	}
	c.ConnectionOptions.DialOptions = append(
		c.ConnectionOptions.DialOptions,
		grpc.WithPerRPCCredentials(tokenAuth{
			token:               c.AuthToken,
			requireTransportSec: !c.ConnectionOptions.Insecure,
		}),
	)

	if c.ConnectionOptions.EnableKeepAliveCheck {
		kap := keepalive.ClientParameters{
			Time:                c.ConnectionOptions.KeepAliveTime,
			Timeout:             c.ConnectionOptions.KeepAliveTimeout,
			PermitWithoutStream: c.ConnectionOptions.KeepAlivePermitWithoutStream,
		}
		c.ConnectionOptions.DialOptions = append(c.ConnectionOptions.DialOptions, grpc.WithKeepaliveParams(kap))
	}

	return nil
}

type tokenAuth struct {
	token               string
	requireTransportSec bool
}

func (t tokenAuth) GetRequestMetadata(ctx context.Context, in ...string) (map[string]string, error) {
	return map[string]string{
		"Authorization": "Bearer " + t.token,
	}, nil
}

func (t tokenAuth) RequireTransportSecurity() bool {
	return t.requireTransportSec
}
