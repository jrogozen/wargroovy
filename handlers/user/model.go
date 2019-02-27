package user

import (
	"github.com/jrogozen/wargroovy/internal/config"
	"github.com/jrogozen/wargroovy/schema"
	u "github.com/jrogozen/wargroovy/utils"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
)

func Create(configuration *config.Config, user *schema.CreateUser) (map[string]interface{}, int) {
	if resp, ok := Validate(configuration, user); !ok {
		return resp, http.StatusBadRequest
	}

	// hash password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)

	insertedID, err := configuration.DB.AddUser(user)

	if err != nil || insertedID == 0 {
		return u.Message(false, err.Error()), http.StatusBadRequest
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

	return response, http.StatusOK
}

func FindUser(configuration *config.Config, userIdString string) (map[string]interface{}, int) {
	userID, _ := strconv.ParseInt(userIdString, 10, 64)
	user, err := configuration.DB.GetUser(userID)

	if err != nil {
		return u.Message(false, err.Error()), http.StatusBadRequest
	}

	// convert user into what we want to return with API
	safeUserView := &schema.UserView{
		ID:       user.ID,
		Email:    user.Email,
		Username: user.Username,
	}

	response := u.Message(true, "User found")
	response["user"] = safeUserView

	return response, http.StatusOK
}

func UpdateUser(configuration *config.Config, claims map[string]interface{}, userIDString string, user *schema.User) (map[string]interface{}, int) {
	userID, _ := strconv.ParseInt(userIDString, 10, 64)

	originalUser, err := configuration.DB.GetUser(userID)

	if err != nil {
		return u.Message(false, err.Error()), http.StatusBadRequest
	}

	log.WithFields(log.Fields{
		"user": originalUser,
	}).Info("UpdateUser model")

	if resp, ok := ValidateUpdate(configuration, claims, originalUser); !ok {
		return resp, http.StatusForbidden
	}

	originalUser.Merge(user)

	insertedID, err := configuration.DB.UpdateUser(originalUser)

	if err != nil {
		return u.Message(false, err.Error()), http.StatusBadRequest
	}

	//TODO missing created at, updated at
	updatedUser := &schema.UserView{
		ID:       insertedID,
		Email:    originalUser.Email,
		Username: originalUser.Username,
	}

	response := u.Message(true, "User updated")
	response["user"] = updatedUser

	return response, http.StatusOK
}
