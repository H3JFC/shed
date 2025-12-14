package store

import (
	"errors"
	"testing"
)

const (
	apiKey = "api_key"
)

func TestAddSecret_OK(t *testing.T) {
	t.Parallel()
	s := prepNewStore(t)

	secret, err := s.AddSecret(apiKey, "secret-value-123", "API key for external service")
	if err != nil {
		t.Fatalf("unexpected error adding secret: %v", err)
	}

	if secret.Key != apiKey {
		t.Fatalf("expected secret key %v, got %v", apiKey, secret.Key)
	}

	if secret.Value != "secret-value-123" {
		t.Fatalf("expected secret value %v, got %v", "secret-value-123", secret.Value)
	}

	if secret.Description != "API key for external service" {
		t.Fatalf("expected secret description %v, got %v", "API key for external service", secret.Description)
	}
}

func TestAddSecret_OKMaxLength(t *testing.T) {
	t.Parallel()
	s := prepNewStore(t)

	secret, err := s.AddSecret(validName32Char, "secret-value", "secret with max length key")
	if err != nil {
		t.Fatalf("unexpected error adding secret with max length key: %v", err)
	}

	if secret.Key != validName32Char {
		t.Fatalf("expected secret key %v, got %v", validName32Char, secret.Key)
	}
}

func TestAddSecret_ErrAlreadyExists(t *testing.T) {
	t.Parallel()
	s := prepNewStore(t)

	_, err := s.AddSecret(apiKey, "secret-value-123", "API key")
	if err != nil {
		t.Fatalf("unexpected error adding secret: %v", err)
	}

	_, err = s.AddSecret(apiKey, "different-value", "another secret with same key")
	if err == nil {
		t.Fatalf("expected error %v, got nil", ErrAlreadyExists)
	}

	if !errors.Is(err, ErrAlreadyExists) {
		t.Fatalf("expected error %v, got %v", ErrAlreadyExists, err)
	}
}

func TestAddSecret_Err(t *testing.T) { // nolint: funlen
	t.Parallel()

	type testcase struct {
		key   string
		value string
		want  error
	}

	tests := map[string]testcase{
		"invalid-starts-with-number": {
			key:   "123secret",
			value: "value",
			want:  ErrInvalidCommandName,
		},
		"invalid-starts-with-underscore": {
			key:   "_secret",
			value: "value",
			want:  ErrInvalidCommandName,
		},
		"invalid-contains-space": {
			key:   "api key",
			value: "value",
			want:  ErrInvalidCommandName,
		},
		"invalid-contains-hyphen": {
			key:   "api-key",
			value: "value",
			want:  ErrInvalidCommandName,
		},
		"invalid-contains-dot": {
			key:   "api.key",
			value: "value",
			want:  ErrInvalidCommandName,
		},
		"invalid-contains-special-chars": {
			key:   "api@key",
			value: "value",
			want:  ErrInvalidCommandName,
		},
		"invalid-empty": {
			key:   "",
			value: "value",
			want:  ErrInvalidCommandName,
		},
		"invalid-too-long-33-chars": {
			key:   invalidName33CharTooLong,
			value: "value",
			want:  ErrInvalidCommandName,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			s := prepNewStore(t)

			_, err := s.AddSecret(tc.key, tc.value, "secret with errors")
			if err == nil {
				t.Fatalf("expected error %v, got nil", tc.want)
			}

			if !errors.Is(err, tc.want) {
				t.Fatalf("expected error %v, got %v", tc.want, err)
			}
		})
	}
}

func TestListSecrets_OK(t *testing.T) {
	t.Parallel()
	s := prepNewStore(t)

	// Create multiple secrets
	_, err := s.AddSecret(apiKey, "secret-value-1", "API key")
	if err != nil {
		t.Fatalf("unexpected error adding secret: %v", err)
	}

	_, err = s.AddSecret("db_password", "secret-value-2", "Database password")
	if err != nil {
		t.Fatalf("unexpected error adding secret: %v", err)
	}

	_, err = s.AddSecret("jwt_secret", "secret-value-3", "JWT signing secret")
	if err != nil {
		t.Fatalf("unexpected error adding secret: %v", err)
	}

	// List all secrets
	secrets, err := s.ListSecrets()
	if err != nil {
		t.Fatalf("unexpected error listing secrets: %v", err)
	}

	if len(secrets) != 3 {
		t.Fatalf("expected 3 secrets, got %v", len(secrets))
	}

	// Verify secret keys exist
	keys := make(map[string]bool)
	for _, secret := range secrets {
		keys[secret.Key] = true
	}

	expectedKeys := []string{apiKey, "db_password", "jwt_secret"}
	for _, key := range expectedKeys {
		if !keys[key] {
			t.Fatalf("expected secret %v to be in list", key)
		}
	}
}

