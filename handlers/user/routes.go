package user

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"github.com/jrogozen/wargroovy/internal/config"
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
