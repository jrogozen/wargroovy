package config

import (
	"bitbucket.org/liamstask/goose/lib/goose"
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/jwtauth"
	"github.com/joho/godotenv"
	"github.com/jrogozen/wargroovy/db"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

type Constants struct {
	PORT              string
	StorageBucketName string
}

type Config struct {
	Constants
	DB            *db.PsqlDB
	TokenAuth     *jwtauth.JWTAuth
	StorageBucket *storage.BucketHandle
}

func Migrate(config *Config) {
	log.Info("Migrating database")

	p, _ := filepath.Abs("../db/migrations")

	migrateConf := &goose.DBConf{
		MigrationsDir: p,
		Env:           "development",
		Driver: goose.DBDriver{
			Name:    "postgres",
			OpenStr: os.Getenv("POSTGRES_CONNECTION"),
			Import:  "https://github.com/lib/pq",
			Dialect: &goose.PostgresDialect{},
		},
	}

	// Get the latest possible migration
	latest, err := goose.GetMostRecentDBVersion(migrateConf.MigrationsDir)

	if err != nil {
		log.Error(err)
		return
	}

	err = goose.RunMigrationsOnDb(migrateConf, migrateConf.MigrationsDir, latest, config.DB.Conn)

	if err != nil {
		log.Error(err)
		return
	}
}

func MakeDB() *db.PsqlDB {
	var (
		connectionString = os.Getenv("POSTGRES_CONNECTION")
	)

	psqlDB, err := db.NewPostgresDB(connectionString)

	if err != nil {
		panic(fmt.Errorf("DB: %v", err))
	}

	return psqlDB
}

func InitJWT() *jwtauth.JWTAuth {
	return jwtauth.New("HS256", []byte(os.Getenv("jwt_secret")), nil)
}

func (config *Config) GetToken(userID int64) string {
	admin := false

	_, tokenString, _ := config.TokenAuth.Encode(jwt.MapClaims{"UserID": userID, "Admin": admin})

	log.WithFields(log.Fields{
		"token":  tokenString,
		"userId": userID,
		"admin":  admin,
	}).Info("Generated token for user")

	return tokenString
}

func configureStorage(bucketID string) (*storage.BucketHandle, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)

	if err != nil {
		return nil, err
	}

	return client.Bucket(bucketID), nil
}

func New() (*Config, error) {
	config := Config{}

	// AppEngine instances won't be able to load this
	err := godotenv.Load("../.env")

	if err != nil {
		log.Warn("Error loading dotenv", err)
	}

	db := MakeDB()
	config.DB = db

	// constants setup
	port := os.Getenv("PORT")

	constants := Constants{PORT: port}

	// attach to config struct
	config.Constants = constants

	config.TokenAuth = InitJWT()

	Migrate(&config)

	StorageBucketName := os.Getenv("STORAGE_BUCKET_NAME")
	StorageBucket, err := configureStorage(StorageBucketName)

	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Could not connect to storage bucket!")
	}

	config.Constants.StorageBucketName = StorageBucketName
	config.StorageBucket = StorageBucket

	return &config, err
}
