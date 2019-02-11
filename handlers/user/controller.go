package user

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/jrogozen/wargroovy/internal/config"
	"github.com/jrogozen/wargroovy/schema"
	u "github.com/jrogozen/wargroovy/utils"
	"net/http"
)

func CreateAUser(configuration *config.Config) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := &schema.User{}

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
		userID := chi.URLParam(r, "userId")
		user := FindUser(configuration, userID)

		if user == nil {
			u.Respond(w, r, u.Message(false, "Could not find user"))
			return
		}

		response := u.Message(true, "User found")
		response["user"] = user

		u.Respond(w, r, response)
		return
	})
}
