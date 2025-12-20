package store

import (
	"context"
	"errors"
	"fmt"

	"h3jfc/shed/db"
)

var ErrSecretNotFound = errors.New("secret not found")

type Secret = db.Secret

func (s *Store) AddSecret(key, value, description string) (*Secret, error) {
	if err := validateName(key); err != nil {
		return nil, err
	}

	_, err := s.GetSecretByKey(key)
	if err == nil {
		return nil, fmt.Errorf("secret with key %q already exists: %w", key, ErrAlreadyExists)
	}

	secret, err := s.queries.CreateSecret(context.Background(), db.CreateSecretParams{
		Key:         key,
		Value:       value,
		Description: description,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create secret: %w", err)
	}

	return &secret, nil
}

func (s *Store) ListSecrets() ([]Secret, error) {
	secrets, err := s.queries.ListSecrets(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to list secrets: %w", err)
	}

	return secrets, nil
}

func (s *Store) UpdateSecret(key, value, description string) (*Secret, error) {
	if err := validateName(key); err != nil {
		return nil, err
	}

	prev, err := s.GetSecretByKey(key)
	if err != nil {
		return nil, fmt.Errorf("failed to update secret: %w", ErrSecretNotFound)
	}

	secret, err := s.queries.UpdateSecret(context.Background(), db.UpdateSecretParams{
		ID:          prev.ID,
		Value:       value,
		Description: description,
		Key:         key,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update secret: %w", err)
	}

	return &secret, nil
}

func (s *Store) RemoveSecret(key string) error {
	secret, err := s.GetSecretByKey(key)
	if err != nil {
		return fmt.Errorf("failed to get secret by key: %w", err)
	}

	err = s.queries.DeleteSecretByID(context.Background(), secret.ID)
	if err != nil {
		return fmt.Errorf("failed to delete secret: %w", err)
	}

	return nil
}

func (s *Store) GetSecretByKey(key string) (*Secret, error) {
	secret, err := s.queries.GetSecretByKey(context.Background(), key)
	if err != nil {
		return nil, fmt.Errorf("failed to get secret by key: %w", err)
	}

	return &secret, nil
}

func (s *Store) GetSecretsByKeys(keys []string) (*[]Secret, error) {
	secrets, err := s.queries.GetSecretsByKeys(context.Background(), keys)
	if err != nil {
		return nil, fmt.Errorf("failed to get secrets by keys: %w", err)
	}

	return &secrets, nil
}
