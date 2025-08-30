package main

import (
	"database/sql"
	"log"
	"ticket-api/internal/env"
	"ticket-api/internal/handler"
	_ "ticket-api/internal/handler"
	"ticket-api/internal/repository"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/mattn/go-sqlite3"
)

type application struct {
	port      int
	jwtSecret string
	repos     repository.AppRepositories
	handlers  handler.AppHandlers
}

func main() {
	db, err := sql.Open("sqlite3", "./data.db")

	fatalIfErr(err)

	repos := repository.NewRepositories(db)
	handlers := handler.NewAppHandlers(repos)

	app := &application{
		port:      env.GetEnvInt("PORT", 8080),
		jwtSecret: env.GetEnvString("JWT_SECRET", "super-storng-key-123456"),
		repos:     *repos,
		handlers:  *handlers,
	}

	if err := app.serve(); err != nil {
		log.Fatal(err)
	}
}

// fatalIfErr logs and exits if err is not nil
func fatalIfErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
