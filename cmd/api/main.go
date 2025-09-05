package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"os"
	"ticket-api/internal/env"
	"ticket-api/internal/handler"
	_ "ticket-api/internal/handler"
	"ticket-api/internal/repository"
	"time"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/mattn/go-sqlite3"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type application struct {
	port      int
	jwtSecret string
	repos     repository.AppRepositories
	handlers  handler.AppHandlers
}

// @title Ticket API
// @version 1.0
// @BasePath /api/v1
func main() {
	dbSql, err := sql.Open("sqlite3", "./data.db")
	fatalIfErr(err)

	//mongodb
	dbMongo, err := ConnectMongo()
	fatalIfErr(err)

	repos := repository.NewRepositories(dbSql, dbMongo)
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

// ConnectMongo connects to MongoDB and returns the database.
func ConnectMongo() (*mongo.Database, error) {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		return nil, errors.New("MONGODB_URI is not set")
	}

	dbName := env.GetEnvString("MONGODB_DB", "ticketdb")
	opts := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(opts)

	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	log.Println("✅ Connected to MongoDB:", dbName)
	return client.Database(dbName), nil
}
