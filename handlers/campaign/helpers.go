package campaign

import (
	"github.com/jrogozen/wargroovy/internal/config"
	"github.com/jrogozen/wargroovy/schema"
	u "github.com/jrogozen/wargroovy/utils"
	"net/http"
	"strconv"
)

//SortOptions used by FindCampaignList
type SortOptions struct {
	limit   int
	offset  int
	orderBy string
}

//Validate validates campaign fields for campaign creation
func Validate(configuration *config.Config, campaign *schema.Campaign) (map[string]interface{}, bool) {
	if campaign.UserID <= 0 {
		return u.Message(false, "Campaigns need to be owned by a user"), false
	}

	return u.Message(true, "Valid"), true
}

//ValidateMap validates map fields for map creation
func ValidateMap(configuration *config.Config, m *schema.Map) (map[string]interface{}, bool) {
	if m.CampaignID <= 0 {
		return u.Message(false, "Campaigns need to be owned by a user"), false
	}

	if m.Name == "" {
		return u.Message(false, "Map must have a name"), false
	}

	return u.Message(true, "Valid"), true
}

//GetSortOptions parse queryParams and return correctly typed SortOptions
func GetSortOptions(r *http.Request) *SortOptions {
	options := &SortOptions{}

	limit := u.StringWithDefault(r.URL.Query().Get("limit"), "20")
	offset := u.StringWithDefault(r.URL.Query().Get("offset"), "-1")
	orderBy := u.StringWithDefault(r.URL.Query().Get("orderBy"), "created_descending")

	limitInt, _ := strconv.Atoi(limit)
	offsetInt, _ := strconv.Atoi(offset)

	options.limit = limitInt
	options.offset = offsetInt

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

	options.orderBy = orderQueryString

	return options
}
