package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Please provide a command: 'up' or 'down'")
	}
	action := os.Args[1]

	// Open DB
	db, err := sql.Open("sqlite3", "./data.db")
	fatalIfErr(err)
	defer db.Close()

	// Migration driver
	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	fatalIfErr(err)

	// Migration source (filesystem)
	src, err := (&file.File{}).Open("cmd/migrate/migrations")
	fatalIfErr(err)

	// Migrate instance
	m, err := migrate.NewWithInstance("file", src, "sqlite3", driver)
	fatalIfErr(err)

	// Execute migration
	switch action {
	case "up":
		err = m.Up()
		if err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
		fmt.Println("Migration applied successfully!")

	case "down":
		err = m.Steps(-1)
		if err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
		fmt.Println("Rolled back one migration!")

	case "downall":
		err = m.Down()
		if err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
		fmt.Println("Rolled back all migrations!")

	default:
		log.Fatal("Unknown command. Use 'up' or 'down'")
	}
}

// fatalIfErr logs and exits if err is not nil
func fatalIfErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
