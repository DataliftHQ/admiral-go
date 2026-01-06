// Package client provides Admiral API clients with connection management,
// authentication, and token validation.
//
// # Quick Start
//
// Create a client and access services:
//
//	c, err := client.New(ctx, client.Config{
//	    AuthToken: "your-token",
//	})
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer c.Close()
//
//	// Access services via accessors
//	resp, err := c.User().GetUser(ctx, req)
//
// # Configuration
//
// The Config struct provides options for customizing the client:
//
//   - AuthToken: Required authentication token
//   - ConnectionOptions: TLS, timeouts, keepalive settings
//   - Logger: Custom logger implementation
//
// # Token Validation
//
// The client validates JWT tokens on creation and provides methods for
// ongoing token validation:
//
//	if err := c.ValidateToken(); err != nil {
//	    // Token is invalid or expired
//	}
//
//	info, _ := c.GetTokenInfo()
//	fmt.Println("Expires in:", info.ExpiresIn())
package client
