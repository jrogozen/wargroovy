package schema

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model

	Email     string     `json:"email"`
	Password  string     `json:"-"`
	Campaigns []Campaign `gorm:"foreignkey:UserID"`
}

//BaseCampaign is safe to edit via API
type BaseCampaign struct {
	Name              string `json:"name"`
	Description       string `json:"description"`
	ThumbPhotoURL     string `json:"thumbPhotoUrl"`
	LargePhotoURL     string `json:"largePhotoUrl"`
	SingleMapCampaign bool   `gorm:"DEFAULT:true" json:"singleMapCampaign"`
}

type Campaign struct {
	gorm.Model
	*BaseCampaign
	UserID int   `json:"userId"`
	Views  int   `json:"views"`
	Maps   []Map `gorm:"foreignkey:CampaignID"`
}

//BaseMap is safe to edit via API
type BaseMap struct {
	Name          string `json:"name"`
	Description   string `json:"description"`
	ThumbPhotoURL string `json:"thumbPhotoUrl"`
	LargePhotoURL string `json:"largePhotoUrl"`
	DownloadCode  string `json:"downloadCode"`
}

type Map struct {
	gorm.Model
	*BaseMap

	CampaignID int `json:"campaignId"`
}
