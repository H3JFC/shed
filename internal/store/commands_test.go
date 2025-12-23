package store

import (
	"errors"
	"testing"

	"github.com/h3jfc/shed/lib/brackets"
)

const (
	validName32Char          = "valid_name_xxxxxxxxxxxxxxxxxxxxx"
	invalidName33CharTooLong = "invalid_name_xxxxxxxxxxxxxxxxxxxx"
)

func init() {
	if len(validName32Char) != 32 {
		panic("validName32Char is not 32 characters long")
	}

	if len(invalidName33CharTooLong) != 33 {
		panic("invalidName33CharTooLong is not 33 characters long")
	}
}

func TestAddCommand_OK(t *testing.T) { // nolint:funlen
	t.Parallel()
	s := prepNewStore(t)

	cmd, err := s.AddCommand("list_files", "ls -la {{path|description}}", "lists files command")
	if err != nil {
		t.Fatalf("unexpected error adding command: %v", err)
	}

	if cmd.Name != "list_files" { //nolint:goconst
		t.Fatalf("expected command name %v, got %v", "list_files", cmd.Name)
	}

	if cmd.Parameters == nil {
		t.Fatalf("expected parameters to be non-nil")
	}

	if len(cmd.Parameters) != 1 {
		t.Fatalf("expected 1 parameter, got %v", len(cmd.Parameters))
	}

	if cmd.Description != "lists files command" {
		t.Fatalf("expected command description %v, got %v", "lists files command", cmd.Description)
	}

	param := cmd.Parameters[0]
	if param.Name != "path" { //nolint:goconst
		t.Fatalf("expected parameter name %v, got %v", "path", param.Name)
	}

	if param.Description != "description" { //nolint:goconst
		t.Fatalf("expected parameter description %v, got %v", "description", param.Description)
	}

	if cmd.Command != "ls -la {{path|description}}" {
		t.Fatalf("expected command %v, got %v", "ls -la", cmd.Command)
	}
}

func TestAddCommand_OKMaxLength(t *testing.T) { // nolint:funlen
	t.Parallel()
	s := prepNewStore(t)

	cmd, err := s.AddCommand(validName32Char, "ls -la", "command with max length name")
	if err != nil {
		t.Fatalf("unexpected error adding command with max length name: %v", err)
	}

	if cmd.Name != validName32Char {
		t.Fatalf("expected command name %v, got %v", validName32Char, cmd.Name)
	}
}

func TestAddCommand_ErrAlreadyExists(t *testing.T) { // nolint:funlen
	t.Parallel()
	s := prepNewStore(t)

	_, err := s.AddCommand("list_files", "ls -la {{path|description}}", "lists files command")
	if err != nil {
		t.Fatalf("unexpected error adding command: %v", err)
	}

	_, err = s.AddCommand("list_files", "ls -la ./", "another command with same name")
	if err == nil {
		t.Fatalf("expected error %v, got nil", ErrAlreadyExists)
	}

	if !errors.Is(err, ErrAlreadyExists) {
		t.Fatalf("expected error %v, got %v", ErrAlreadyExists, err)
	}
}

func TestAddCommand_Err(t *testing.T) { // nolint:funlen
	t.Parallel()

	type testcase struct {
		commandName string
		command     string
		want        error
	}

	tests := map[string]testcase{
		"invalid-starts-with-number": {
			commandName: "123command",
			command:     "ls",
			want:        ErrInvalidCommandName,
		},
		"invalid-starts-with-underscore": {
			commandName: "_command",
			command:     "ls",
			want:        ErrInvalidCommandName,
		},
		"invalid-contains-space": {
			commandName: "list files",
			command:     "ls -la",
			want:        ErrInvalidCommandName,
		},
		"invalid-contains-hyphen": {
			commandName: "list-files",
			command:     "ls",
			want:        ErrInvalidCommandName,
		},
		"invalid-contains-dot": {
			commandName: "list.files",
			command:     "ls",
			want:        ErrInvalidCommandName,
		},
		"invalid-contains-special-chars": {
			commandName: "list@files",
			command:     "ls",
			want:        ErrInvalidCommandName,
		},
		"invalid-empty": {
			commandName: "",
			command:     "ls",
			want:        ErrInvalidCommandName,
		},
		"invalid-too-long-33-chars": {
			commandName: invalidName33CharTooLong,
			command:     "ls",
			want:        ErrInvalidCommandName,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			s := prepNewStore(t)

			_, err := s.AddCommand(tc.commandName, tc.command, "commands with errors")
			if err == nil {
				t.Fatalf("expected error %v, got nil", tc.want)
			}

			if !errors.Is(err, tc.want) {
				t.Fatalf("expected error %v, got %v", tc.want, err)
			}
		})
	}
}

