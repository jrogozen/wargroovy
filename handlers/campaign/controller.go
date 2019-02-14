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

		// requires jwt-auth middleware to be used in part of the router stack
		_, claims, err := jwtauth.FromContext(r.Context())

		if err != nil {
			u.Respond(w, r, u.Message(false, "Error authorizing user"))
			return
		}

		err = render.DecodeJSON(r.Body, campaign)

		if err != nil {
			u.Respond(w, r, u.Message(false, "Invalid request"))
			return
		}

		resp := Create(configuration, claims, campaign)

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

		resp := UpdateCampaign(configuration, claims, originalCampaign, campaignUpdate)

		u.Respond(w, r, resp)
		return
	})
}

func CreateAMap(configuration *config.Config) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		campaignIDString := chi.URLParam(r, "campaignId")

		// need claims to verify that a user can create this map
		_, claims, err := jwtauth.FromContext(r.Context())

		if err != nil {
			u.Respond(w, r, u.Message(false, "Error authorizing user"))
			return
		}

		m := &schema.Map{}

		// find campaign
		campaign := FindCampaign(configuration, campaignIDString)

		if campaign == nil {
			u.Respond(w, r, u.Message(false, "Could not find campaign"))
			return
		}

		// marshal request body into map schema
		err = render.DecodeJSON(r.Body, m)

		if err != nil {
			u.Respond(w, r, u.Message(false, "Invalid request"))
			return
		}

		campaignID, _ := strconv.Atoi(campaignIDString)

		// set the campaign id for map schema
		m.CampaignID = uint(campaignID)

		// create and validate map
		resp := CreateMap(configuration, claims, m, campaign)

		u.Respond(w, r, resp)
		return
	})
}

func EditAMap(configuration *config.Config) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		campaignIDString := chi.URLParam(r, "campaignId")
		mapIDString := chi.URLParam(r, "mapId")

		// requires jwt-auth middleware to be used in part of the router stack
		_, claims, err := jwtauth.FromContext(r.Context())

		if err != nil {
			u.Respond(w, r, u.Message(false, "Error authorizing user"))
			return
		}

		// find campaign
		campaign := FindCampaign(configuration, campaignIDString)

		if campaign == nil {
			u.Respond(w, r, u.Message(false, "Could not find campaign"))
			return
		}

		m := &schema.BaseMap{}

		// marshal json body into map struct
		err = render.DecodeJSON(r.Body, m)

		if err != nil {
			u.Respond(w, r, u.Message(false, "Error updating map"))
			return
		}

		response := UpdateMap(configuration, claims, campaign, mapIDString, m)

		u.Respond(w, r, response)
		return
	})
}
