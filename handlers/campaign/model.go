package campaign

import (
	"github.com/jrogozen/wargroovy/internal/config"
	"github.com/jrogozen/wargroovy/schema"
	u "github.com/jrogozen/wargroovy/utils"
)

func UpdateCampaign(configuration *config.Config, claims map[string]interface{}, campaign *schema.Campaign, updatedCampaign *schema.BaseCampaign) map[string]interface{} {
	if resp, ok := Validate(configuration, claims, campaign); !ok {
		return resp
	}

	configuration.Database.Model(&campaign).Updates(&updatedCampaign)

	if campaign.ID <= 0 {
		response := u.Message(false, "Campaign could not be created")

		return response
	}

	response := u.Message(true, "Campaign updated")
	response["campaign"] = campaign

	return response
}

func FindMap(configuration *config.Config, id string) *schema.Map {
	m := &schema.Map{}

	configuration.Database.Table("maps").Where("id = ?", id).First(m)

	if m.ID == 0 {
		return nil
	}

	return m
}

func UpdateMap(configuration *config.Config, claims map[string]interface{}, campaign *schema.Campaign, mapID string, updatedMap *schema.BaseMap) map[string]interface{} {
	m := FindMap(configuration, mapID)

	if m == nil {
		return u.Message(false, "Could not find map to update")
	}

	if campaign.ID != m.CampaignID {
		return u.Message(false, "Map does not belong to given campaign")
	}

	if resp, ok := ValidateMap(configuration, claims, m, campaign); !ok {
		return resp
	}

	configuration.Database.Model(&m).Updates(&updatedMap)

	response := u.Message(true, "Map updated")
	response["map"] = m

	return response
}

func FindCampaign(configuration *config.Config, id string) *schema.Campaign {
	campaign := &schema.Campaign{}

	configuration.Database.Preload("Maps").Table("campaigns").Where("id = ?", id).First(campaign)

	if campaign.ID == 0 {
		return nil
	}

	return campaign
}

func Create(configuration *config.Config, claims map[string]interface{}, campaign *schema.Campaign) map[string]interface{} {
	if resp, ok := Validate(configuration, claims, campaign); !ok {
		return resp
	}

	configuration.Database.Create(campaign)

	if campaign.ID <= 0 {
		response := u.Message(false, "Campaign could not be created")

		return response
	}

	response := u.Message(true, "Campaign created")
	response["campaign"] = campaign

	return response
}

func CreateMap(configuration *config.Config, claims map[string]interface{}, m *schema.Map, campaign *schema.Campaign) map[string]interface{} {
	if resp, ok := ValidateMap(configuration, claims, m, campaign); !ok {
		return resp
	}

	configuration.Database.Create(m)

	if m.ID <= 0 {
		response := u.Message(false, "Failed to create map")

		return response
	}

	response := u.Message(true, "Map created")
	response["map"] = m

	return response
}

func IncrementCampaignView(configuration *config.Config, campaign schema.Campaign) schema.Campaign {
	configuration.Database.Model(&campaign).Update("views", campaign.Views+1)

	return campaign
}

//FindCampaignList returns ordered list of campaigns
func FindCampaignList(configuration *config.Config, options *SortOptions) []*schema.Campaign {
	campaigns := make([]*schema.Campaign, 0, options.limit)

	configuration.Database.
		Preload("Maps").
		Table("campaigns").
		Order(options.orderBy).
		Limit(options.limit).
		Offset(options.offset).
		Find(&campaigns)

	return campaigns
}