func TestCopyCommand_Err(t *testing.T) { // nolint:funlen
	t.Parallel()

	srcName := "list_files"
	valueParams := `{"path":"/home/user"}`

	type testcase struct {
		destName string
		want     error
	}

	tests := map[string]testcase{
		"invalid-starts-with-number": {
			destName: "123command",
			want:     ErrInvalidCommandName,
		},
		"invalid-starts-with-underscore": {
			destName: "_command",
			want:     ErrInvalidCommandName,
		},
		"invalid-contains-space": {
			destName: "list files",
			want:     ErrInvalidCommandName,
		},
		"invalid-contains-hyphen": {
			destName: "list-files",
			want:     ErrInvalidCommandName,
		},
		"invalid-contains-dot": {
			destName: "list.files",
			want:     ErrInvalidCommandName,
		},
		"invalid-contains-special-chars": {
			destName: "list@files",
			want:     ErrInvalidCommandName,
		},
		"invalid-empty": {
			destName: "",
			want:     ErrInvalidCommandName,
		},
		"invalid-too-long-33-chars": {
			destName: invalidName33CharTooLong,
			want:     ErrInvalidCommandName,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			s := prepNewStore(t)

			_, err := s.AddCommand(srcName, "ls -la {{path|description}}", "lists files command")
			if err != nil {
				t.Fatalf("unexpected error adding command: %v", err)
			}

			_, err = s.CopyCommand(srcName, tc.destName, valueParams)
			if err == nil {
				t.Fatalf("expected error %v, got nil", tc.want)
			}

			if !errors.Is(err, tc.want) {
				t.Fatalf("expected error %v, got %v", tc.want, err)
			}
		})
	}
}

func TestCopyCommand_OK(t *testing.T) { // nolint:funlen
	t.Parallel()
	s := prepNewStore(t)

	if _, err := s.AddCommand("list_files", "ls -la {{path|description}}", ""); err != nil {
		t.Fatalf("unexpected error adding command: %v", err)
	}

	copiedCmd, err := s.CopyCommand("list_files", "list_files_home", `{"path":"/Users/user"}`)
	if err != nil {
		t.Fatalf("unexpected error copying command: %v", err)
	}

	if copiedCmd.Name != "list_files_home" { //nolint:goconst
		t.Fatalf("expected command name %v, got %v", "list_files_home", copiedCmd.Name)
	}

	if copiedCmd.Parameters == nil {
		t.Fatalf("expected parameters to be non-nil")
	}

	if len(copiedCmd.Parameters) != 0 {
		t.Fatalf("expected 0 parameters, got %v", len(copiedCmd.Parameters))
	}
}

func TestCopyCommand_OKMaxLength(t *testing.T) { // nolint:funlen
	t.Parallel()
	s := prepNewStore(t)

	if _, err := s.AddCommand("list_files", "ls -la {{path|description}}", "list files"); err != nil {
		t.Fatalf("unexpected error adding command: %v", err)
	}

	copiedCmd, err := s.CopyCommand("list_files", validName32Char, `{"path":"/Users/user"}`)
	if err != nil {
		t.Fatalf("unexpected error copying command with max length name: %v", err)
	}

	if copiedCmd.Name != validName32Char {
		t.Fatalf("expected command name %v, got %v", validName32Char, copiedCmd.Name)
	}
}

