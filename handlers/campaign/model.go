package campaign

import (
	"github.com/jrogozen/wargroovy/internal/config"
	"github.com/jrogozen/wargroovy/schema"
	u "github.com/jrogozen/wargroovy/utils"
)

func UpdateCampaign(configuration *config.Config, claims map[string]interface{}, campaign *schema.Campaign, updatedCampaign *schema.BaseCampaign) *schema.Campaign {
	campaignUserID := campaign.UserID

	if _, ok := u.IsUserAuthorized(campaignUserID, claims); !ok {
		return nil
	}

	configuration.Database.Model(&campaign).Updates(&updatedCampaign)

	return campaign
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

func FindCampaign(configuration *config.Config, id string) *schema.Campaign {
	campaign := &schema.Campaign{}

	configuration.Database.Preload("Maps").Table("campaigns").Where("id = ?", id).First(campaign)

	if campaign.ID == 0 {
		return nil
	}

	return campaign
}

func Create(configuration *config.Config, campaign *schema.Campaign) map[string]interface{} {
	if resp, ok := Validate(configuration, campaign); !ok {
		return resp
	}

	configuration.Database.Create(campaign)

	if campaign.ID <= 0 {
		return nil
	}

	response := u.Message(true, "Campaign created")
	response["campaign"] = campaign

	return response
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
