package auth

import (
	"github.com/go-chi/render"
	"github.com/jrogozen/wargroovy/internal/config"
	"github.com/jrogozen/wargroovy/schema"
	u "github.com/jrogozen/wargroovy/utils"
	"net/http"
)

func LoginAUser(configuration *config.Config) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := &schema.User{}

		// decode request body into user struct
		err := render.DecodeJSON(r.Body, user)

		if err != nil {
			u.Respond(w, r, u.Message(false, "Invalid request"))
			return
		}

		rsp := Login(configuration, user.Email, user.Password)

		u.Respond(w, r, rsp)
		return
	})
}