func TestCopyCommand_OKUnusedParams(t *testing.T) { // nolint:funlen
	t.Parallel()
	s := prepNewStore(t)

	if _, err := s.AddCommand("list_files", "ls -la {{path|description}}", "list files"); err != nil {
		t.Fatalf("unexpected error adding command: %v", err)
	}

	copiedCmd, err := s.CopyCommand("list_files", "list_files_home", `{"unused":"not-used-param"}`)
	if err != nil {
		t.Fatalf("unexpected error copying command: %v", err)
	}

	if copiedCmd.Name != "list_files_home" {
		t.Fatalf("expected command name %v, got %v", "list_files_home", copiedCmd.Name)
	}

	if copiedCmd.Parameters == nil {
		t.Fatalf("expected parameters to be non-nil")
	}

	if len(copiedCmd.Parameters) != 1 {
		t.Fatalf("expected 1 parameters, got %v", len(copiedCmd.Parameters))
	}

	if copiedCmd.Parameters[0].Name != "path" {
		t.Fatalf("expected parameter name %v, got %v", "path", copiedCmd.Parameters[0].Name)
	}

	if copiedCmd.Parameters[0].Description != "description" {
		t.Fatalf("expected parameter description %v, got %v", "description", copiedCmd.Parameters[0].Description)
	}
}

func TestCopyCommand_OKEmptyParams(t *testing.T) { // nolint:funlen
	t.Parallel()
	s := prepNewStore(t)

	if _, err := s.AddCommand("list_files", "ls -la {{path|description}}", "list files"); err != nil {
		t.Fatalf("unexpected error adding command: %v", err)
	}

	copiedCmd, err := s.CopyCommand("list_files", "list_files_home", ``)
	if err != nil {
		t.Fatalf("unexpected error copying command: %v", err)
	}

	if copiedCmd.Name != "list_files_home" {
		t.Fatalf("expected command name %v, got %v", "list_files_home", copiedCmd.Name)
	}

	if copiedCmd.Parameters == nil {
		t.Fatalf("expected parameters to be non-nil")
	}

	if len(copiedCmd.Parameters) != 1 {
		t.Fatalf("expected 1 parameters, got %v", len(copiedCmd.Parameters))
	}

	if copiedCmd.Parameters[0].Name != "path" {
		t.Fatalf("expected parameter name %v, got %v", "path", copiedCmd.Parameters[0].Name)
	}

	if copiedCmd.Parameters[0].Description != "description" {
		t.Fatalf("expected parameter description %v, got %v", "description", copiedCmd.Parameters[0].Description)
	}
}

func TestCopyCommand_ErrParsing(t *testing.T) { // nolint:funlen
	t.Parallel()
	s := prepNewStore(t)

	if _, err := s.AddCommand("list_files", "ls -la {{path|description}}", "list files"); err != nil {
		t.Fatalf("unexpected error adding command: %v", err)
	}

	_, err := s.CopyCommand("list_files", "list_files_home", `{"invalid":not-valid-json`)
	if !errors.Is(err, ErrParsingValueParams) {
		t.Fatalf("expected error %v, got %v", ErrParsingValueParams, err)
	}
}

func TestCopyCommand_DoesNotExist(t *testing.T) { // nolint:funlen
	t.Parallel()
	s := prepNewStore(t)

	_, err := s.CopyCommand("list_files", "list_files_home", `{}`)
	if !errors.Is(err, ErrCommandNotFound) {
		t.Fatalf("expected error %v, got %v", ErrCommandNotFound, err)
	}
}

func TestRemoveCommand_OK(t *testing.T) { // nolint:funlen
	t.Parallel()
	s := prepNewStore(t)

	if _, err := s.AddCommand("list_files", "ls -la {{path|description}}", "list files"); err != nil {
		t.Fatalf("unexpected error adding command: %v", err)
	}

	err := s.RemoveCommand("list_files")
	if err != nil {
		t.Fatalf("unexpected error copying command: %v", err)
	}
}

