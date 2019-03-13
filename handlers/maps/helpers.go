package maps

import (
	"fmt"
	"github.com/jrogozen/wargroovy/internal/config"
	"github.com/jrogozen/wargroovy/schema"
	u "github.com/jrogozen/wargroovy/utils"
	// log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"strings"
)

//Validate validates map fields for map creation
func Validate(configuration *config.Config, claims map[string]interface{}, m *schema.Map) (map[string]interface{}, bool) {
	if claims["UserID"] == nil {
		return u.Message(false, "Maps need to be owned by a user"), false
	}

	if m.UserID <= 0 && claims["UserID"] != nil {
		// map is missing userID, set it based on claim
		m.UserID = int64(claims["UserID"].(float64))
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

//TODO: refactor this to be used elsewhere
func appendToQueryCondition(s string, append string) string {
	if s == "" {
		return fmt.Sprintf("where (%s", append)
	}

	return fmt.Sprintf("%s OR %s", s, append)
}

func appendToExistingQueryCondition(s string, append string) string {
	if s == "" {
		return fmt.Sprintf("and (%s", append)
	}

	return fmt.Sprintf("%s OR %s", s, append)
}

func constructTypeQueryString(types []string) string {
	const identifier = "m.type = '%s'"

	typeQueryString := ""

	for _, t := range types {
		typeQueryString = appendToQueryCondition(typeQueryString, fmt.Sprintf(identifier, t))
	}

	return typeQueryString
}

func constructTagQueryString(tags []string) string {
	const identifierStart = "tags like '%"
	const identifierEnd = "%'"
	tagQueryString := ""

	for _, tg := range tags {
		tagQueryString = appendToExistingQueryCondition(tagQueryString, fmt.Sprint(identifierStart, tg, identifierEnd))
	}

	return tagQueryString
}

func GetTagSortOptions(r *http.Request) *schema.TagSortOptions {
	options := &schema.TagSortOptions{}

	limit := u.StringWithDefault(r.URL.Query().Get("limit"), "20")
	orderBy := u.StringWithDefault(r.URL.Query().Get("orderBy"), "count_descending")

	limitInt, _ := strconv.Atoi(limit)

	options.Limit = limitInt

	orderQueryString := ""

	switch orderBy {
	case "count_descending":
		orderQueryString = "count desc"
	case "count_ascending":
		orderQueryString = "count asc"
	case "name_ascending":
		orderQueryString = "tag_name asc"
	case "name_descending":
		orderQueryString = "tag_name desc"
	default:
		orderQueryString = "count desc"
	}

	options.OrderBy = orderQueryString

	return options
}

//GetSortOptions parse queryParams and return correctly typed SortOptions
func GetSortOptions(r *http.Request) *schema.SortOptions {
	options := &schema.SortOptions{}

	limit := u.StringWithDefault(r.URL.Query().Get("limit"), "20")
	offset := u.StringWithDefault(r.URL.Query().Get("offset"), "0")
	orderBy := u.StringWithDefault(r.URL.Query().Get("orderBy"), "created_descending")
	t := u.StringWithDefault(r.URL.Query().Get("type"), "all")
	tg := u.StringWithDefault(r.URL.Query().Get("tags"), "all")

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
	case "alphabetical_ascending":
		orderQueryString = "name asc"
	case "alphabetical_descending":
		orderQueryString = "name desc"
	case "rating_ascending":
		orderQueryString = "rating asc nulls first"
	case "rating_descending":
		orderQueryString = "rating desc nulls last"
	default:
		orderQueryString = "created_at desc"
	}

	options.OrderBy = orderQueryString

	typeQueryString := ""

	if t == "all" {
		typeQueryString = "WHERE (m.type is not null)"
	} else {
		filterFunc := func(s string) bool {
			return strings.Contains(s, "scenario") || strings.Contains(s, "skirmish") || strings.Contains(s, "puzzle")
		}

		types := u.Choose(strings.Split(t, ","), filterFunc)

		typeQueryString = constructTypeQueryString(types) + ")"
	}

	options.Type = typeQueryString

	tagsQueryString := ""

	if tg != "all" {
		tags := strings.Split(tg, ",")

		tagsQueryString = constructTagQueryString(tags) + ")"
	}

	options.Tags = tagsQueryString

	return options
}
