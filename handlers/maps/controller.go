package maps

import (
	// "github.com/go-chi/chi"
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

// func GetAMap(configuration *config.Config) http.HandlerFunc {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		mapID := chi.URLParam(r, "mapId")
// 		m := FindMap(configuration, mapID)

// 		IncrementMapView(configuration, *m)

// 		if m == nil {
// 			u.Respond(w, r, u.Message(false, "Could not find map"))
// 		} else {
// 			response := u.Message(true, "Map found")
// 			response["map"] = m

// 			u.Respond(w, r, response)
// 		}
// 	})
// }

// func GetMapList(configuration *config.Config) http.HandlerFunc {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		sortOptions := GetSortOptions(r)
// 		maps := FindMapList(configuration, sortOptions)

// 		response := u.Message(true, "Maps found")
// 		response["maps"] = maps

// 		u.Respond(w, r, response)
// 	})
// }

// func EditAMap(configuration *config.Config) http.HandlerFunc {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		mapID := chi.URLParam(r, "mapId")
// 		mapUpdate := &schema.BaseMap{}

// 		// requires jwt-auth middleware to be used in part of the router stack
// 		_, claims, err := jwtauth.FromContext(r.Context())

// 		if err != nil {
// 			u.Respond(w, r, u.Message(false, "Error authorizing user"))
// 			return
// 		}

// 		err = render.DecodeJSON(r.Body, mapUpdate)

// 		if err != nil {
// 			u.Respond(w, r, u.Message(false, "Error updating map"))
// 			return
// 		}

// 		originalMap := FindMap(configuration, mapID)

// 		if originalMap == nil {
// 			u.Respond(w, r, u.Message(false, "Could not find map"))
// 			return
// 		}

// 		resp := UpdateMap(configuration, claims, originalMap, mapUpdate)

// 		u.Respond(w, r, resp)
// 		return
// 	})
// }
