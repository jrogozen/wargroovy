package config

import (
	"bitbucket.org/liamstask/goose/lib/goose"
	"fmt"
	"github.com/go-chi/jwtauth"
	"github.com/joho/godotenv"
	"github.com/jrogozen/wargroovy/db"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"

	"os"
	"path/filepath"
)

type Constants struct {
	PORT string
}

type Config struct {
	Constants
	DB        *db.PsqlDB
	TokenAuth *jwtauth.JWTAuth
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

	return &config, err
}
