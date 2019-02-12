package user

import (
	"github.com/jrogozen/wargroovy/internal/config"
	"github.com/jrogozen/wargroovy/schema"
	u "github.com/jrogozen/wargroovy/utils"
	"golang.org/x/crypto/bcrypt"
)

func Create(configuration *config.Config, user *schema.User) map[string]interface{} {
	if resp, ok := Validate(configuration, user); !ok {
		return resp
	}

	// hash password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)

	// add to db
	configuration.Database.Create(user)

	if user.ID <= 0 {
		return u.Message(false, "Failed to create user")
	}

	// jwt
	u.AttachToken(configuration, user)

	user.Password = ""

	response := u.Message(true, "User created")
	response["user"] = user

	return response
}

func FindUser(configuration *config.Config, id string) *schema.UserWithOutPassword {
	user := &schema.UserWithOutPassword{}

	configuration.Database.
		Preload("Campaigns").
		Table("users").
		Where("id = ?", id).
		First(user)

	if user.Email == "" {
		return nil
	}

	return user
}
