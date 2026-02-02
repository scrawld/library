package util

import (
	"context"
	"fmt"
	"net"
	"net/http"
)

// StartTestServer starts a test server with a dynamic port and returns the host URL,
// a cleanup function to shut down the server, and any error encountered.
func StartTestServer(handler http.Handler) (string, func(), error) {
	// Create a listener with a dynamic port
	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create listener: %w", err)
	}

	// Extract the dynamic port assigned by the OS
	addr := ln.Addr().String()
	host := fmt.Sprintf("http://%s", addr)

	// Start the server in a goroutine
	srv := &http.Server{
		Handler: handler,
	}
	go func() {
		if err := srv.Serve(ln); err != nil {
			fmt.Printf("Error starting server: %v\n", err)
		}
	}()

	// Return the host and a cleanup function to shut down the server
	cleanup := func() {
		srv.Shutdown(context.Background())
		ln.Close()
	}
	return host, cleanup, nil
}
