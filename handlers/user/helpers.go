package user

import (
	"github.com/jrogozen/wargroovy/internal/config"
	"github.com/jrogozen/wargroovy/schema"
	u "github.com/jrogozen/wargroovy/utils"
	"strings"
)

/*Validate validates user fields for user creation
=> { "status": true, "message": "ok" }, true
*/
func Validate(configuration *config.Config, user *schema.CreateUser) (map[string]interface{}, bool) {
	if !strings.Contains(user.Email, "@") {
		return u.Message(false, "Email address not valid"), false
	}

	if len(user.Password) < 6 {
		return u.Message(false, "Invalid password"), false
	}

	return u.Message(true, "Valid"), true
}
