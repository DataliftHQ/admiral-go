package client

import "time"

// DefaultHostPort is the default API endpoint.
const DefaultHostPort = "api.admiral.io:443"

// DefaultDialTimeout is the default timeout for establishing a gRPC connection.
const DefaultDialTimeout = 30 * time.Second

// DefaultKeepAliveTime is the default interval for sending keepalive pings.
const DefaultKeepAliveTime = 30 * time.Second

// DefaultKeepAliveTimeout is the default timeout for keepalive ping responses.
const DefaultKeepAliveTimeout = 90 * time.Second
