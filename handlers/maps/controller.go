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
			w.WriteHeader(http.StatusForbidden)
			u.Respond(w, r, u.Message(false, "Error authorizing user"))
			return
		}

		err = render.DecodeJSON(r.Body, m)

		if err != nil {
			log.Info(err)

			w.WriteHeader(http.StatusBadRequest)
			u.Respond(w, r, u.Message(false, "Invalid request"))
			return
		}

		resp, status := Create(configuration, claims, m)

		w.WriteHeader(status)
		u.Respond(w, r, resp)
		return
	})
}

func GetAMap(configuration *config.Config) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mapID := chi.URLParam(r, "mapId")
		response, status := FindMap(configuration, mapID)

		w.WriteHeader(status)
		u.Respond(w, r, response)
		return
	})
}

func GetAMapBySlug(configuration *config.Config) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slug := chi.URLParam(r, "slug")

		response, status := FindMapBySlug(configuration, slug)

		w.WriteHeader(status)
		u.Respond(w, r, response)
		return
	})
}

func GetMapList(configuration *config.Config) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sortOptions := GetSortOptions(r)
		response, status := FindMapList(configuration, sortOptions)

		w.WriteHeader(status)
		u.Respond(w, r, response)
		return
	})
}

func EditAMap(configuration *config.Config) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mapIDString := chi.URLParam(r, "mapId")
		mapUpdate := &schema.Map{}

		// requires jwt-auth middleware to be used in part of the router stack
		_, claims, err := jwtauth.FromContext(r.Context())

		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			u.Respond(w, r, u.Message(false, "Error authorizing user"))
			return
		}

		err = render.DecodeJSON(r.Body, mapUpdate)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			u.Respond(w, r, u.Message(false, "Error updating map"))
			return
		}

		resp, status := UpdateMap(configuration, claims, mapIDString, mapUpdate)

		w.WriteHeader(status)
		u.Respond(w, r, resp)
		return
	})
}

func DeleteMapPhoto(configuration *config.Config) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mapIDString := chi.URLParam(r, "mapId")

		// requires jwt-auth middleware to be used in part of the router stack
		_, claims, err := jwtauth.FromContext(r.Context())

		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			u.Respond(w, r, u.Message(false, "Error authorizing user"))
			return
		}

		type deleteMapOptions struct {
			URL string `json:"url"`
		}

		options := &deleteMapOptions{}

		err = render.DecodeJSON(r.Body, options)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			u.Respond(w, r, u.Message(false, "Error deleting photo"))
			return
		}

		resp, status := DeletePhoto(configuration, claims, mapIDString, options.URL)

		w.WriteHeader(status)
		u.Respond(w, r, resp)
		return
	})
}

func DeleteAMap(configuration *config.Config) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mapIDString := chi.URLParam(r, "mapId")

		// requires jwt-auth middleware to be used in part of the router stack
		_, claims, err := jwtauth.FromContext(r.Context())

		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			u.Respond(w, r, u.Message(false, "Error authorizing user"))
			return
		}

		resp, status := Delete(configuration, claims, mapIDString)

		w.WriteHeader(status)
		u.Respond(w, r, resp)
		return
	})
}

func RateAMap(configuration *config.Config) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mapIDString := chi.URLParam(r, "mapId")

		// requires jwt-auth middleware to be used in part of the router stack
		_, claims, err := jwtauth.FromContext(r.Context())

		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			u.Respond(w, r, u.Message(false, "Error authorizing user"))
			return
		}

		type ratingBody struct {
			Rating int64 `json:"rating"`
		}

		decodedBody := &ratingBody{}

		err = render.DecodeJSON(r.Body, decodedBody)

		log.WithField("decodedRating", decodedBody.Rating).Info("Decoded rating")

		resp, status := Rate(configuration, claims, mapIDString, decodedBody.Rating)

		w.WriteHeader(status)
		u.Respond(w, r, resp)
		return
	})
}