func TestRemoveCommand_ErrCommandNotFound(t *testing.T) { // nolint:funlen
	t.Parallel()
	s := prepNewStore(t)

	err := s.RemoveCommand("does_not_exist")
	if !errors.Is(err, ErrCommandNotFound) {
		t.Fatalf("expected error %v, got %v", ErrCommandNotFound, err)
	}
}

func TestUpdateCommand_OK(t *testing.T) { // nolint:funlen
	t.Parallel()
	s := prepNewStore(t)

	// Create initial command
	cmd, err := s.AddCommand("list_files", "ls -la {{path|description}}", "list files")
	if err != nil {
		t.Fatalf("unexpected error adding command: %v", err)
	}

	// Update command with new parameters
	params := brackets.Parameters{
		{Name: "path", Description: "updated description"},
	}

	updatedCmd, err := s.UpdateCommand(
		cmd.ID,
		"list_files",
		"ls -la {{path|updated description}}",
		"updated list files command",
		params,
		"{}",
	)
	if err != nil {
		t.Fatalf("unexpected error updating command: %v", err)
	}

	if updatedCmd.Name != "list_files" {
		t.Fatalf("expected command name %v, got %v", "list_files", updatedCmd.Name)
	}

	if updatedCmd.Command != "ls -la {{path|updated description}}" {
		t.Fatalf("expected command %v, got %v", "ls -la {{path|updated description}}", updatedCmd.Command)
	}

	if len(updatedCmd.Parameters) != 1 {
		t.Fatalf("expected 1 parameter, got %v", len(updatedCmd.Parameters))
	}

	if updatedCmd.Parameters[0].Description != "updated description" {
		t.Fatalf("expected parameter description %v, got %v", "updated description", updatedCmd.Parameters[0].Description)
	}
}

func TestUpdateCommand_ChangeName(t *testing.T) { // nolint:funlen
	t.Parallel()
	s := prepNewStore(t)

	// Create initial command
	cmd, err := s.AddCommand("list_files", "ls -la {{path|description}}", "list files")
	if err != nil {
		t.Fatalf("unexpected error adding command: %v", err)
	}

	// Update command name
	updatedCmd, err := s.UpdateCommand(
		cmd.ID,
		"show_files",
		"ls -la {{path|description}}",
		"show files command",
		cmd.Parameters,
		"{}",
	)
	if err != nil {
		t.Fatalf("unexpected error updating command: %v", err)
	}

	if updatedCmd.Name != "show_files" {
		t.Fatalf("expected command name %v, got %v", "show_files", updatedCmd.Name)
	}
}

func TestUpdateCommand_OKMaxLength(t *testing.T) { // nolint:funlen
	t.Parallel()
	s := prepNewStore(t)

	// Create initial command
	cmd, err := s.AddCommand("list_files", "ls -la {{path|description}}", "list files")
	if err != nil {
		t.Fatalf("unexpected error adding command: %v", err)
	}

	// Update command with max length name
	updatedCmd, err := s.UpdateCommand(
		cmd.ID,
		validName32Char,
		"ls -la {{path|description}}",
		"command with max length name",
		cmd.Parameters,
		"{}",
	)
	if err != nil {
		t.Fatalf("unexpected error updating command with max length name: %v", err)
	}

	if updatedCmd.Name != validName32Char {
		t.Fatalf("expected command name %v, got %v", validName32Char, updatedCmd.Name)
	}
}

