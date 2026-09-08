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
	"ticket-api/internal/security"
	"ticket-api/internal/services"
	"time"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/mattn/go-sqlite3"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type application struct {
	port     int
	mongo    *mongo.Database
	sql      *sql.DB
	redis    *redis.Client
	minio    *minio.Client
	services *services.AppServices
	repos    *repository.AppRepositories
	handlers *handler.AppHandlers
	security *security.SecurityRegistry
}

// @title Ticket API
// @version 1.0
// @BasePath /api/v1
func main() {
	config.Load("config.yaml")
	dbSQL, err := sql.Open("sqlite3", "file:./data.db?_foreign_keys=on")
	fatalIfErr(err)
	errx.NewRegistry(dbSQL)

	defer dbSQL.Close()

	// mongodb
	var dbMongo *mongo.Database = nil
	if config.Get().Mongo.Enable {
		dbMongo, err = ConnectMongo()
		fatalIfErr(err)
	}

	dbRedis, err := ConnectRedis()
	fatalIfErr(err)

	// MinIO
	var minioClient *minio.Client = nil
	if config.Get().Minio.Enable {
		minioClient, err = ConnectMinio()
		fatalIfErr(err)
	}

	repos := repository.NewRepositories(dbSQL, dbMongo, dbRedis)
	services := services.NewAppService(dbRedis, minioClient, repos)
	handlers := handler.NewAppHandlers(repos, services)

	// In-memory security registry initialization & load
	securityRegistry := security.NewSecurityRegistry(repos.RolesRelations)
	if err := securityRegistry.Reload(context.Background()); err != nil {
		log.Printf("⚠️ Warning: failed to load initial security registry: %v", err)
	}

	app := &application{
		port:     config.Get().App.Port,
		sql:      dbSQL,
		mongo:    dbMongo,
		redis:    dbRedis,
		minio:    minioClient,
		services: services,
		repos:    repos,
		handlers: handlers,
		security: securityRegistry,
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
		return nil, errors.New("MONGODB_URI environment variable not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(opts)
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	dbName := os.Getenv("MONGO_INITDB_DATABASE")
	if dbName == "" {
		dbName = "ticketdb"
	}

	log.Printf("✅ Connected to MongoDB: %s", dbName)
	return client.Database(dbName), nil
}

// ConnectRedis initializes the Redis client.
func ConnectRedis() (*redis.Client, error) {
	host := env.GetEnvString("REDIS_HOST", "localhost")
	port := env.GetEnvString("REDIS_PORT", "6379")
	password := os.Getenv("REDIS_PASSWORD")
	dbStr := env.GetEnvString("REDIS_DB", "0")
	db, err := strconv.Atoi(dbStr)
	if err != nil {
		db = 0
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: password,
		DB:       db,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	log.Printf("✅ Connected to Redis: %s:%s DB: %d", host, port, db)
	return rdb, nil
}

// ConnectMinio initializes the MinIO client and ensures the default bucket exists.
func ConnectMinio() (*minio.Client, error) {
	endpoint := env.GetEnvString("MINIO_ENDPOINT", "localhost:9000")
	accessKey := os.Getenv("ACCESS_KEY_MINIO")
	if accessKey == "" {
		accessKey = os.Getenv("MINIO_ROOT_USER")
	}
	secretKey := os.Getenv("SECRET_KEY_MINIO")
	if secretKey == "" {
		secretKey = os.Getenv("MINIO_ROOT_PASSWORD")
	}
	useSSL := config.Get().Minio.UseSSL
	bucketName := config.Get().Minio.Bucket

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	exists, err := minioClient.BucketExists(ctx, bucketName)
	if err != nil {
		return nil, err
	}

	if !exists {
		err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return nil, err
		}
		log.Printf("✅ Bucket %s created successfully", bucketName)
	} else {
		log.Printf("✅ Bucket %s already exists", bucketName)
	}

	log.Printf("✅ Connected to MinIO: %s", endpoint)
	return minioClient, nil
}
