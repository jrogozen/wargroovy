package campaign

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

		router.Post("/{campaignId}/map", CreateAMap(configuration))
		router.Put("/{campaignId}/map/{mapId}", EditAMap(configuration))
		router.Post("/", CreateACampaign(configuration))

		router.Put("/{campaignId}", EditACampaign(configuration))

	})

	router.Group(func(router chi.Router) {
		router.Get("/{campaignId}", GetACampaign(configuration))
		router.Get("/list", GetCampaignsList(configuration))

	})

	return router
}