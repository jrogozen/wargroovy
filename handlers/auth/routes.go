package auth

import (
	"github.com/go-chi/chi"
	"github.com/jrogozen/wargroovy/internal/config"
)

func Routes(configuration *config.Config) *chi.Mux {
	router := chi.NewRouter()

	router.Post("/login", LoginAUser(configuration))

	return router
}
