// Package store provides functions to manage a collection of items.
package store

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"slices"

	"github.com/spf13/viper"

	"h3jfc/shed/db"
	"h3jfc/shed/lib/brackets"
	"h3jfc/shed/lib/itertools"
	"h3jfc/shed/lib/sqlite3"
)

const (
	nameMaxLength = 32
	nameDetails   = "command names may only contain letters, numbers, hyphens, and underscores"
	nameLength    = "it must be between 1 and 32 characters long"
)

var (
	ErrNotFound           = errors.New("database could not be found or opened")
	ErrAlreadyExists      = errors.New("item already exists")
	ErrInvalidCommandName = errors.New("invalid command name")
	ErrParsingValueParams = errors.New("failed to parse value parameters")
	ErrCommandNotFound    = errors.New("command not found")
	ErrNameTooLong        = errors.New("command name is too long, it must be 40 characters or less")
)

type Store struct {
	queries *db.Queries
	dbtx    db.DBTX
}

func NewStoreFromConfig() (*Store, error) {
	dbPath := viper.GetString("shed-db.location")
	encryptionKey := viper.GetString("shed-db.password")

	dbtx, err := sqlite3.DB(dbPath, encryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w, %w", ErrNotFound, err)
	}

	queries := db.New(dbtx)

	return &Store{queries: queries, dbtx: dbtx}, nil
}

func NewStore(dbtx db.DBTX) *Store {
	queries := db.New(dbtx)

	return &Store{queries: queries, dbtx: dbtx}
}

type Command struct {
	ID          int64
	Name        string
	Command     string
	Description string
	Parameters  brackets.Parameters
	CreatedAt   string
	UpdatedAt   string
}

func (s *Store) AddCommand(name, command, description string) (*Command, error) {
	if err := validateName(name); err != nil {
		return nil, err
	}

	c, err := brackets.ParseCommand(command)
	if err != nil {
		return nil, fmt.Errorf("failed to parse command: %w", err)
	}

	p, err := brackets.ParseParameters(c)
	if err != nil {
		return nil, fmt.Errorf("failed to parse command parameters: %w", err)
	}

	_, err = s.GetCommandByName(name)
	if err == nil {
		return nil, fmt.Errorf("command with name %q already exists: %w", name, ErrAlreadyExists)
	}

	cmd, err := s.createCommand(name, c, description, p)
	if err != nil {
		return nil, fmt.Errorf("failed to create command: %w", err)
	}

	return cmd, nil
}

func (s *Store) CopyCommand(srcName, destName, jsonValueParams string) (*Command, error) {
	if jsonValueParams == "" {
		jsonValueParams = "{}"
	}

	vp, err := brackets.ValuedParametersFromJSON(jsonValueParams)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrParsingValueParams, err)
	}

	c, err := s.GetCommandByName(srcName)
	if err != nil {
		return nil, fmt.Errorf("failed to get source command: %w", err)
	}

	cmdStr := brackets.HydrateStringSafe(c.Command, vp)

	cmd, err := s.AddCommand(destName, cmdStr, c.Description)
	if err != nil {
		return nil, fmt.Errorf("failed to add copied command: %w", err)
	}

	return cmd, nil
}

func (s *Store) RemoveCommand(name string) error {
	if _, err := s.GetCommandByName(name); err != nil {
		return fmt.Errorf("command %q does not exist: %w", name, ErrCommandNotFound)
	}

	if err := s.queries.DeleteCommandByName(context.Background(), name); err != nil {
		return fmt.Errorf("failed to delete command: %w", err)
	}

	return nil
}

func (s *Store) UpdateCommand(id int64, name, command, description string, params brackets.Parameters) (*Command, error) {
	if err := validateName(name); err != nil {
		return nil, err
	}

	c, err := brackets.ParseCommand(command)
	if err != nil {
		return nil, fmt.Errorf("failed to parse command: %w", err)
	}

	prev, err := s.GetCommand(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing command: %w", err)
	}

	priority, err := brackets.ParseParameters(c) // just to validate
	if err != nil {
		return nil, fmt.Errorf("failed to parse command parameters: %w", err)
	}

	priority.ThreeWayMerge(&prev.Parameters, &params)

	cmd, err := s.updateCommand(id, name, c, description, priority)
	if err != nil {
		return nil, fmt.Errorf("failed to update command: %w", err)
	}

	return cmd, nil
}

