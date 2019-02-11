package auth

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/jinzhu/gorm"
	"github.com/jrogozen/wargroovy/internal/config"
	"github.com/jrogozen/wargroovy/schema"
	u "github.com/jrogozen/wargroovy/utils"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func Routes(configuration *config.Config) *chi.Mux {
	router := chi.NewRouter()

	router.Post("/login", LoginAUser(configuration))

	return router
}

func Login(configuration *config.Config, email string, password string) map[string]interface{} {
	user := &schema.User{}

	err := configuration.Database.Preload("Campaigns").Table("users").Where("email = ?", email).First(user).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.Message(false, "Email address not found")
		}
		return u.Message(false, "Connection error")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return u.Message(false, "Invalid credentials")
	}

	// delete sensitive information
	user.Password = ""

	// jwt
	u.AttachToken(user)

	resp := u.Message(true, "Logged in")
	resp["user"] = user

	return resp
}

func LoginAUser(configuration *config.Config) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := &schema.User{}

		// decode request body into user struct
		err := render.DecodeJSON(r.Body, user)

		if err != nil {
			u.Respond(w, r, u.Message(false, "Invalid request"))
		} else {
			rsp := Login(configuration, user.Email, user.Password)

			u.Respond(w, r, rsp)
		}
	})
}
