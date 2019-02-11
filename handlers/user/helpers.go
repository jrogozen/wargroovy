package user

import (
	"github.com/jinzhu/gorm"
	"github.com/jrogozen/wargroovy/internal/config"
	"github.com/jrogozen/wargroovy/schema"
	u "github.com/jrogozen/wargroovy/utils"
	"strings"
)

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