func (s *Store) GetCommandByName(name string) (*Command, error) {
	cmd, err := s.queries.GetCommandByName(context.Background(), name)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrCommandNotFound, err)
	}

	return ToCommand(cmd)
}

func (s *Store) GetCommand(id int64) (*Command, error) {
	cmd, err := s.queries.GetCommandByID(context.Background(), id)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrCommandNotFound, err)
	}

	return ToCommand(cmd)
}

func (s *Store) ListCommands() ([]Command, error) {
	cc, err := s.queries.ListCommands(context.Background())
	if err != nil {
		return []Command{}, fmt.Errorf("failed to list commands: %w", err)
	}

	return ToCommands(cc)
}

// validateName checks if a command name is valid.
// Valid names must:
// - Start with a letter (a-z, A-Z)
// - Contain only alphanumeric characters and underscores
// - Not be empty
// - Not exceed the maximum length
func validateName(name string) error {
	if len(name) == 0 {
		return fmt.Errorf("%w: %s", ErrInvalidCommandName, nameLength)
	}

	if len(name) > nameMaxLength {
		return fmt.Errorf("%w: %s", ErrInvalidCommandName, nameLength)
	}

	// Check first character is a letter
	first := name[0]
	if !((first >= 'a' && first <= 'z') || (first >= 'A' && first <= 'Z')) {
		return fmt.Errorf("%w: %s", ErrInvalidCommandName, nameDetails)
	}

	// Check remaining characters are alphanumeric or underscore
	for i := 1; i < len(name); i++ {
		c := name[i]
		isLetter := (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
		isDigit := c >= '0' && c <= '9'
		isUnderscore := c == '_'

		if !isLetter && !isDigit && !isUnderscore {
			return fmt.Errorf("%w: %s", ErrInvalidCommandName, nameDetails)
		}
	}

	return nil
}

func ToParameters(raw json.RawMessage) (brackets.Parameters, error) {
	var params brackets.Parameters
	if err := json.Unmarshal(raw, &params); err != nil {
		return nil, fmt.Errorf("failed to unmarshal parameters: %w", err)
	}

	return params, nil
}

func ToCommand(c db.Command) (*Command, error) {
	params, err := ToParameters(c.Parameters)
	if err != nil {
		return nil, fmt.Errorf("failed to convert to command: %w", err)
	}

	return &Command{
		ID:          c.ID,
		Name:        c.Name,
		Command:     c.Command,
		Description: c.Description,
		Parameters:  params,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}, nil
}

func ToCommands(cc []db.Command) ([]Command, error) {
	mapped := itertools.Map2(slices.Values(cc), func(c db.Command) (Command, error) {
		cmd, err := ToCommand(c)
		if err != nil {
			return Command{}, fmt.Errorf("failed to convert to command: %w", err)
		}

		return *cmd, nil
	})

	var out []Command
	for c, err := range mapped {
		if err != nil {
			return nil, fmt.Errorf("failed to convert commands: %w", err)
		}

		out = append(out, c)
	}

	return out, nil
}

func (s *Store) createCommand(name, command, description string, params brackets.Parameters) (*Command, error) {
	bb, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal parameters to json: %w", err)
	}

	args := db.CreateCommandParams{
		Name:        name,
		Command:     command,
		Description: description,
		Parameters:  bb,
	}

	c, err := s.queries.CreateCommand(context.Background(), args)
	if err != nil {
		return nil, fmt.Errorf("failed to create command: %w", err)
	}

	return ToCommand(c)
}

func (s *Store) updateCommand(id int64, name, command, description string, params brackets.Parameters) (*Command, error) {
	bb, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal parameters to json: %w", err)
	}

	args := db.UpdateCommandParams{
		ID:          id,
		Name:        name,
		Command:     command,
		Description: description,
		Parameters:  bb,
	}

	c, err := s.queries.UpdateCommand(context.Background(), args)
	if err != nil {
		return nil, fmt.Errorf("failed to create command: %w", err)
	}

	return ToCommand(c)
}

func checkNameLength(name string) error {
	if len(name) > nameMaxLength {
		return ErrNameTooLong
	}

	return nil
}
