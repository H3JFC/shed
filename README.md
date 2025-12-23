# Shed

<p align="center">
  <a href="https://github.com/H3JFC/shed/actions/workflows/main.yml"><img src="https://github.com/H3JFC/shed/actions/workflows/main.yml/badge.svg?branch=main" alt="Go"></a>
  <a href="https://goreportcard.com/report/github.com/h3jfc/shed"><img src="https://goreportcard.com/badge/github.com/h3jfc/shed" alt="Go Report Card"></a>
  <a href="https://pkg.go.dev/github.com/h3jfc/shed"><img src="https://pkg.go.dev/badge/github.com/h3jfc/shed.svg" alt="Go Reference"></a>
</p>

A secure CLI tool for storing and executing commonly used commands with parameter templating and secrets management.

## Overview

Shed is a command-line interface (CLI) tool designed to help developers and system administrators store, organize, and execute frequently used commands. It features parameter templating, integrated secrets management, and encrypted storage using SQLCipher.

### Key Features

- **Command Storage**: Save commonly used commands with descriptive names
- **Parameter Templating**: Define parameters in commands using `{{name|description}}` syntax
- **Secrets Management**: Store sensitive information securely and reference with `{{!key}}` syntax
- **Encrypted Database**: All data is stored in an encrypted SQLite database using SQLCipher
- **Interactive Prompts**: Automatically prompts for parameter values when running commands
- **Command Operations**: Add, list, edit, copy, describe, and remove commands
- **Cross-Platform**: Supports Linux, macOS, and Windows

## Installation

### Pre-built Binaries

