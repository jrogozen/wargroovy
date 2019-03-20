package utils

import (
	"errors"
	"github.com/go-chi/jwtauth"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func AttachAuthCookie(token string, w http.ResponseWriter) {
	expire := time.Now().AddDate(0, 0, 7*52) // 1 year

	responseCookie := &http.Cookie{
		Name:     "jwt",
		Value:    token,
		MaxAge:   0,
		Expires:  expire,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		Domain:   ".wargroovy.com",
	}

	http.SetCookie(w, responseCookie)
}

func GetUserIdFromClaims(claims map[string]interface{}) (int64, error) {
	if claims["UserID"] == nil {
		return 0, errors.New("No valid UserID found in claims")
	}

	actualUserID := int64(claims["UserID"].(float64))

	return actualUserID, nil
}

func IsUserAuthorized(attemptedUserID int64, claims map[string]interface{}) (map[string]interface{}, bool) {
	actualUserID := int64(claims["UserID"].(float64))

	if actualUserID != attemptedUserID {
		log.WithFields(log.Fields{
			"mapUserID":    attemptedUserID,
			"claimsUserID": actualUserID,
			"claims":       claims,
		}).Info("could not match userID and claimsUserID")

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
