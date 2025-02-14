package auth

import (
	"github.com/go-chi/render"
	"github.com/jrogozen/wargroovy/internal/config"
	"github.com/jrogozen/wargroovy/schema"
	u "github.com/jrogozen/wargroovy/utils"
	// log "github.com/sirupsen/logrus"
	"net/http"
)

func LoginAUser(configuration *config.Config) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := &schema.User{}

		// decode request body into user struct
		err := render.DecodeJSON(r.Body, user)

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			u.Respond(w, r, u.Message(false, "Invalid request"))
			return
		}

		rsp, status := Login(configuration, user.Email, user.Password)

		if rsp["status"].(bool) {
			token := rsp["user"].(*schema.SecureUserView).Token

			u.AttachAuthCookie(token, w)
		}

		w.WriteHeader(status)

		u.Respond(w, r, rsp)
		return
	})
}
