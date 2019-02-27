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

func FindMap(configuration *config.Config, mapIDString string) (map[string]interface{}, int) {
	mapID, _ := strconv.ParseInt(mapIDString, 10, 64)

	m, err := configuration.DB.GetMap(mapID)

	if err != nil {
		return u.Message(false, err.Error()), http.StatusBadRequest
	}

	response := u.Message(true, "Map found")
	response["map"] = m

	return response, http.StatusOK

}

func FindMapBySlug(configuration *config.Config, slug string) (map[string]interface{}, int) {
	m, err := configuration.DB.GetMapBySlug(slug)

	if err != nil {
		return u.Message(false, err.Error()), http.StatusBadRequest
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

// func IncrementMapView(configuration *config.Config, m schema.Map) schema.Map {
// 	configuration.Database.Model(&m).Update("views", m.Views+1)

// 	return m
// }

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
