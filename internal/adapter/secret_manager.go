package adapter

import (
	"context"
	"fmt"
	"sync"
	"time"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
)

type cacheEntry struct {
	value  string
	expiry time.Time
}

// SecretManager handles secure retrieval of infrastructure credentials.
type SecretManager struct {
	client *secretmanager.Client
	cache  map[string]cacheEntry
	mu     sync.RWMutex
}

// NewSecretManager initializes the GCP Secret Manager client.
func NewSecretManager(ctx context.Context) (*SecretManager, error) {
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to configure Secret Manager client: %w", err)
	}
	return &SecretManager{
		client: client,
		cache:  make(map[string]cacheEntry),
	}, nil
}

// Close cleans up the underlying client connection.
func (sm *SecretManager) Close() error {
	return sm.client.Close()
}

// GetSecret fetches the latest version of a secret from GCP mapping transparent TTL optimization natively.
// It fails completely ("Always-Fail") if the secret cannot be found.
func (sm *SecretManager) GetSecret(ctx context.Context, projectID, secretName string) (string, error) {
	name := fmt.Sprintf("projects/%s/secrets/%s/versions/latest", projectID, secretName)

	sm.mu.RLock()
	entry, exists := sm.cache[name]
	sm.mu.RUnlock()

	// High-Scale burst optimization: Returning memory-mapped TTL directly
	if exists && time.Now().Before(entry.expiry) {
		return entry.value, nil
	}

	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: name,
	}

	result, err := sm.client.AccessSecretVersion(ctx, req)
	if err != nil {
		return "", fmt.Errorf("CRITICAL: Failed to access secret '%s' - %w", secretName, err)
	}

	payload := string(result.Payload.Data)

	sm.mu.Lock()
	sm.cache[name] = cacheEntry{
		value:  payload,
		expiry: time.Now().Add(5 * time.Minute), // Explicit 5M cache burst
	}
	sm.mu.Unlock()

	return payload, nil
}
