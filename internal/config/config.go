package config

import (
	"fmt"
	"github.com/go-chi/jwtauth"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	"github.com/jrogozen/wargroovy/schema"
	log "github.com/sirupsen/logrus"
	"os"
)

type Constants struct {
	PORT string
}

type Config struct {
	Constants
	Database  *gorm.DB
	TokenAuth *jwtauth.JWTAuth
}

func Migrate(config *Config) {
	log.Info("Migrating database")

	config.Database.AutoMigrate(&schema.User{})
	config.Database.AutoMigrate(&schema.Campaign{})
	config.Database.AutoMigrate(&schema.Map{})
}

func DB() *gorm.DB {
	var (
		connectionString = os.Getenv("POSTGRES_CONNECTION")
	)

	log.Info(connectionString)

	db, err := gorm.Open("postgres", connectionString)

	if err != nil {
		panic(fmt.Sprintf("DB: %v", err))
	}

	return db
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

	db := DB()
	config.Database = db

	// constants setup
	port := os.Getenv("PORT")

	constants := Constants{PORT: port}

	// attach to config struct
	config.Constants = constants

	config.TokenAuth = InitJWT()

	log.Info("migrating soon")
	Migrate(&config)

	return &config, err
}
