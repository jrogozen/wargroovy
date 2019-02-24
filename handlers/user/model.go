package user

import (
	"github.com/jrogozen/wargroovy/internal/config"
	"github.com/jrogozen/wargroovy/schema"
	u "github.com/jrogozen/wargroovy/utils"
	// log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"strconv"
)

func Create(configuration *config.Config, user *schema.CreateUser) map[string]interface{} {
	if resp, ok := Validate(configuration, user); !ok {
		return resp
	}

	// hash password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)

	insertedID, err := configuration.DB.AddUser(user)

	if err != nil || insertedID == 0 {
		return u.Message(false, err.Error())
	}

	//TODO probably can move this into the db package
	createdUser := &schema.CreatedUser{
		ID:       insertedID,
		Email:    user.Email,
		Username: user.Username,
		Token:    u.GetToken(configuration, insertedID),
	}

	response := u.Message(true, "User created")
	response["user"] = createdUser

	return response
}

func FindUser(configuration *config.Config, userIdString string) map[string]interface{} {
	userID, _ := strconv.ParseInt(userIdString, 10, 64)
	user, err := configuration.DB.GetUser(userID)

	if err != nil {
		return u.Message(false, err.Error())
	}

	response := u.Message(true, "User found")
	response["user"] = user

	return response
}