func TestUpdateCommand_AddNewParameter(t *testing.T) { // nolint:funlen
	t.Parallel()
	s := prepNewStore(t)

	// Create initial command with one parameter
	cmd, err := s.AddCommand("list_files", "ls -la {{path|description}}", "list files")
	if err != nil {
		t.Fatalf("unexpected error adding command: %v", err)
	}

	// Update command to add a new parameter
	newParams := brackets.Parameters{
		{Name: "path", Description: "description"},
		{Name: "filter", Description: "file filter"},
	}

	updatedCmd, err := s.UpdateCommand(
		cmd.ID,
		"list_files",
		"ls -la {{path|description}} {{filter|file filter}}",
		"list files with filter",
		newParams,
		"{}",
	)
	if err != nil {
		t.Fatalf("unexpected error updating command: %v", err)
	}

	if len(updatedCmd.Parameters) != 2 {
		t.Fatalf("expected 2 parameters, got %v", len(updatedCmd.Parameters))
	}

	// Parameters should be sorted by name
	if updatedCmd.Parameters[0].Name != "filter" {
		t.Fatalf("expected first parameter name %v, got %v", "filter", updatedCmd.Parameters[0].Name)
	}

	if updatedCmd.Parameters[1].Name != "path" {
		t.Fatalf("expected second parameter name %v, got %v", "path", updatedCmd.Parameters[1].Name)
	}
}

func TestUpdateCommand_ErrCommandNotFound(t *testing.T) { // nolint:funlen
	t.Parallel()
	s := prepNewStore(t)

	_, err := s.UpdateCommand(999, "list_files", "ls -la", "description", brackets.Parameters{}, "{}")
	if !errors.Is(err, ErrCommandNotFound) {
		t.Fatalf("expected error %v, got %v", ErrCommandNotFound, err)
	}
}

func TestUpdateCommand_ErrInvalidCommandName(t *testing.T) { // nolint:funlen
	t.Parallel()

	type testcase struct {
		commandName string
		want        error
	}

	tests := map[string]testcase{
		"invalid-starts-with-number": {
			commandName: "123command",
			want:        ErrInvalidCommandName,
		},
		"invalid-starts-with-underscore": {
			commandName: "_command",
			want:        ErrInvalidCommandName,
		},
		"invalid-contains-space": {
			commandName: "list files",
			want:        ErrInvalidCommandName,
		},
		"invalid-contains-hyphen": {
			commandName: "list-files",
			want:        ErrInvalidCommandName,
		},
		"invalid-contains-dot": {
			commandName: "list.files",
			want:        ErrInvalidCommandName,
		},
		"invalid-contains-special-chars": {
			commandName: "list@files",
			want:        ErrInvalidCommandName,
		},
		"invalid-empty": {
			commandName: "",
			want:        ErrInvalidCommandName,
		},
		"invalid-too-long-33-chars": {
			commandName: invalidName33CharTooLong,
			want:        ErrInvalidCommandName,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			s := prepNewStore(t)

			// Create initial command
			cmd, err := s.AddCommand("list_files", "ls -la", "list files")
			if err != nil {
				t.Fatalf("unexpected error adding command: %v", err)
			}

			// Attempt update with invalid name
			_, err = s.UpdateCommand(cmd.ID, tc.commandName, "ls -la", "description", brackets.Parameters{}, "{}")
			if err == nil {
				t.Fatalf("expected error %v, got nil", tc.want)
			}

			if !errors.Is(err, tc.want) {
				t.Fatalf("expected error %v, got %v", tc.want, err)
			}
		})
	}
}

func TestUpdateCommand_ThreeWayMerge(t *testing.T) { // nolint:funlen
	t.Parallel()
	s := prepNewStore(t)

	// Create initial command with parameter
	cmd, err := s.AddCommand("test_cmd", "echo {{param1|original desc}}", "test command")
	if err != nil {
		t.Fatalf("unexpected error adding command: %v", err)
	}

	// Update with longer description in the updated params
	updatedParams := brackets.Parameters{
		{Name: "param1", Description: "much longer updated description"},
	}

	updatedCmd, err := s.UpdateCommand(
		cmd.ID,
		"test_cmd",
		"echo {{param1|new command desc}}",
		"updated test command",
		updatedParams,
		"{}",
	)
	if err != nil {
		t.Fatalf("unexpected error updating command: %v", err)
	}

	// The longer description should be chosen
	if len(updatedCmd.Parameters) != 1 {
		t.Fatalf("expected 1 parameter, got %v", len(updatedCmd.Parameters))
	}

	if updatedCmd.Parameters[0].Description != "much longer updated description" {
		t.Fatalf("expected parameter description %v, got %v", "much longer updated description", updatedCmd.Parameters[0].Description) //nolint:lll
	}
}

