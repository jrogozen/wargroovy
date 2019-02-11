package user

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
	"github.com/jinzhu/gorm"
	"github.com/jrogozen/wargroovy/internal/config"
	"github.com/jrogozen/wargroovy/schema"
	u "github.com/jrogozen/wargroovy/utils"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
)

func Routes(configuration *config.Config) *chi.Mux {
	router := chi.NewRouter()

	router.Group(func(router chi.Router) {
		/* looks for tokens in this order:
		'jwt' URI query parameter
		'Authorization: BEARER T' request header
		'jwt' Cookie value
		*/
		router.Use(jwtauth.Verifier(configuration.TokenAuth))
		router.Use(jwtauth.Authenticator)

	})

	router.Group(func(router chi.Router) {
		router.Post("/", CreateAUser(configuration))
		router.Get("/{userId}", GetAUser(configuration))
	})

	return router
}

/*Validate validates user fields for user creation
=> { "status": true, "message": "ok" }, true
*/
func Validate(configuration *config.Config, user *schema.User) (map[string]interface{}, bool) {
	if !strings.Contains(user.Email, "@") {
		return u.Message(false, "Email address not valid"), false
	}

	if len(user.Password) < 6 {
		return u.Message(false, "Invalid password"), false
	}

	temp := &schema.User{}

	err := configuration.Database.Table("users").Where("email = ?", user.Email).First(temp).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return u.Message(false, "Connection error"), false
	}
	if temp.Email != "" {
		return u.Message(false, "Email address already in use"), false
	}

	return u.Message(true, "Valid"), true
}

func Create(configuration *config.Config, user *schema.User) map[string]interface{} {
	if resp, ok := Validate(configuration, user); !ok {
		return resp
	}

	// hash password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)

	fmt.Println(*user)
	// add to db
	configuration.Database.Create(user)

	if user.ID <= 0 {
		return u.Message(false, "Failed to create user")
	}

	// jwt
	u.AttachToken(user)

	response := u.Message(true, "User created")
	response["user"] = user

	return response
}

func FindUser(configuration *config.Config, id string) *schema.UserWithOutPassword {
	user := &schema.UserWithOutPassword{}

	configuration.Database.Preload("Campaigns").Table("users").Where("id = ?", id).First(user)

	if user.Email == "" {
		return nil
	}

	return user
}

func CreateAUser(configuration *config.Config) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := &schema.User{}

		// decode request body into user struct
		err := render.DecodeJSON(r.Body, user)

		if err != nil {
			u.Respond(w, r, u.Message(false, "Invalid request"))
		} else {
			resp := Create(configuration, user)
			u.Respond(w, r, resp)
		}
	})
}

func GetAUser(configuration *config.Config) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "userId")
		user := FindUser(configuration, userID)

		if user == nil {
			u.Respond(w, r, u.Message(false, "Could not find user"))
		} else {
			response := u.Message(true, "User found")
			response["user"] = user

			u.Respond(w, r, response)
		}
	})
}
