package campaign

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
	"github.com/jrogozen/wargroovy/internal/config"
	"github.com/jrogozen/wargroovy/schema"
	u "github.com/jrogozen/wargroovy/utils"
	"net/http"
	"strconv"
)

func CreateACampaign(configuration *config.Config) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		campaign := &schema.Campaign{}

		err := render.DecodeJSON(r.Body, campaign)

		if err != nil {
			u.Respond(w, r, u.Message(false, "Invalid request"))
			return
		}

		resp := Create(configuration, campaign)
		u.Respond(w, r, resp)
		return
	})
}

func GetACampaign(configuration *config.Config) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		campaignID := chi.URLParam(r, "campaignId")
		campaign := FindCampaign(configuration, campaignID)

		IncrementCampaignView(configuration, *campaign)

		if campaign == nil {
			u.Respond(w, r, u.Message(false, "Could not find campaign"))
		} else {
			response := u.Message(true, "Campaign found")
			response["campaign"] = campaign

			u.Respond(w, r, response)
		}
	})
}

func GetCampaignsList(configuration *config.Config) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sortOptions := GetSortOptions(r)
		campaigns := FindCampaignList(configuration, sortOptions)

		response := u.Message(true, "Campaigns found")
		response["campaigns"] = campaigns

		u.Respond(w, r, response)
	})
}

func EditACampaign(configuration *config.Config) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		campaignID := chi.URLParam(r, "campaignId")
		campaignUpdate := &schema.BaseCampaign{}

		// requires jwt-auth middleware to be used in part of the router stack
		_, claims, err := jwtauth.FromContext(r.Context())

		if err != nil {
			u.Respond(w, r, u.Message(false, "Error authorizing user"))
			return
		}

		err = render.DecodeJSON(r.Body, campaignUpdate)

		if err != nil {
			u.Respond(w, r, u.Message(false, "Error updating campaign"))
			return
		}

		originalCampaign := FindCampaign(configuration, campaignID)

		if originalCampaign == nil {
			u.Respond(w, r, u.Message(false, "Could not find campaign"))
			return
		}

		updatedCampaign := UpdateCampaign(configuration, claims, originalCampaign, campaignUpdate)

		if updatedCampaign == nil {
			response := u.Message(false, "Not authorized to edit this campaign")
			u.Respond(w, r, response)
			return
		}
		response := u.Message(true, "Campaign updated")
		response["campaign"] = updatedCampaign

		u.Respond(w, r, response)
		return
	})
}

func CreateAMap(configuration *config.Config) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		campaignID, _ := strconv.Atoi(chi.URLParam(r, "campaignId"))
		m := &schema.Map{}

		err := render.DecodeJSON(r.Body, m)

		m.CampaignID = uint(campaignID)

		if err != nil {
			u.Respond(w, r, u.Message(false, "Invalid request"))
			return
		}

		resp := CreateMap(configuration, m)
		u.Respond(w, r, resp)
		return
	})
}

func EditAMap(configuration *config.Config) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		campaignID := chi.URLParam(r, "campaignId")
		mapID := chi.URLParam(r, "mapId")
		m := &schema.BaseMap{}

		err := render.DecodeJSON(r.Body, m)

		if err != nil {
			u.Respond(w, r, u.Message(false, "Error updating map"))
			return
		}

		updatedMap := UpdateMap(configuration, campaignID, mapID, m)

		if updatedMap == nil {
			response := u.Message(false, "Could not find map")
			u.Respond(w, r, response)
			return
		}

		response := u.Message(true, "Map updated")
		response["map"] = updatedMap

		u.Respond(w, r, response)
		return
	})
}