func TestUpdateCommand_WithHydration(t *testing.T) { // nolint:funlen
	t.Parallel()
	s := prepNewStore(t)

	// Create initial command with parameters
	cmd, err := s.AddCommand("api_call", "curl -XGET {{url|api url}} -H {{header|auth header}}", "make api call")
	if err != nil {
		t.Fatalf("unexpected error adding command: %v", err)
	}

	// Update command and hydrate the url parameter
	params := brackets.Parameters{
		{Name: "header", Description: "auth header"},
	}

	updatedCmd, err := s.UpdateCommand(
		cmd.ID,
		"api_call",
		"curl -XGET {{url|api url}} -H {{header|auth header}}",
		"make api call",
		params,
		`{"url":"https://api.example.com"}`,
	)
	if err != nil {
		t.Fatalf("unexpected error updating command: %v", err)
	}

	// The command should have the url hydrated
	if updatedCmd.Command != "curl -XGET https://api.example.com -H {{header|auth header}}" {
		t.Fatalf("expected command %q, got %q",
			"curl -XGET https://api.example.com -H {{header|auth header}}",
			updatedCmd.Command)
	}

	// Should only have one parameter left (header)
	if len(updatedCmd.Parameters) != 1 {
		t.Fatalf("expected 1 parameter, got %v", len(updatedCmd.Parameters))
	}

	if updatedCmd.Parameters[0].Name != "header" {
		t.Fatalf("expected parameter name %q, got %q", "header", updatedCmd.Parameters[0].Name)
	}
}

func TestUpdateCommand_WithHydrationMultiple(t *testing.T) { // nolint:funlen
	t.Parallel()
	s := prepNewStore(t)

	// Create initial command
	cmd, err := s.AddCommand(
		"post_data",
		"curl -XPOST --data '{{data}}' -H {{auth}} {{url}}",
		"post data to api",
	)
	if err != nil {
		t.Fatalf("unexpected error adding command: %v", err)
	}

	// Update and hydrate multiple parameters
	updatedCmd, err := s.UpdateCommand(
		cmd.ID,
		"post_data",
		"curl -XPOST --data '{{data}}' -H {{auth}} {{url}}",
		"post data to api",
		brackets.Parameters{},
		`{"data":"{\"foo\":\"bar\"}","url":"https://api.example.com"}`,
	)
	if err != nil {
		t.Fatalf("unexpected error updating command: %v", err)
	}

	// The command should have data and url hydrated
	expectedCmd := `curl -XPOST --data '{"foo":"bar"}' -H {{auth}} https://api.example.com`
	if updatedCmd.Command != expectedCmd {
		t.Fatalf("expected command %q, got %q", expectedCmd, updatedCmd.Command)
	}

	// Should only have auth parameter left
	if len(updatedCmd.Parameters) != 1 {
		t.Fatalf("expected 1 parameter, got %v", len(updatedCmd.Parameters))
	}

	if updatedCmd.Parameters[0].Name != "auth" {
		t.Fatalf("expected parameter name %q, got %q", "auth", updatedCmd.Parameters[0].Name)
	}
}

