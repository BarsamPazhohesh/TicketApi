package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"os"
	"strconv"
	"ticket-api/internal/config"
	"ticket-api/internal/env"
	"ticket-api/internal/errx"
	"ticket-api/internal/handler"
	"ticket-api/internal/repository"
	"ticket-api/internal/services"
	"time"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/mattn/go-sqlite3"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type application struct {
	port     int
	mongo    *mongo.Database
	sql      *sql.DB
	redis    *redis.Client
	services *services.AppServices
	repos    *repository.AppRepositories
	handlers *handler.AppHandlers
}

// @title Ticket API
// @version 1.0
// @BasePath /api/v1
func main() {
	config.Load("config.yaml")
	dbSql, err := sql.Open("sqlite3", "file:./data.db?_foreign_keys=on")
	fatalIfErr(err)
	errx.NewRegistry(dbSql)

	defer dbSql.Close()

	// mongodb
	var dbMongo *mongo.Database = nil
	if config.Get().Mongo.Enable {
		dbMongo, err = ConnectMongo()
		fatalIfErr(err)
	}

	dbRedis, err := ConnectRedis()
	fatalIfErr(err)

	services := services.NewAppService()
	repos := repository.NewRepositories(dbSql, dbMongo)
	handlers := handler.NewAppHandlers(repos, services)

	app := &application{
		port:     config.Get().App.Port,
		sql:      dbSql,
		mongo:    dbMongo,
		redis:    dbRedis,
		services: services,
		repos:    repos,
		handlers: handlers,
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

	dbName := env.GetEnvString("MONGODB_DB", config.Get().Mongo.DBName)
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

// ConnectRedis connects to Redis and returns the client.
func ConnectRedis() (*redis.Client, error) {
	redisCfg := config.Get().Redis

	addr := redisCfg.Host + ":" + strconv.Itoa(redisCfg.Port)
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: env.GetEnvString("REDIS_PASSWORD", ""),
		DB:       redisCfg.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Test the connection
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	log.Println("✅ Connected to Redis:", addr, "DB:", redisCfg.DB)
	return rdb, nil
}
