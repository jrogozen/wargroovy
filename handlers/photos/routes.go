package photos

import (
	"github.com/go-chi/chi"
	// "github.com/go-chi/jwtauth"
	"github.com/jrogozen/wargroovy/internal/config"
	// u "github.com/jrogozen/wargroovy/utils"
)

func Routes(configuration *config.Config) *chi.Mux {
	router := chi.NewRouter()

	router.Group(func(router chi.Router) {
		router.Post("/", UploadPhotos(configuration))
	})

	return router
}
