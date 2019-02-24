package maps

import (
	"github.com/jrogozen/wargroovy/internal/config"
	"github.com/jrogozen/wargroovy/schema"
	u "github.com/jrogozen/wargroovy/utils"
	log "github.com/sirupsen/logrus"
	"strconv"
)

// func UpdateMap(configuration *config.Config, claims map[string]interface{}, m *schema.Map, updatedMap *schema.BaseMap) map[string]interface{} {
// 	if resp, ok := Validate(configuration, claims, m); !ok {
// 		return resp
// 	}

// 	configuration.Database.Model(&m).Updates(&updatedMap)

// 	if m.ID <= 0 {
// 		response := u.Message(false, "Map could not be created")

// 		return response
// 	}

// 	response := u.Message(true, "Map updated")
// 	response["map"] = m

// 	return response
// }

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
// func FindMapList(configuration *config.Config, options *SortOptions) []*schema.Map {
// 	maps := make([]*schema.Map, 0, options.limit)

// 	configuration.Database.
// 		Preload("Photos").
// 		Table("maps").
// 		Order(options.orderBy).
// 		Limit(options.limit).
// 		Offset(options.offset).
// 		Find(&maps)

// 		// mapsWithUser := make([]*schema.MapWithUser, 0, options.limit)

// 		// for _, m := range maps {
// 		// 	user := schema.User{}
// 		// 	// mapWithUser := &schema.MapWithUser{}

// 		// 	configuration.Database.Model(&m).Related(&user)

// 		// 	userForMap := schema.UserForMap{
// 		// 		Username: user.Username,
// 		// 	}

// 	// mapWithUser.User = *user
// 	// mapWithUser.Name = m.Name
// 	log.Info("hello")
// 	// mapsWithUser = append(mapsWithUser, mapWithUser)
// 	// }

// 	return maps
// }
