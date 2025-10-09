package db

import "embed"

// Migrations is an embedded template.
//
//go:embed migrations/*.sql
var Migrations embed.FS
