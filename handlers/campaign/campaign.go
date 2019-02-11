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

func Routes(configuration *config.Config) *chi.Mux {
	router := chi.NewRouter()

	// var tokenAuth *jwtauth.JWTAuth

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
			return
		}

		resp := Create(configuration, campaign)
		u.Respond(w, r, resp)
		return
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

func IncrementCampaignView(configuration *config.Config, campaign schema.Campaign) schema.Campaign {
	configuration.Database.Model(&campaign).Update("views", campaign.Views+1)

	return campaign
}

/*FindCampaignList returns ordered list of campaigns
{ orderBy, limit, start }
*/
func FindCampaignList(configuration *config.Config, orderBy string, limit string, offset string) []*schema.Campaign {
	limitInt, _ := strconv.Atoi(limit)
	offsetInt, _ := strconv.Atoi(offset)
	campaigns := make([]*schema.Campaign, 0, limitInt)

	orderQueryString := ""

	switch orderBy {
	case "created_ascending":
		orderQueryString = "created_at asc"
	case "created_descending":
		orderQueryString = "created_at desc"
	case "views_ascending":
		orderQueryString = "views asc"
	case "views_descending":
		orderQueryString = "views desc"
	case "alphabetical_asc":
		orderQueryString = "name asc"
	case "alphabetical_desc":
		orderQueryString = "name desc"
	default:
		orderQueryString = "created_at desc"
	}

	configuration.Database.
		Preload("Maps").
		Table("campaigns").
		Order(orderQueryString).
		Limit(limitInt).
		Offset(offsetInt).
		Find(&campaigns)

	return campaigns
}

func UpdateCampaign(configuration *config.Config, claims map[string]interface{}, campaign *schema.Campaign, updatedCampaign *schema.BaseCampaign) *schema.Campaign {
	campaignUserID := campaign.UserID

	if _, ok := u.IsUserAuthorized(campaignUserID, claims); !ok {
		return nil
	}

	configuration.Database.Model(&campaign).Updates(&updatedCampaign)

	return campaign
}

func GetACampaign(configuration *config.Config) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		campaignID := chi.URLParam(r, "campaignId")
		campaign := FindCampaign(configuration, campaignID)

		// increment view counter
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
		limit := u.StringWithDefault(r.URL.Query().Get("limit"), "20")
		offset := u.StringWithDefault(r.URL.Query().Get("offset"), "-1")
		orderBy := u.StringWithDefault(r.URL.Query().Get("orderBy"), "created_descending")

		campaigns := FindCampaignList(configuration, orderBy, limit, offset)

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

func FindMap(configuration *config.Config, id string) *schema.Map {
	m := &schema.Map{}

	configuration.Database.Table("maps").Where("id = ?", id).First(m)

	if m.ID == 0 {
		return nil
	}

	return m
}

func UpdateMap(configuration *config.Config, campaignId string, mapId string, updatedMap *schema.BaseMap) *schema.Map {
	// find campaign to get userId so we can verify that this user can update this map
	// once we implement jwt
	campaign := FindCampaign(configuration, campaignId)

	if campaign == nil {
		return nil
	}

	m := FindMap(configuration, mapId)

	if m == nil || uint(m.CampaignID) != campaign.ID {
		return nil
	} else {
		configuration.Database.Model(&m).Updates(&updatedMap)

		return m
	}
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
