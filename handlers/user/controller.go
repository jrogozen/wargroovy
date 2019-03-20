package user

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
	"github.com/jrogozen/wargroovy/internal/config"
	"github.com/jrogozen/wargroovy/schema"
	u "github.com/jrogozen/wargroovy/utils"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func CreateAUser(configuration *config.Config) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := &schema.CreateUser{}

		// decode request body into user struct
		err := render.DecodeJSON(r.Body, user)

		if err != nil {
			log.Error(err)
			u.Respond(w, r, u.Message(false, "Invalid request"))
			return
		}

		resp, status := Create(configuration, user)

		if resp["status"].(bool) {
			token := resp["user"].(*schema.CreatedUser).Token

			u.AttachAuthCookie(token, w)
		}

		w.WriteHeader(status)
		u.Respond(w, r, resp)
		return
	})
}

func GetAUser(configuration *config.Config) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// requires jwt-auth middleware to be used in part of the router stack
		_, claims, err := jwtauth.FromContext(r.Context())

		if err != nil || claims["UserID"] == nil {
			w.WriteHeader(http.StatusUnauthorized)
			u.Respond(w, r, u.Message(false, "Error authorizing user"))
			return
		}

		log.WithField("claims", claims).Info("Found claims")

		userIDString := fmt.Sprintf("%d", int(claims["UserID"].(float64)))

		log.WithFields(log.Fields{
			"userIDString": userIDString,
		}).Info("Found UserID from claims")

		resp, status := FindUser(configuration, userIDString)

		w.WriteHeader(status)
		u.Respond(w, r, resp)
		return
	})
}

func EditAUser(configuration *config.Config) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userIDString := chi.URLParam(r, "userId")
		userUpdate := &schema.User{}

		// requires jwt-auth middleware to be used in part of the router stack
		_, claims, err := jwtauth.FromContext(r.Context())

		if err != nil || claims["UserID"] == nil {
			w.WriteHeader(http.StatusUnauthorized)
			u.Respond(w, r, u.Message(false, "Error authorizing user"))
			return
		}

		err = render.DecodeJSON(r.Body, userUpdate)

		if err != nil {
			u.Respond(w, r, u.Message(false, "Error updating user"))
			return
		}

		log.WithFields(log.Fields{
			"userIdString": userIDString,
		}).Info("Edit A User controller")

		resp, status := UpdateUser(configuration, claims, userIDString, userUpdate)

		w.WriteHeader(status)
		u.Respond(w, r, resp)
		return
	})
}
