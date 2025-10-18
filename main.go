package main

import (
	"fmt"
	"runtime"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"

	"h3jfc/shed/cmd"
	"h3jfc/shed/lib/sqlite3"
)

const (
	defaultEncryptionKey = "my_secret_key" // In a real application, use a secure method to manage encryption keys.
)

func main() {
	fmt.Println(runtime.GOOS)

	if err := sqlite3.Migrate(defaultEncryptionKey); err != nil {
		return
	}

	cmd.Execute()
}
