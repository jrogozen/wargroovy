package maps

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"github.com/jrogozen/wargroovy/internal/config"
	u "github.com/jrogozen/wargroovy/utils"
)

func Routes(configuration *config.Config) *chi.Mux {
	router := chi.NewRouter()

	router.Group(func(router chi.Router) {
		router.Use(jwtauth.Verifier(configuration.TokenAuth))
		router.Use(u.Authenticator)

		router.Post("/", CreateAMap(configuration))
		router.Put("/{mapId}", EditAMap(configuration))
		router.Delete("/{mapId}/photo", DeleteMapPhoto(configuration))

	})

	router.Group(func(router chi.Router) {
		router.Get("/list", GetMapList(configuration))
		router.Get("/bySlug/{slug}", GetAMapBySlug(configuration))
		router.Get("/{mapId}", GetAMap(configuration))
	})

	return router
}
