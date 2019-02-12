package auth

import (
	"github.com/jinzhu/gorm"
	"github.com/jrogozen/wargroovy/internal/config"
	"github.com/jrogozen/wargroovy/schema"
	u "github.com/jrogozen/wargroovy/utils"
	"golang.org/x/crypto/bcrypt"
)

func Login(configuration *config.Config, email string, password string) map[string]interface{} {
	user := &schema.User{}

	err := configuration.Database.
		Preload("Campaigns").
		Table("users").
		Where("email = ?", email).
		First(user).
		Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.Message(false, "Email address not found")
		}
		return u.Message(false, "Connection error")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return u.Message(false, "Invalid credentials")
	}

	// delete sensitive information
	user.Password = ""

	// jwt
	u.AttachToken(configuration, user)

	resp := u.Message(true, "Logged in")
	resp["user"] = user

	return resp
}
