package utils

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/jwtauth"
	"github.com/jrogozen/wargroovy/internal/config"
	"github.com/jrogozen/wargroovy/schema"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func AttachToken(configuration *config.Config, user *schema.User) *schema.User {
	admin := false

	_, tokenString, _ := configuration.TokenAuth.Encode(jwt.MapClaims{"UserID": user.ID, "Admin": admin})

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

func Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		token, _, err := jwtauth.FromContext(ctx)

		log.WithField("token", token).Trace("authenticating token")

		if token == nil || !token.Valid {
			w.WriteHeader(http.StatusForbidden)
			Respond(w, r, Message(false, "invalid token"))
			return
		}

		if err != nil {
			w.WriteHeader(http.StatusForbidden)

			// get string of error message
			Respond(w, r, Message(false, err.Error()))
			return
		}

		// Token is authenticated, pass it through
		next.ServeHTTP(w, r)
	})
}
