package schema

type UserWithOutPassword struct {
	Email    string `json:"-"`
	Username string `json:"username"`
	Token    string `json:"token" sql:"-"`
	Maps     []Map  `gorm:"foreignkey:UserID" json:"maps"`
	Password string `json:"-"`
}

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

// //BaseMap is safe to edit via API
// type BaseMap struct {
// 	Name         string     `json:"name"`
// 	Description  string     `json:"description"`
// 	DownloadCode string     `json:"downloadCode"`
// 	Type         string     `json:"type"`
// 	Photos       []MapPhoto `json:"photos"`
// }

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
}
