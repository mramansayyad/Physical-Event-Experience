package adapter

import (
	"context"
	"fmt"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
)

// SecretManager handles secure retrieval of infrastructure credentials.
type SecretManager struct {
	client *secretmanager.Client
}

// NewSecretManager initializes the GCP Secret Manager client.
func NewSecretManager(ctx context.Context) (*SecretManager, error) {
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to configure Secret Manager client: %w", err)
	}
	return &SecretManager{client: client}, nil
}

// Close cleans up the underlying client connection.
func (sm *SecretManager) Close() error {
	return sm.client.Close()
}

// GetSecret fetches the latest version of a secret from GCP.
// It fails completely ("Always-Fail") if the secret cannot be found.
func (sm *SecretManager) GetSecret(ctx context.Context, projectID, secretName string) (string, error) {
	name := fmt.Sprintf("projects/%s/secrets/%s/versions/latest", projectID, secretName)
	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: name,
	}

	result, err := sm.client.AccessSecretVersion(ctx, req)
	if err != nil {
		return "", fmt.Errorf("CRITICAL: Failed to access secret '%s' - %w", secretName, err)
	}

	return string(result.Payload.Data), nil
}
