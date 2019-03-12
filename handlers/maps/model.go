package maps

import (
	"github.com/jrogozen/wargroovy/internal/config"
	"github.com/jrogozen/wargroovy/schema"
	u "github.com/jrogozen/wargroovy/utils"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

func UpdateMap(configuration *config.Config, claims map[string]interface{}, mapIDString string, m *schema.Map) (map[string]interface{}, int) {
	mapID, _ := strconv.ParseInt(mapIDString, 10, 64)

	originalMap, err := configuration.DB.GetMap(mapID)

	if err != nil {
		return u.Message(false, err.Error()), http.StatusBadRequest
	}

	if resp, ok := ValidateUpdate(configuration, claims, originalMap); !ok {
		return resp, http.StatusForbidden
	}

	originalMap.Merge(m)

	insertedID, err := configuration.DB.UpdateMap(originalMap)

	if err != nil {
		return u.Message(false, err.Error()), http.StatusBadRequest
	}

	//TODO missing created at, updated at
	updatedMap := &schema.Map{
		ID:           insertedID,
		Name:         originalMap.Name,
		Description:  originalMap.Description,
		DownloadCode: originalMap.DownloadCode,
		Type:         originalMap.Type,
		UserID:       originalMap.UserID,
		Views:        originalMap.Views,
		Slug:         originalMap.Slug,
		Photos:       originalMap.Photos,
	}

	response := u.Message(true, "Map updated")
	response["map"] = updatedMap

	return response, http.StatusOK
}

func FindMap(configuration *config.Config, mapIDString string, claims map[string]interface{}) (map[string]interface{}, int) {
	mapID, _ := strconv.ParseInt(mapIDString, 10, 64)

	m, err := configuration.DB.GetMap(mapID)

	if err != nil {
		return u.Message(false, err.Error()), http.StatusBadRequest
	}

	_, err = configuration.DB.IncrementMapView(mapID)

	if err != nil {
		log.WithField("mapId", mapID).Warn("Could not increment map view")
	}

	if claims["UserID"] != nil {
		userID := int64(claims["UserID"].(float64))

		rating, err := configuration.DB.GetMapUserRating(mapID, userID)

		if err == nil {
			m.UserRating = "thumbs_down"

			if rating > 1 {
				m.UserRating = "thumbs_up"
			}
		} else {
			m.UserRating = "not_rated"
		}
	} else {
		m.UserRating = "not_rated"
	}

	response := u.Message(true, "Map found")
	response["map"] = m

	return response, http.StatusOK

}

func FindMapBySlug(configuration *config.Config, slug string, claims map[string]interface{}) (map[string]interface{}, int) {
	m, err := configuration.DB.GetMapBySlug(slug)

	if err != nil {
		return u.Message(false, err.Error()), http.StatusBadRequest
	}

	if claims["UserID"] != nil {
		userID := int64(claims["UserID"].(float64))

		rating, err := configuration.DB.GetMapUserRating(m.ID, userID)

		if err == nil {
			m.UserRating = "thumbs_down"

			if rating > 1 {
				m.UserRating = "thumbs_up"
			}
		} else {
			m.UserRating = "not_rated"
		}
	} else {
		m.UserRating = "not_rated"
	}

	_, err = configuration.DB.IncrementMapView(m.ID)

	if err != nil {
		log.WithField("mapId", m.ID).Warn("Could not increment map view")
	}

	response := u.Message(true, "Map found")
	response["map"] = m

	return response, http.StatusOK
}

func Create(configuration *config.Config, claims map[string]interface{}, m *schema.Map) (map[string]interface{}, int) {
	if resp, ok := Validate(configuration, claims, m); !ok {
		return resp, http.StatusForbidden
	}

	insertedID, err := configuration.DB.AddMap(m)

	if err != nil || insertedID == 0 {
		log.Error(err)

		return u.Message(false, err.Error()), http.StatusBadRequest
	}
	response := u.Message(true, "Map created")
	response["map"] = insertedID

	return response, http.StatusOK
}

// //FindMapList returns ordered list of maps
func FindMapList(configuration *config.Config, options *schema.SortOptions) (map[string]interface{}, int) {
	maps, err := configuration.DB.ListByMap(options)

	if err != nil {
		log.Error(err)

		return u.Message(false, err.Error()), http.StatusBadRequest
	}

	response := u.Message(true, "Maps found")
	response["maps"] = maps

	return response, http.StatusOK
}

func DeletePhoto(configuration *config.Config, claims map[string]interface{}, mapIDString string, url string) (map[string]interface{}, int) {
	mapID, _ := strconv.ParseInt(mapIDString, 10, 64)

	originalMap, err := configuration.DB.GetMap(mapID)

	if err != nil {
		return u.Message(false, err.Error()), http.StatusBadRequest
	}

	if resp, ok := ValidateUpdate(configuration, claims, originalMap); !ok {
		return resp, http.StatusForbidden
	}

	numPhotosDeleted, err := configuration.DB.DeleteMapPhoto(mapID, url)

	if err != nil {
		return u.Message(false, err.Error()), http.StatusBadRequest
	}

	response := u.Message(true, "Photo deleted")
	response["deleted"] = numPhotosDeleted

	return response, http.StatusOK
}

func Delete(configuration *config.Config, claims map[string]interface{}, mapIDString string) (map[string]interface{}, int) {
	mapID, _ := strconv.ParseInt(mapIDString, 10, 64)

	originalMap, err := configuration.DB.GetMap(mapID)

	if err != nil {
		return u.Message(false, err.Error()), http.StatusBadRequest
	}

	if resp, ok := ValidateUpdate(configuration, claims, originalMap); !ok {
		return resp, http.StatusForbidden
	}

	numMapsDeleted, err := configuration.DB.DeleteMap(mapID)

	if err != nil {
		return u.Message(false, err.Error()), http.StatusBadRequest
	}

	response := u.Message(true, "Map deleted")
	response["deleted"] = numMapsDeleted

	return response, http.StatusOK
}

func Rate(configuration *config.Config, claims map[string]interface{}, mapIDString string, rating int64) (map[string]interface{}, int) {
	mapID, _ := strconv.ParseInt(mapIDString, 10, 64)
	userID, err := u.GetUserIdFromClaims(claims)

	if err != nil {
		return u.Message(false, err.Error()), http.StatusForbidden
	}

	insertedRating, err := configuration.DB.RateMap(mapID, userID, rating)

	if err != nil {
		return u.Message(false, err.Error()), http.StatusBadRequest
	}

	response := u.Message(true, "Map rated")
	response["rating"] = insertedRating

	return response, http.StatusOK
}

func FindMapListTags(configuration *config.Config, options *schema.TagSortOptions) (map[string]interface{}, int) {
	tags, err := configuration.DB.GetMapListTags(options.OrderBy, options.Limit)

	if err != nil {
		log.Error(err)

		return u.Message(false, err.Error()), http.StatusBadRequest
	}

	response := u.Message(true, "Tags found")
	response["tags"] = tags

	return response, http.StatusOK
}
