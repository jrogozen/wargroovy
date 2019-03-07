package schema

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"github.com/rs/xid"
	// log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

//User not safe to return ever. corresponds to every row in table
// used when fetching from DB
type User struct {
	ID        int64  `json:"id"`
	CreatedAt int    `json:"created_at"`
	UpdatedAt int    `json:"updated_at"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Token     string `json:"token"`
}

// handle what properties can be updated by a user and when
func (user *User) Merge(update *User) {
	if update.Password != "" {
		// need to hash a new password
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(update.Password), bcrypt.DefaultCost)
		user.Password = string(hashedPassword)
	}

	if update.Email != "" {
		user.Email = update.Email
	}

	if update.Username != "" {
		user.Username = update.Username
	}
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
	ID           int64          `json:"id"`
	CreatedAt    int            `json:"created_at"`
	UpdatedAt    int            `json:"updated_at"`
	Name         string         `json:"name"`
	Description  DescriptionMap `json:"description"`
	DownloadCode string         `json:"download_code"`
	Type         string         `json:"type"`
	UserID       int64          `json:"userId" sql:"type:integer REFERENCES users(id)"`
	Views        int            `json:"views"`
	Photos       []string       `json:"photos"`
	Slug         string         `json:"slug"`
	Author       string         `json:"author"`
	Rating       float64        `json:"rating"`
	UserRating   string         `json:"user_rating"`
}

func (m *Map) Merge(u *Map) {
	if u.Name != "" {
		m.Name = u.Name

		// need to generate a new slug
		slug := strings.ToLower(strings.Replace(m.Name, " ", "-", -1)) + "-" + xid.New().String()

		m.Slug = slug
	}

	if u.Description != nil {
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
	Description  DescriptionMap
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
	Type    string
}

type DescriptionMap map[string]interface{}

func (d DescriptionMap) Value() (driver.Value, error) {
	j, err := json.Marshal(d)
	return j, err
}

func (d *DescriptionMap) Scan(src interface{}) error {
	source, ok := src.([]byte)

	if !ok {
		return nil
	}

	var i interface{}
	err := json.Unmarshal(source, &i)

	if err != nil {
		return err
	}

	*d, ok = i.(map[string]interface{})

	if !ok {
		return errors.New("Type assertion .(map[string]interface{}) failed.")
	}

	return nil
}
