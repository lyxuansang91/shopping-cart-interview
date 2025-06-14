package core

import (
	"context"
	"fmt"
	"time"

	"go.temporal.io/sdk/client"
)

func NewTemporalClient(ctx context.Context, host string) (client.Client, error) {
	// Add retry logic for Temporal client connection
	var temporalClient client.Client
	var err error
	for i := 0; i < 60; i++ { // Increase to 60 seconds
		temporalClient, err = client.Dial(client.Options{
			HostPort:  host,      // Use the docker service name
			Namespace: "default", // Add explicit namespace
		})
		if err == nil {
			break
		}
		time.Sleep(time.Second)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to create temporal client after 60 attempts: %w", err)
	}
	return temporalClient, nil
}
