package maps

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
	"github.com/jrogozen/wargroovy/internal/config"
	"github.com/jrogozen/wargroovy/schema"
	u "github.com/jrogozen/wargroovy/utils"
	log "github.com/sirupsen/logrus"
	"net/http"
	// "strconv"
)

func CreateAMap(configuration *config.Config) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m := &schema.Map{}

		// requires jwt-auth middleware to be used in part of the router stack
		_, claims, err := jwtauth.FromContext(r.Context())

		if err != nil {
			u.Respond(w, r, u.Message(false, "Error authorizing user"))
			return
		}

		err = render.DecodeJSON(r.Body, m)

		if err != nil {
			log.Info(err)

			u.Respond(w, r, u.Message(false, "Invalid request"))
			return
		}

		resp := Create(configuration, claims, m)

		u.Respond(w, r, resp)
		return
	})
}

func GetAMap(configuration *config.Config) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mapID := chi.URLParam(r, "mapId")
		response := FindMap(configuration, mapID)

		u.Respond(w, r, response)

		return
	})
}

func GetMapList(configuration *config.Config) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sortOptions := GetSortOptions(r)
		response := FindMapList(configuration, sortOptions)

		u.Respond(w, r, response)

		return
	})
}

func EditAMap(configuration *config.Config) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mapIdString := chi.URLParam(r, "mapId")
		mapUpdate := &schema.Map{}

		// requires jwt-auth middleware to be used in part of the router stack
		_, claims, err := jwtauth.FromContext(r.Context())

		if err != nil {
			u.Respond(w, r, u.Message(false, "Error authorizing user"))
			return
		}

		err = render.DecodeJSON(r.Body, mapUpdate)

		if err != nil {
			u.Respond(w, r, u.Message(false, "Error updating map"))
			return
		}

		resp := UpdateMap(configuration, claims, mapIdString, mapUpdate)

		u.Respond(w, r, resp)
		return
	})
}
