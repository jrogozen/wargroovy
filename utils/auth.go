package utils

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/jwtauth"
	"github.com/jrogozen/wargroovy/internal/config"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func GetToken(configuration *config.Config, userID int64) string {
	admin := false

	_, tokenString, _ := configuration.TokenAuth.Encode(jwt.MapClaims{"UserID": userID, "Admin": admin})

	log.WithFields(log.Fields{
		"token":  tokenString,
		"userId": userID,
		"admin":  admin,
	}).Info("Generated token for user")

	return tokenString
}

func AttachAuthCookie(token string, w http.ResponseWriter) {
	responseCookie := &http.Cookie{
		Name:     "jwt",
		Value:    token,
		MaxAge:   0,
		HttpOnly: true,
	}

	http.SetCookie(w, responseCookie)
}

func IsUserAuthorized(attemptedUserID int64, claims map[string]interface{}) (map[string]interface{}, bool) {
	actualUserID := int64(claims["UserID"].(float64))

	if actualUserID != attemptedUserID {
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
