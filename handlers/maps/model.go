package maps

import (
	"github.com/jrogozen/wargroovy/internal/config"
	"github.com/jrogozen/wargroovy/schema"
	u "github.com/jrogozen/wargroovy/utils"
	log "github.com/sirupsen/logrus"
	"strconv"
)

func UpdateMap(configuration *config.Config, claims map[string]interface{}, mapIdString string, m *schema.Map) map[string]interface{} {
	mapID, _ := strconv.ParseInt(mapIdString, 10, 64)

	originalMap, err := configuration.DB.GetMap(mapID)

	if resp, ok := ValidateUpdate(configuration, claims, originalMap); !ok {
		return resp
	}

	if err != nil {
		return u.Message(false, err.Error())
	}

	originalMap.Merge(m)

	insertedID, err := configuration.DB.UpdateMap(originalMap)

	if err != nil {
		return u.Message(false, err.Error())
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
		Photos:       originalMap.Photos,
	}

	response := u.Message(true, "Map updated")
	response["map"] = updatedMap

	return response
}

func FindMap(configuration *config.Config, mapIdString string) map[string]interface{} {
	mapID, _ := strconv.ParseInt(mapIdString, 10, 64)
	m, err := configuration.DB.GetMap(mapID)

	if err != nil {
		return u.Message(false, err.Error())
	}

	response := u.Message(true, "Map found")
	response["map"] = m

	return response

}

func Create(configuration *config.Config, claims map[string]interface{}, m *schema.Map) map[string]interface{} {
	if resp, ok := Validate(configuration, claims, m); !ok {
		return resp
	}

	insertedID, err := configuration.DB.AddMap(m)

	if err != nil || insertedID == 0 {
		log.Error(err)

		return u.Message(false, err.Error())
	}

	//TODO missing created at, updated at
	createdMap := &schema.Map{
		ID:           insertedID,
		Name:         m.Name,
		Description:  m.Description,
		DownloadCode: m.DownloadCode,
		Type:         m.Type,
		UserID:       m.UserID,
		Views:        0,
		Photos:       m.Photos,
	}

	response := u.Message(true, "Map created")
	response["map"] = createdMap

	return response
}

// func IncrementMapView(configuration *config.Config, m schema.Map) schema.Map {
// 	configuration.Database.Model(&m).Update("views", m.Views+1)

// 	return m
// }

// //FindMapList returns ordered list of maps
func FindMapList(configuration *config.Config, options *schema.SortOptions) map[string]interface{} {
	maps, err := configuration.DB.ListByMap(options)

	if err != nil {
		log.Error(err)

		return u.Message(false, err.Error())
	}

	response := u.Message(true, "Maps found")
	response["maps"] = maps

	return response
}
