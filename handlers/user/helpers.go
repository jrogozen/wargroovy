package user

import (
	"github.com/jrogozen/wargroovy/internal/config"
	"github.com/jrogozen/wargroovy/schema"
	u "github.com/jrogozen/wargroovy/utils"
	log "github.com/sirupsen/logrus"
	"strings"
)

/*Validate validates user fields for user creation
=> { "status": true, "message": "ok" }, true
*/
func Validate(configuration *config.Config, user *schema.CreateUser) (map[string]interface{}, bool) {
	hasEmail := user.Email != ""
	hasUsername := user.Username != ""
	creatingRealAccount := hasEmail || hasUsername

	if creatingRealAccount && !strings.Contains(user.Email, "@") {
		return u.Message(false, "Email address not valid"), false
	}

	if creatingRealAccount && len(user.Password) < 6 {
		return u.Message(false, "Invalid password"), false
	}

	if creatingRealAccount && len(user.Username) < 1 {
		return u.Message(false, "Must create a username"), false
	}

	return u.Message(true, "Valid"), true
}

func ValidateUpdate(configuration *config.Config, claims map[string]interface{}, user *schema.User) (map[string]interface{}, bool) {
	log.WithFields(log.Fields{
		"userID": user.ID,
	}).Info("validate update")

	if _, ok := u.IsUserAuthorized(user.ID, claims); !ok {
		return u.Message(false, "Can only edit user with valid token"), false
	}

	return u.Message(true, "Valid"), true
}