func TestUpdateCommand_WithEmptyString(t *testing.T) { // nolint:funlen
	t.Parallel()
	s := prepNewStore(t)

	// Create initial command
	cmd, err := s.AddCommand("test_cmd", "echo {{param|test parameter}}", "test command")
	if err != nil {
		t.Fatalf("unexpected error adding command: %v", err)
	}

	// Update with empty string should work the same as empty JSON object
	updatedCmd, err := s.UpdateCommand(
		cmd.ID,
		"test_cmd",
		"echo {{param|test parameter}}",
		"test command updated",
		brackets.Parameters{{Name: "param", Description: "test parameter"}},
		"",
	)
	if err != nil {
		t.Fatalf("unexpected error updating command with empty string: %v", err)
	}

	// Command should remain unchanged (no hydration)
	if updatedCmd.Command != "echo {{param|test parameter}}" {
		t.Fatalf("expected command %q, got %q", "echo {{param|test parameter}}", updatedCmd.Command)
	}

	// Should still have the parameter
	if len(updatedCmd.Parameters) != 1 {
		t.Fatalf("expected 1 parameter, got %v", len(updatedCmd.Parameters))
	}

	if updatedCmd.Parameters[0].Name != "param" {
		t.Fatalf("expected parameter name %q, got %q", "param", updatedCmd.Parameters[0].Name)
	}
}

func TestUpdateCommand_ErrInvalidJSON(t *testing.T) { // nolint:funlen
	t.Parallel()
	s := prepNewStore(t)

	// Create initial command
	cmd, err := s.AddCommand("test_cmd", "echo {{param}}", "test command")
	if err != nil {
		t.Fatalf("unexpected error adding command: %v", err)
	}

	// Attempt update with invalid JSON
	_, err = s.UpdateCommand(
		cmd.ID,
		"test_cmd",
		"echo {{param}}",
		"test command",
		brackets.Parameters{},
		`{invalid json}`,
	)
	if err == nil {
		t.Fatalf("expected error for invalid JSON, got nil")
	}

	if !errors.Is(err, brackets.ErrParsingValueParams) {
		t.Fatalf("expected ErrParsingValueParams, got %v", err)
	}
}

func TestGetCommand_OK(t *testing.T) { // nolint:funlen,cyclop
	t.Parallel()
	s := prepNewStore(t)

	// Create a command
	addedCmd, err := s.AddCommand("list_files", "ls -la {{path|directory path}}", "list files in directory")
	if err != nil {
		t.Fatalf("unexpected error adding command: %v", err)
	}

	// Retrieve command by ID
	cmd, err := s.GetCommand(addedCmd.ID)
	if err != nil {
		t.Fatalf("unexpected error getting command: %v", err)
	}

	if cmd.ID != addedCmd.ID {
		t.Fatalf("expected command ID %v, got %v", addedCmd.ID, cmd.ID)
	}

	if cmd.Name != "list_files" {
		t.Fatalf("expected command name %v, got %v", "list_files", cmd.Name)
	}

	if cmd.Command != "ls -la {{path|directory path}}" {
		t.Fatalf("expected command %v, got %v", "ls -la {{path|directory path}}", cmd.Command)
	}

	if len(cmd.Parameters) != 1 {
		t.Fatalf("expected 1 parameter, got %v", len(cmd.Parameters))
	}

	if cmd.Parameters[0].Name != "path" {
		t.Fatalf("expected parameter name %v, got %v", "path", cmd.Parameters[0].Name)
	}

	if cmd.Parameters[0].Description != "directory path" {
		t.Fatalf("expected parameter description %v, got %v", "directory path", cmd.Parameters[0].Description)
	}

	if cmd.CreatedAt == "" {
		t.Fatalf("expected CreatedAt to be set")
	}

	if cmd.UpdatedAt == "" {
		t.Fatalf("expected UpdatedAt to be set")
	}
}

func TestGetCommand_ErrCommandNotFound(t *testing.T) { // nolint:funlen
	t.Parallel()
	s := prepNewStore(t)

	_, err := s.GetCommand(999)
	if !errors.Is(err, ErrCommandNotFound) {
		t.Fatalf("expected error %v, got %v", ErrCommandNotFound, err)
	}
}

