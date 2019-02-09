package config

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Constants struct {
	PORT   string
	DB_URI string
}

type Config struct {
	Constants
	Database *gorm.DB
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

	return &config, err
}
