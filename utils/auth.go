package utils

import (
	"github.com/jrogozen/wargroovy/internal/config"
	"github.com/jrogozen/wargroovy/schema"
	log "github.com/sirupsen/logrus"
)

func AttachToken(configuration *config.Config, user *schema.User) *schema.User {
	admin := false
	_, tokenString, _ := configuration.TokenAuth.Encode(&schema.TokenClaims{UserID: user.ID, Admin: admin})

	log.WithFields(log.Fields{
		"token":  tokenString,
		"userId": user.ID,
		"admin":  admin,
	}).Trace("Generated token for user")

	user.Token = tokenString

	return user
}

func IsUserAuthorized(attemptedUserID uint, claims map[string]interface{}) (map[string]interface{}, bool) {
	actualUserID := claims["UserID"].(float64)

	if uint(actualUserID) != attemptedUserID {
		return Message(false, "Mismatched token userId and requestId"), false
	}

	return Message(true, "Authorized"), true
}
