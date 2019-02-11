package utils

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/jwtauth"
	"github.com/jrogozen/wargroovy/schema"
	"os"
)

func AttachToken(user *schema.User) *schema.User {
	claim := &schema.TokenClaims{UserID: user.ID}

	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), claim)
	tokenString, _ := token.SignedString([]byte(os.Getenv("jwt_secret")))
	user.Token = tokenString

	return user
}

/*InitJWT this is pretty lame.
it's because the jwtauth package requires passing this struct to the verifier middleware
currently we're attaching the token to login/create api endpoints manually
need to either not use the jwtauth provided middleware, or switch the AttachToken code
to however jwtauth wants us to make the claim

with doing it two separate ways, they should still resolve the same
as long as they are constructed using the same secret and signing method
*/
func InitJWT() *jwtauth.JWTAuth {
	return jwtauth.New("HS256", []byte(os.Getenv("jwt_secret")), nil)
}

func IsUserAuthorized(attemptedUserID uint, claims map[string]interface{}) (map[string]interface{}, bool) {
	actualUserID := claims["UserID"]

	if actualUserID != attemptedUserID {
		return Message(false, "Mismatched token userId and requestId"), false
	}

	return Message(true, "Authorized"), true
}
