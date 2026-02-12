# admiral-go

Go client library for the Admiral API.

## Installation

```bash
go get go.admiral.io/sdk
```

## Usage

```go
package main

import (
	"context"
	"log"

	"go.admiral.io/sdk/client"
)

func main() {
	ctx := context.Background()

	// Create client
	c, err := client.New(ctx, client.Config{
		HostPort:  "api.admiral.io:443",
		AuthToken: "your-token-here",
	})
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	// Access services via accessors
	// c.Agent().MethodName(ctx, req)
	// c.Cluster().MethodName(ctx, req)
	// c.Healthcheck().MethodName(ctx, req)
	// c.Runner().MethodName(ctx, req)
	// c.ServiceAccount().MethodName(ctx, req)
	// c.User().MethodName(ctx, req)
}
```

## Available Services

| Service | Accessor | Import |
|---------|----------|--------|
| AgentAPI | `Agent()` | `go.admiral.io/sdk/proto/agent/v1` |
| ClusterAPI | `Cluster()` | `go.admiral.io/sdk/proto/cluster/v1` |
| HealthcheckAPI | `Healthcheck()` | `go.admiral.io/sdk/proto/healthcheck/v1` |
| RunnerAPI | `Runner()` | `go.admiral.io/sdk/proto/runner/v1` |
| ServiceAccountAPI | `ServiceAccount()` | `go.admiral.io/sdk/proto/serviceaccount/v1` |
| UserAPI | `User()` | `go.admiral.io/sdk/proto/user/v1` |

## Example

```go
package main

import (
	"context"
	"fmt"
	"log"

	"go.admiral.io/sdk/client"
	agentv1 "go.admiral.io/sdk/proto/agent/v1"
)

func main() {
	ctx := context.Background()

	c, err := client.New(ctx, client.Config{
		HostPort:  "api.admiral.io:443",
		AuthToken: os.Getenv("ADMIRAL_TOKEN"),
	})
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	// Validate token before making requests
	if err := c.ValidateToken(); err != nil {
		log.Fatal("Invalid token:", err)
	}

	// Call a service method
	resp, err := c.Agent().ListMethod(ctx, &agentv1.ListMethodRequest{})
	if err != nil {
		log.Fatal("Request failed:", err)
	}
	fmt.Printf("Response: %+v\n", resp)
}
```

## Configuration

```go
import (
	"crypto/tls"
	"time"

	"go.admiral.io/sdk/client"
)

cfg := client.Config{
	// Required: Server address
	HostPort: "api.admiral.io:443",

	// Required: Authentication token (JWT)
	AuthToken: "your-jwt-token",

	// Optional: Connection options
	ConnectionOptions: client.ConnectionOptions{
		// Use insecure connection (no TLS) - default: false
		Insecure: false,

		// Connection timeout - default: 10s
		DialTimeout: 10 * time.Second,

		// Custom TLS configuration
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},

		// Additional gRPC dial options
		DialOptions: []grpc.DialOption{},
	},

	// Optional: Custom logger (default: no-op logger)
	Logger: client.NewDefaultLogger(),
}

c, err := client.New(ctx, cfg)
```

## Token Validation

```go
// Validate token format and expiration
if err := c.ValidateToken(); err != nil {
	log.Fatal("Invalid token:", err)
}

// Get detailed token information
info, err := c.GetTokenInfo()
if err != nil {
	log.Fatal("Failed to parse token:", err)
}

fmt.Printf("Subject: %s\n", info.Subject)
fmt.Printf("Issuer: %s\n", info.Issuer)
fmt.Printf("Expires in: %v\n", info.ExpiresIn())
fmt.Printf("Is expired: %v\n", info.IsExpired())
```

## Version Information

```go
// Get library version
fmt.Println("Version:", client.Version())

// Get user agent string (useful for debugging)
fmt.Println("User-Agent:", client.ClientUserAgent())
```

## Requirements

- Go 1.25 or later

## License

MIT License - see [LICENSE](LICENSE.txt) for details.
