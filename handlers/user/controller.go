package user

import (
	"fmt"
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
			u.Respond(w, r, u.Message(false, "Invalid request"))
			return
		}

		resp := Create(configuration, user)
		u.Respond(w, r, resp)
		return
	})
}

func GetAUser(configuration *config.Config) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// requires jwt-auth middleware to be used in part of the router stack
		_, claims, err := jwtauth.FromContext(r.Context())

		if err != nil || claims["UserID"] == nil {
			u.Respond(w, r, u.Message(false, "Error authorizing user"))
			return
		}

		userIDString := fmt.Sprintf("%d", int(claims["UserID"].(float64)))

		log.WithFields(log.Fields{
			"userIDString": userIDString,
		}).Info("Found UserID from jwt")

		resp := FindUser(configuration, userIDString)

		u.Respond(w, r, resp)
		return
	})
}
