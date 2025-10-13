package main

import (
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"

	"h3jfc/shed/cmd"
	"h3jfc/shed/lib/sqlite3"
)

func main() {
	if err := sqlite3.Migrate(); err != nil {
		panic(err)
	}

	cmd.Execute()
}