Download the latest release for your platform from the [Releases](https://github.com/H3JFC/shed/releases) page:

- `shed-linux-amd64` - Linux (x86_64)
- `shed-darwin-amd64` - macOS (Intel)
- `shed-darwin-arm64` - macOS (Apple Silicon)
- `shed-windows-amd64.exe` - Windows (x86_64)

Make the binary executable (Linux/macOS):

```bash
chmod +x shed-*
sudo mv shed-* /usr/local/bin/shed
```

### Build from Source

Requirements:

- Go 1.25.1 or later
- GCC (for CGO/SQLCipher support)
- SQLCipher development libraries

```bash
# Clone the repository
git clone https://github.com/H3JFC/shed.git
cd shed

# Build
go build -tags="sqlcipher,linux" -o shed main.go

# Install (optional)
sudo mv shed /usr/local/bin/
```

Platform-specific build tags:

- Linux: `-tags="sqlcipher,linux"`
- macOS: `-tags="sqlcipher,darwin"`
- Windows: `-tags="sqlcipher,windows"`

## Quick Start

```bash
# Initialize shed (creates configuration and database)
shed init

# Add a simple command
shed add hello "echo 'Hello, World!'"

# Add a command with parameters
shed add greet "echo 'Hello, {{name|person's name}}!'" -d "Greet someone"

# List all commands
shed list

# Run a command (will prompt for parameters)
shed run greet

# Add a secret
shed secret add api_key

# Add a command using a secret
shed add deploy "curl -H 'Authorization: Bearer {{!api_key}}' https://api.example.com/deploy"
```

## Usage

### Commands

#### `shed init`

Initialize shed configuration and database.

```bash
shed init
```

#### `shed add <name> <command>`

Add a new command to shed.

```bash
# Basic command
shed add list_files "ls -la"

# Command with parameters
shed add list_dir "ls -la {{path|directory path}}" -d "List directory contents"

# Command with multiple parameters
shed add git_commit "git add . && git commit -m '{{message|commit message}}' && git push {{branch|branch name}}"

# Command with secrets
shed add deploy "kubectl apply -f {{file|manifest file}} --token={{!k8s_token}}"
```

**Parameter Syntax**: `{{name|description}}`

- `name`: Parameter identifier (used internally)
- `description`: Optional human-readable description shown in prompts

**Secret Syntax**: `{{!key}}`

- `key`: The secret key stored in shed

Options:

- `-d, --description`: Description of the command

#### `shed list`

List all stored commands.

```bash
shed list
```

Output includes:

- Command name
- Command string
- Description
- Parameters (with descriptions)
- Created/Updated timestamps

#### `shed run <name>`

Execute a stored command.

```bash
shed run greet
# Prompts: Enter value for name (person's name):
# Executes: echo 'Hello, John!'
```

Shed will interactively prompt for any parameters or secrets needed by the command.

#### `shed describe <name>`

Show detailed information about a command.

```bash
shed describe git_commit
```

#### `shed edit <name>`

Edit an existing command.

```bash
shed edit greet
# Opens editor to modify command and description
```

#### `shed cp <source> <destination>`

Copy a command to a new name.

```bash
shed cp greet welcome
```

#### `shed rm <name>`

Remove a command.

```bash
shed rm old_command
```

### Secret Management

Secrets are stored encrypted in the database and can be referenced in commands.

#### `shed secret add <key>`

Add a new secret.

```bash
shed secret add github_token -d "GitHub Personal Access Token"
# Prompts for secret value (input hidden)
```

Options:

- `-d, --description`: Description of the secret

#### `shed secret list`

List all secrets (values are hidden).

```bash
shed secret list
```

#### `shed secret edit <key>`

Update a secret's value or description.

```bash
shed secret edit github_token
```

#### `shed secret rm <key>`

Remove a secret.

```bash
shed secret rm old_api_key
```

## Configuration

Shed looks for configuration in the following locations (in order):

1. `$SHED_DIR` environment variable
2. `~/.config/shed/` (Linux/macOS)
3. `~/Library/Application Support/shed/` (macOS)
4. `%APPDATA%\shed\` (Windows)

### Configuration File

The configuration file is `config.toml`:

```toml
[shed-db]
location = "/path/to/shed.db"
password = "encryption-key"
```

### Environment Variables

- `SHED_DIR`: Override default configuration directory
- `SHED_SHED_DB_LOCATION`: Override database location
- `SHED_SHED_DB_PASSWORD`: Override database encryption key

### Command-Line Flags

Global flags:

- `--shed-dir`: Path to shed configuration directory
- `-v, --verbose`: Enable verbose logging

## Architecture

### Database Schema

Shed uses SQLCipher (encrypted SQLite) with the following schema:

- **commands**: Stores command definitions
  - id, name, command, description, created_at, updated_at
- **parameters**: Stores command parameters
  - id, command_id, name, description, position
- **secrets**: Stores encrypted secrets
  - id, key, value (encrypted), description, created_at, updated_at

### Security

- **Encryption**: All data is encrypted at rest using SQLCipher
- **Secret Storage**: Secrets are double-encrypted within the database
- **Password Protection**: Database requires an encryption key
- **No Plaintext**: Secrets are never stored or logged in plaintext

## Development

### Requirements

- Go 1.25.1+
- GCC/Clang (for CGO)
- SQLCipher development libraries
- Make (optional, for using Makefile)

### Setup

```bash
# Clone repository
git clone https://github.com/H3JFC/shed.git
cd shed

# Install dependencies
go mod download

# Run tests
make test

# Run tests with coverage
make test-coverage

# Build
make build

# Run linter
make lint
```

### Project Structure

```
shed/
├── cmd/                   # Command definitions
│   ├── command/           # Command management commands
│   ├── secret/            # Secret management commands
│   ├── init.go            # Initialization command
│   └── root.go            # Root command and CLI setup
├── internal/              # Internal packages
│   ├── commands/          # Command execution logic
│   ├── config/            # Configuration management
│   ├── execute/           # Command execution engine
│   ├── logger/            # Logging utilities
│   └── store/             # Database operations
├── migrations/            # Database migrations
├── main.go               # Application entry point
└── go.mod                # Go module definition
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test
go test -v ./internal/store/...
```

### Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Troubleshooting

### SQLCipher Not Found

If you get SQLCipher errors during build:

**Linux**:

```bash
sudo apt-get install libsqlcipher-dev
```

**macOS**:

```bash
brew install sqlcipher
```

**Windows**:
Use MSYS2:

```bash
pacman -S mingw-w64-x86_64-sqlcipher
```

### Configuration Not Found

Ensure shed is initialized:

```bash
shed init
```

Or set the `SHED_DIR` environment variable:

```bash
export SHED_DIR="/path/to/config"
```

## License

This project is dual-licensed:

- **AGPL-3.0**: For open-source use (see [LICENSE.md](LICENSE.md))
- **Commercial License**: For proprietary use (see [COMMERCIAL-LICENSE.md](COMMERCIAL-LICENSE.md))

## Acknowledgments

- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [Viper](https://github.com/spf13/viper) - Configuration management
- [SQLCipher](https://www.zetetic.net/sqlcipher/) - Encrypted SQLite
- [golang-migrate](https://github.com/golang-migrate/migrate) - Database migrations

## Status

⚠️ **Early Development**: Shed is in active development. Use at your own risk.

## Contact

- GitHub: [@H3JFC](https://github.com/H3JFC)
- Issues: [GitHub Issues](https://github.com/H3JFC/shed/issues)
