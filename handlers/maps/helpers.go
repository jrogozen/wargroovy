package maps

import (
	"github.com/jrogozen/wargroovy/internal/config"
	"github.com/jrogozen/wargroovy/schema"
	u "github.com/jrogozen/wargroovy/utils"
	"net/http"
	"strconv"
)

//Validate validates map fields for map creation
func Validate(configuration *config.Config, claims map[string]interface{}, m *schema.Map) (map[string]interface{}, bool) {
	if m.UserID <= 0 {
		return u.Message(false, "maps need to be owned by a user"), false
	}

	if _, ok := u.IsUserAuthorized(m.UserID, claims); !ok {
		return u.Message(false, "Can only create or edit a map associated with valid userId"), false
	}

	if m.Name == "" {
		return u.Message(false, "Maps must have a name"), false
	}

	return u.Message(true, "Valid"), true
}

func ValidateUpdate(configuration *config.Config, claims map[string]interface{}, m *schema.Map) (map[string]interface{}, bool) {
	if _, ok := u.IsUserAuthorized(m.UserID, claims); !ok {
		return u.Message(false, "Can only create or edit a map associated with valid userId"), false
	}

	return u.Message(true, "Valid"), true
}

//GetSortOptions parse queryParams and return correctly typed SortOptions
func GetSortOptions(r *http.Request) *schema.SortOptions {
	options := &schema.SortOptions{}

	limit := u.StringWithDefault(r.URL.Query().Get("limit"), "20")
	offset := u.StringWithDefault(r.URL.Query().Get("offset"), "0")
	orderBy := u.StringWithDefault(r.URL.Query().Get("orderBy"), "created_descending")

	limitInt, _ := strconv.Atoi(limit)
	offsetInt, _ := strconv.Atoi(offset)

	options.Limit = limitInt
	options.Offset = offsetInt

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

	options.OrderBy = orderQueryString

	return options
}