func TestGetCommandByName_OK(t *testing.T) { // nolint:funlen
	t.Parallel()
	s := prepNewStore(t)

	// Create a command
	_, err := s.AddCommand("list_files", "ls -la {{path|directory path}}", "list files in directory")
	if err != nil {
		t.Fatalf("unexpected error adding command: %v", err)
	}

	// Retrieve command by name
	cmd, err := s.GetCommandByName("list_files")
	if err != nil {
		t.Fatalf("unexpected error getting command by name: %v", err)
	}

	if cmd.Name != "list_files" {
		t.Fatalf("expected command name %v, got %v", "list_files", cmd.Name)
	}

	if cmd.Command != "ls -la {{path|directory path}}" {
		t.Fatalf("expected command %v, got %v", "ls -la {{path|directory path}}", cmd.Command)
	}

	if len(cmd.Parameters) != 1 {
		t.Fatalf("expected 1 parameter, got %v", len(cmd.Parameters))
	}

	if cmd.Parameters[0].Name != "path" {
		t.Fatalf("expected parameter name %v, got %v", "path", cmd.Parameters[0].Name)
	}

	if cmd.Parameters[0].Description != "directory path" {
		t.Fatalf("expected parameter description %v, got %v", "directory path", cmd.Parameters[0].Description)
	}
}

func TestGetCommandByName_ErrCommandNotFound(t *testing.T) { // nolint:funlen
	t.Parallel()
	s := prepNewStore(t)

	_, err := s.GetCommandByName("does_not_exist")
	if !errors.Is(err, ErrCommandNotFound) {
		t.Fatalf("expected error %v, got %v", ErrCommandNotFound, err)
	}
}

func TestListCommands_OK(t *testing.T) { // nolint:funlen,cyclop
	t.Parallel()
	s := prepNewStore(t)

	// Create multiple commands
	_, err := s.AddCommand("list_files", "ls -la {{path|directory path}}", "list files in directory")
	if err != nil {
		t.Fatalf("unexpected error adding command: %v", err)
	}

	_, err = s.AddCommand("show_date", "date +%Y-%m-%d", "display current date")
	if err != nil {
		t.Fatalf("unexpected error adding command: %v", err)
	}

	_, err = s.AddCommand("grep_logs", "grep {{pattern|search pattern}} {{file|log file}}", "search logs for pattern")
	if err != nil {
		t.Fatalf("unexpected error adding command: %v", err)
	}

	// List all commands
	commands, err := s.ListCommands()
	if err != nil {
		t.Fatalf("unexpected error listing commands: %v", err)
	}

	if len(commands) != 3 {
		t.Fatalf("expected 3 commands, got %v", len(commands))
	}

	// Verify command names exist
	names := make(map[string]bool)
	for _, cmd := range commands {
		names[cmd.Name] = true
	}

	expectedNames := []string{"list_files", "show_date", "grep_logs"}
	for _, name := range expectedNames {
		if !names[name] {
			t.Fatalf("expected command %v to be in list", name)
		}
	}

	// Verify grep_logs has correct parameters
	var grepCmd *Command

	for i := range commands {
		if commands[i].Name == "grep_logs" {
			grepCmd = &commands[i]

			break
		}
	}

	if grepCmd == nil {
		t.Fatalf("expected to find grep_logs command")
	}

	if len(grepCmd.Parameters) != 2 {
		t.Fatalf("expected 2 parameters for grep_logs, got %v", len(grepCmd.Parameters))
	}

	// Parameters should be sorted by name
	if grepCmd.Parameters[0].Name != "file" {
		t.Fatalf("expected first parameter name %v, got %v", "file", grepCmd.Parameters[0].Name)
	}

	if grepCmd.Parameters[1].Name != "pattern" {
		t.Fatalf("expected second parameter name %v, got %v", "pattern", grepCmd.Parameters[1].Name)
	}
}

func TestListCommands_Empty(t *testing.T) { // nolint:funlen
	t.Parallel()
	s := prepNewStore(t)

	// List commands when none exist
	commands, err := s.ListCommands()
	if err != nil {
		t.Fatalf("unexpected error listing commands: %v", err)
	}

	if len(commands) != 0 {
		t.Fatalf("expected 0 commands, got %v", len(commands))
	}
}
