package main

import (
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/h3jfc/shed/cmd"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	cmd.Execute()
}
