package main

import (
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"

	"h3jfc/shed/cmd"
)

func main() {
	cmd.Execute()
}
