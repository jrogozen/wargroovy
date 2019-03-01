package auth

import (
	"github.com/jrogozen/wargroovy/internal/config"
	"github.com/jrogozen/wargroovy/schema"
	u "github.com/jrogozen/wargroovy/utils"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func Login(configuration *config.Config, email string, password string) (map[string]interface{}, int) {
	user, err := configuration.DB.GetUserByLogin(email)

	if err != nil {
		return u.Message(false, err.Error()), http.StatusForbidden
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return u.Message(false, "Invalid credentials"), http.StatusForbidden
	}

	secureUser := &schema.SecureUserView{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
		Username:  user.Username,
		Token:     configuration.GetToken(user.ID),
	}

	resp := u.Message(true, "Logged in")
	resp["user"] = secureUser

	return resp, http.StatusOK
}