func TestListSecrets_Empty(t *testing.T) {
	t.Parallel()
	s := prepNewStore(t)

	// List secrets when none exist
	secrets, err := s.ListSecrets()
	if err != nil {
		t.Fatalf("unexpected error listing secrets: %v", err)
	}

	if len(secrets) != 0 {
		t.Fatalf("expected 0 secrets, got %v", len(secrets))
	}
}

func TestUpdateSecret_OK(t *testing.T) {
	t.Parallel()
	s := prepNewStore(t)

	// Create initial secret
	_, err := s.AddSecret(apiKey, "old-value", "API key")
	if err != nil {
		t.Fatalf("unexpected error adding secret: %v", err)
	}

	// Update secret
	updatedSecret, err := s.UpdateSecret(apiKey, "new-value", "Updated API key")
	if err != nil {
		t.Fatalf("unexpected error updating secret: %v", err)
	}

	if updatedSecret.Key != apiKey {
		t.Fatalf("expected secret key %v, got %v", apiKey, updatedSecret.Key)
	}

	if updatedSecret.Value != "new-value" {
		t.Fatalf("expected secret value %v, got %v", "new-value", updatedSecret.Value)
	}

	if updatedSecret.Description != "Updated API key" {
		t.Fatalf("expected secret description %v, got %v", "Updated API key", updatedSecret.Description)
	}
}

func TestUpdateSecret_OKMaxLength(t *testing.T) {
	t.Parallel()
	s := prepNewStore(t)

	// Create initial secret with max length key
	_, err := s.AddSecret(validName32Char, "old-value", "secret with max length key")
	if err != nil {
		t.Fatalf("unexpected error adding secret: %v", err)
	}

	// Update secret with max length key
	updatedSecret, err := s.UpdateSecret(validName32Char, "new-value", "updated secret with max length key")
	if err != nil {
		t.Fatalf("unexpected error updating secret with max length key: %v", err)
	}

	if updatedSecret.Key != validName32Char {
		t.Fatalf("expected secret key %v, got %v", validName32Char, updatedSecret.Key)
	}
}

func TestUpdateSecret_ErrSecretNotFound(t *testing.T) {
	t.Parallel()
	s := prepNewStore(t)

	_, err := s.UpdateSecret("nonexistent", "value", "description")
	if err == nil {
		t.Fatalf("expected error %v, got nil", ErrSecretNotFound)
	}

	if !errors.Is(err, ErrSecretNotFound) {
		t.Fatalf("expected error %v, got %v", ErrSecretNotFound, err)
	}
}

func TestUpdateSecret_Err(t *testing.T) {
	t.Parallel()

	type testcase struct {
		key  string
		want error
	}

	tests := map[string]testcase{
		"invalid-starts-with-number": {
			key:  "123secret",
			want: ErrInvalidCommandName,
		},
		"invalid-starts-with-underscore": {
			key:  "_secret",
			want: ErrInvalidCommandName,
		},
		"invalid-contains-space": {
			key:  "api key",
			want: ErrInvalidCommandName,
		},
		"invalid-contains-hyphen": {
			key:  "api-key",
			want: ErrInvalidCommandName,
		},
		"invalid-contains-dot": {
			key:  "api.key",
			want: ErrInvalidCommandName,
		},
		"invalid-contains-special-chars": {
			key:  "api@key",
			want: ErrInvalidCommandName,
		},
		"invalid-empty": {
			key:  "",
			want: ErrInvalidCommandName,
		},
		"invalid-too-long-33-chars": {
			key:  invalidName33CharTooLong,
			want: ErrInvalidCommandName,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			s := prepNewStore(t)

			// No need to create a secret first since validation happens before lookup
			_, err := s.UpdateSecret(tc.key, "value", "description")
			if err == nil {
				t.Fatalf("expected error %v, got nil", tc.want)
			}

			if !errors.Is(err, tc.want) {
				t.Fatalf("expected error %v, got %v", tc.want, err)
			}
		})
	}
}

