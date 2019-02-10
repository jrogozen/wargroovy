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

type Campaign struct {
	gorm.Model

	Name              string `json:"name"`
	Description       string `json:"description"`
	ThumbPhotoURL     string `json:"thumbPhotoUrl"`
	LargePhotoURL     string `json:"largePhotoUrl"`
	SingleMapCampaign bool   `gorm:"DEFAULT:true" json:"singleMapCampaign"`
	Views             int    `json:"views"`
	UserID            int    `json:"userId"`
	Maps              []Map  `gorm:"foreignkey:CampaignID"`
}

type Map struct {
	gorm.Model

	Name          string `json:"name"`
	Description   string `json:"description"`
	ThumbPhotoURL string `json:"thumbPhotoUrl"`
	LargePhotoURL string `json:"largePhotoUrl"`
	DownloadCode  string `json:"downloadCode"`
	CampaignID    int    `json:"campaignId"`
}
