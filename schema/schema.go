package schema

import (
	"github.com/rs/xid"
	"strings"
)

//User not safe to return ever. corresponds to every row in table
// used when fetching from DB
type User struct {
	ID        int64
	CreatedAt int
	UpdatedAt int
	Email     string
	Username  string
	Password  string
	Token     string `json:"token"`
}

//CreateUser what to accept to create a user
type CreateUser struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

//CreatedUser user struct to send back in create/login flows
type CreatedUser struct {
	ID       int64  `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Token    string `json:"token"`
}

//UserView user struct when a user queries herself
type UserView struct {
	ID       int64  `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

type SecureUserView struct {
	ID        int64  `json:"ID"`
	CreatedAt int    `json:"createdAt"`
	UpdatedAt int    `json:"updatedAt"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	Token     string `json:"token"`
}

type MapPhoto struct {
	ID    int64  `json:"id"`
	MapID int64  `json:"map_id"`
	URL   string `json:"url"`
}

//Map not safe to return ever. corresponds to every row in table
// used when fetching from DB
type Map struct {
	ID           int64    `json:"id"`
	CreatedAt    int      `json:"created_at"`
	UpdatedAt    int      `json:"updated_at"`
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	DownloadCode string   `json:"downloadCode"`
	Type         string   `json:"type"`
	UserID       int64    `json:"userId" sql:"type:integer REFERENCES users(id)"`
	Views        int      `json:"views"`
	Photos       []string `json:"photos"`
	Slug         string   `json:"slug"`
}

func (m *Map) Merge(u *Map) {
	if u.Name != "" {
		m.Name = u.Name

		// need to generate a new slug
		slug := strings.Replace(m.Name, " ", "-", -1) + "-" + xid.New().String()

		m.Slug = slug
	}

	if u.Description != "" {
		m.Description = u.Description
	}

	if u.DownloadCode != "" {
		m.DownloadCode = u.DownloadCode
	}

	if u.Type != "" {
		m.Type = u.Type
	}
}

type MapFromDB struct {
	ID           int64
	CreatedAt    int
	UpdatedAt    int
	Name         string
	Description  string
	DownloadCode string
	Type         string
	UserID       int64
	Views        int
	Photos       string
}

//SortOptions used by FindCampaignList
type SortOptions struct {
	Limit   int
	Offset  int
	OrderBy string
}
