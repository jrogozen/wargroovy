package config

import (
	"fmt"
	"github.com/go-chi/jwtauth"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	"github.com/jrogozen/wargroovy/schema"
	u "github.com/jrogozen/wargroovy/utils"
	log "github.com/sirupsen/logrus"
	"os"
)

type Constants struct {
	PORT   string
	DB_URI string
}

type Config struct {
	Constants
	Database  *gorm.DB
	TokenAuth *jwtauth.JWTAuth
}

func Migrate(config *Config) {
	log.Info("Migrating database.")

	config.Database.AutoMigrate(&schema.User{})
	config.Database.AutoMigrate(&schema.Campaign{})
	config.Database.AutoMigrate(&schema.Map{})
}

func NewDB(dbUri string) (*gorm.DB, error) {
	db, err := gorm.Open("postgres", dbUri)

	if err != nil {
		return nil, err
	}

	return db, nil
}

func New() (*Config, error) {
	config := Config{}

	err := godotenv.Load()

	if err != nil {
		log.Panicln("Error loading dotenv", err)

		return &config, err
	}

	// get constants from dotEnv
	port := os.Getenv("port")
	dbName := os.Getenv("db_name")
	dbHost := os.Getenv("db_host")

	dbUri := fmt.Sprintf("host=%s dbname=%s sslmode=disable", dbHost, dbName)

	fmt.Println(dbUri)

	constants := Constants{PORT: port, DB_URI: dbUri}

	// attach to config struct
	config.Constants = constants

	db, err := NewDB(dbUri)

	if err != nil {
		log.Panicln("Error connecting to db", err)

		return &config, err
	}

	config.Database = db

	config.TokenAuth = u.InitJWT()

	Migrate(&config)

	return &config, err
}
