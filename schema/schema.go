package schema

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
)

type UserWithOutPassword struct {
	*User

	Password string `json:"-"`
}

type User struct {
	gorm.Model

	Email     string     `json:"email"`
	Token     string     `json:"token" sql:"-"`
	Password  string     `json:"password"`
	Campaigns []Campaign `gorm:"foreignkey:UserID"`
}

type TokenClaims struct {
	jwt.StandardClaims

	UserID uint
	Admin  bool
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
	UserID uint  `json:"userId"`
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

	CampaignID uint `json:"campaignId"`
}