func TestRemoveSecret_OK(t *testing.T) {
	t.Parallel()
	s := prepNewStore(t)

	if _, err := s.AddSecret(apiKey, "secret-value", "API key"); err != nil {
		t.Fatalf("unexpected error adding secret: %v", err)
	}

	err := s.RemoveSecret(apiKey)
	if err != nil {
		t.Fatalf("unexpected error removing secret: %v", err)
	}

	// Verify secret is removed
	_, err = s.GetSecretByKey(apiKey)
	if err == nil {
		t.Fatalf("expected error getting removed secret, got nil")
	}
}

func TestRemoveSecret_ErrSecretNotFound(t *testing.T) {
	t.Parallel()
	s := prepNewStore(t)

	err := s.RemoveSecret("does_not_exist")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestGetSecretByKey_OK(t *testing.T) {
	t.Parallel()
	s := prepNewStore(t)

	// Create a secret
	addedSecret, err := s.AddSecret(apiKey, "secret-value-123", "API key for service")
	if err != nil {
		t.Fatalf("unexpected error adding secret: %v", err)
	}

	// Retrieve secret by key
	secret, err := s.GetSecretByKey(apiKey)
	if err != nil {
		t.Fatalf("unexpected error getting secret: %v", err)
	}

	if secret.ID != addedSecret.ID {
		t.Fatalf("expected secret ID %v, got %v", addedSecret.ID, secret.ID)
	}

	if secret.Key != apiKey {
		t.Fatalf("expected secret key %v, got %v", apiKey, secret.Key)
	}

	if secret.Value != "secret-value-123" {
		t.Fatalf("expected secret value %v, got %v", "secret-value-123", secret.Value)
	}

	if secret.Description != "API key for service" {
		t.Fatalf("expected secret description %v, got %v", "API key for service", secret.Description)
	}

	if secret.CreatedAt == "" {
		t.Fatalf("expected CreatedAt to be set")
	}

	if secret.UpdatedAt == "" {
		t.Fatalf("expected UpdatedAt to be set")
	}
}

func TestGetSecretByKey_ErrSecretNotFound(t *testing.T) {
	t.Parallel()
	s := prepNewStore(t)

	_, err := s.GetSecretByKey("does_not_exist")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestSecretOperations_MultipleUpdates(t *testing.T) {
	t.Parallel()
	s := prepNewStore(t)

	// Create secret
	_, err := s.AddSecret(apiKey, "value1", "description1")
	if err != nil {
		t.Fatalf("unexpected error adding secret: %v", err)
	}

	// Update multiple times
	_, err = s.UpdateSecret(apiKey, "value2", "description2")
	if err != nil {
		t.Fatalf("unexpected error on first update: %v", err)
	}

	_, err = s.UpdateSecret(apiKey, "value3", "description3")
	if err != nil {
		t.Fatalf("unexpected error on second update: %v", err)
	}

	secret, err := s.UpdateSecret(apiKey, "value4", "description4")
	if err != nil {
		t.Fatalf("unexpected error on third update: %v", err)
	}

	if secret.Value != "value4" {
		t.Fatalf("expected final value %v, got %v", "value4", secret.Value)
	}

	if secret.Description != "description4" {
		t.Fatalf("expected final description %v, got %v", "description4", secret.Description)
	}
}

func TestSecretOperations_EmptyValues(t *testing.T) {
	t.Parallel()
	s := prepNewStore(t)

	// Test with empty value (should be allowed)
	secret, err := s.AddSecret(apiKey, "", "API key with empty value")
	if err != nil {
		t.Fatalf("unexpected error adding secret with empty value: %v", err)
	}

	if secret.Value != "" {
		t.Fatalf("expected empty value, got %v", secret.Value)
	}

	// Test with empty description (should be allowed)
	secret2, err := s.AddSecret("db_password", "secret123", "")
	if err != nil {
		t.Fatalf("unexpected error adding secret with empty description: %v", err)
	}

	if secret2.Description != "" {
		t.Fatalf("expected empty description, got %v", secret2.Description)
	}
}
