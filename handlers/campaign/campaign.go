package campaign

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/jrogozen/wargroovy/internal/config"
	"github.com/jrogozen/wargroovy/schema"
	u "github.com/jrogozen/wargroovy/utils"
	"net/http"
	"strconv"
)

func Routes(configuration *config.Config) *chi.Mux {
	router := chi.NewRouter()

	// maps
	router.Post("/{campaignId}/map", CreateAMap(configuration))

	// campaigns
	router.Post("/", CreateACampaign(configuration))
	router.Get("/{campaignId}", GetACampaign(configuration))

	return router
}

//Validate validates campaign fields for campaign creation
func Validate(configuration *config.Config, campaign *schema.Campaign) (map[string]interface{}, bool) {
	if campaign.UserID <= 0 {
		return u.Message(false, "Campaigns need to be owned by a user"), false
	}

	return u.Message(true, "Valid"), true
}

func Create(configuration *config.Config, campaign *schema.Campaign) map[string]interface{} {
	if resp, ok := Validate(configuration, campaign); !ok {
		return resp
	}

	configuration.Database.Create(campaign)

	if campaign.ID <= 0 {
		return u.Message(false, "Failed to create campaign")
	}

	// remove user info

	response := u.Message(true, "Campaign created")
	response["campaign"] = campaign

	return response
}

func CreateACampaign(configuration *config.Config) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		campaign := &schema.Campaign{}

		err := render.DecodeJSON(r.Body, campaign)

		if err != nil {
			u.Respond(w, r, u.Message(false, "Invalid request"))
		} else {
			resp := Create(configuration, campaign)
			u.Respond(w, r, resp)
		}
	})
}

func FindCampaign(configuration *config.Config, id string) *schema.Campaign {
	campaign := &schema.Campaign{}

	configuration.Database.Preload("Maps").Table("campaigns").Where("id = ?", id).First(campaign)

	if campaign.ID == 0 {
		return nil
	}

	return campaign
}

func GetACampaign(configuration *config.Config) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		campaignID := chi.URLParam(r, "campaignId")
		campaign := FindCampaign(configuration, campaignID)

		if campaign == nil {
			u.Respond(w, r, u.Message(false, "Could not find campaign"))
		} else {
			response := u.Message(true, "Campaign found")
			response["campaign"] = campaign

			u.Respond(w, r, response)
		}
	})
}

/** MAP **/

func ValidateMap(configuration *config.Config, m *schema.Map) (map[string]interface{}, bool) {
	if m.CampaignID <= 0 {
		return u.Message(false, "Campaigns need to be owned by a user"), false
	}

	if m.Name == "" {
		return u.Message(false, "Map must have a name"), false
	}

	return u.Message(true, "Valid"), true
}

func CreateMap(configuration *config.Config, m *schema.Map) map[string]interface{} {
	if resp, ok := ValidateMap(configuration, m); !ok {
		return resp
	}

	configuration.Database.Create(m)

	if m.ID <= 0 {
		return u.Message(false, "Failed to create map")
	}

	response := u.Message(true, "Map created")
	response["map"] = m

	return response
}

func CreateAMap(configuration *config.Config) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		campaignID, _ := strconv.Atoi(chi.URLParam(r, "campaignId"))
		m := &schema.Map{}

		err := render.DecodeJSON(r.Body, m)

		m.CampaignID = campaignID

		if err != nil {
			u.Respond(w, r, u.Message(false, "Invalid request"))
		} else {
			resp := CreateMap(configuration, m)
			u.Respond(w, r, resp)
		}
	})
}
