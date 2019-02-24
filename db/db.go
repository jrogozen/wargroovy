package db

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jrogozen/wargroovy/schema"
	"github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"time"
)

type PsqlDB struct {
	Conn *sql.DB

	users userSQL
	maps  mapSQL
}

type userSQL struct {
	// list       *sql.Stmt
	// listBy     *sql.Stmt
	insert     *sql.Stmt
	get        *sql.Stmt
	getByLogin *sql.Stmt
	// update     *sql.Stmt
	// delete     *sql.Stmt
}

type mapSQL struct {
	insert       *sql.Stmt
	insertPhotos *sql.Stmt
	get          *sql.Stmt
	getPhotos    *sql.Stmt

	// list   *sql.Stmt
	// listBy *sql.Stmt
	// update *sql.Stmt
}

func NewPostgresDB(connectionString string) (*PsqlDB, error) {
	log.Info(connectionString)

	conn, err := sql.Open("postgres", connectionString)

	if err != nil {
		panic(fmt.Errorf("DB: %v", err))
	}

	db := &PsqlDB{
		Conn:  conn,
		users: userSQL{},
	}

	if err = db.Conn.Ping(); err != nil {
		db.Conn.Close()
		return nil, fmt.Errorf("psql: could not establish a connection: %v", err)
	}

	var userInsert *sql.Stmt
	var userGet *sql.Stmt
	var userGetByLogin *sql.Stmt

	var mapInsert *sql.Stmt
	var mapPhotosInsert *sql.Stmt
	var mapGet *sql.Stmt
	var mapPhotosGet *sql.Stmt

	if userInsert, err = db.Conn.Prepare(insertUserStatement); err != nil {
		return nil, fmt.Errorf("psql: prepare user insert: %v", err)
	}

	if userGet, err = db.Conn.Prepare(getUserStatement); err != nil {
		return nil, fmt.Errorf("psql: prepare user get: %v", err)
	}

	if userGetByLogin, err = db.Conn.Prepare(getUserByLoginStatement); err != nil {
		return nil, fmt.Errorf("psql: prepare user getByLogin: %v", err)
	}

	if mapInsert, err = db.Conn.Prepare(insertMapStatement); err != nil {
		return nil, fmt.Errorf("psql: prepare map insert: %v", err)
	}

	if mapPhotosInsert, err = db.Conn.Prepare(insertMapPhotosStatement); err != nil {
		return nil, fmt.Errorf("psql: prepare map photos insert: %v", err)
	}

	if mapGet, err = db.Conn.Prepare(getMapStatement); err != nil {
		return nil, fmt.Errorf("psql: prepare map get: %v", err)
	}

	if mapPhotosGet, err = db.Conn.Prepare(getMapPhotosStatement); err != nil {
		return nil, fmt.Errorf("psql: prepare map photos get: %v", err)
	}

	db.users.insert = userInsert
	db.users.get = userGet
	db.users.getByLogin = userGetByLogin

	db.maps.insert = mapInsert
	db.maps.insertPhotos = mapPhotosInsert
	db.maps.get = mapGet
	db.maps.getPhotos = mapPhotosGet

	return db, nil
}

const insertUserStatement = `
	INSERT INTO users (
		created_at, updated_at, email, username, password
	) VALUES ($1, $2, $3, $4, $5) RETURNING id`

func (db *PsqlDB) AddUser(u *schema.CreateUser) (id int64, err error) {
	now := time.Now().UnixNano()

	var insertedID int64

	err = QueryRow(db.users.insert, now, now, u.Email, u.Username, u.Password).
		Scan(&insertedID)

	if err != nil || insertedID == 0 {
		if pgerr, ok := err.(*pq.Error); ok {
			if pgerr.Code == "23505" {
				log.WithFields(log.Fields{
					"code":   pgerr.Code,
					"const":  pgerr.Constraint,
					"detail": pgerr.Detail,
				}).Error("psqlError")

				if pgerr.Constraint == "users_email_key" {
					return 0, errors.New("Email is not unique")
				} else if pgerr.Constraint == "users_username_key" {
					return 0, errors.New("Username is not unique")
				}
			}
		} else {
			log.Error(err)
		}

		return 0, errors.New("Could not create user")
	}

	return insertedID, nil
}

const getUserStatement = "SELECT * from users WHERE id = $1"

func (db *PsqlDB) GetUser(id int64) (*schema.UserView, error) {
	user, err := scanUser(db.users.get.QueryRow(id))

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("psql: could not find user with id %d", id)
	}

	if err != nil {
		return nil, fmt.Errorf("psql: could not get user: %v", err)
	}

	// convert user into what we want to return with API
	safeUserView := &schema.UserView{
		ID:       user.ID,
		Email:    user.Email,
		Username: user.Username,
	}

	return safeUserView, nil
}

const getUserByLoginStatement = "SELECT * from users WHERE email = $1"

func (db *PsqlDB) GetUserByLogin(email string) (*schema.User, error) {
	user, err := scanUser(db.users.getByLogin.QueryRow(email))

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("psql: could not find user with email %s", email)
	}

	if err != nil {
		return nil, fmt.Errorf("psql: could not get user: %v", err)
	}

	return user, nil
}

const insertMapStatement = `INSERT into maps (
		created_at, updated_at, name, description, download_code, type, views, user_id
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`
const insertMapPhotosStatement = `INSERT into map_photos (
		map_id, url
	) VALUES ($1, $2) RETURNING id`

func (db *PsqlDB) AddMap(m *schema.Map) (int64, error) {
	now := time.Now().UnixNano()

	var insertedID int64

	err := QueryRow(db.maps.insert,
		now, now, m.Name, m.Description, m.DownloadCode, m.Type, 0, m.UserID).
		Scan(&insertedID)

	if err != nil {
		return 0, fmt.Errorf("psql: could not create map: %v", err)
	}

	if len(m.Photos) > 0 {
		var insertedPhotoID int64

		for _, url := range m.Photos {
			err = QueryRow(db.maps.insertPhotos, insertedID, url).Scan(&insertedPhotoID)

			if err != nil {
				log.WithFields(log.Fields{
					"error": err,
					"url":   url,
					"mapId": insertedID,
				}).Error("Could not save photo")
			} else {
				log.WithFields(log.Fields{
					"photoId": insertedPhotoID,
				}).Info("Photo saved")
			}
		}
	}

	return insertedID, nil
}

const getMapStatement = "SELECT * from maps WHERE id = $1"
const getMapPhotosStatement = "SELECT url from map_photos WHERE map_id = $1"

func (db *PsqlDB) GetMap(id int64) (*schema.Map, error) {
	m, err := scanMap(db.maps.get.QueryRow(id))

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("psql: could not find map with id %d", id)
	}

	if err != nil {
		return nil, fmt.Errorf("psql: could not get map: %v", err)
	}

	rows, err := db.maps.getPhotos.Query(id)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var photos []string

	for rows.Next() {
		photo, err := scanPhoto(rows)

		if err != nil {
			return nil, fmt.Errorf("psql: could not read row: %v", err)
		}

		photos = append(photos, photo)
	}

	m.Photos = photos

	return m, nil
}

type rowScanner interface {
	Scan(dest ...interface{}) error
}

func scanUser(s rowScanner) (*schema.User, error) {
	var (
		id        int64
		createdAt int
		updatedAt int
		email     string
		username  string
		password  string
	)

	if err := s.Scan(&id, &createdAt, &updatedAt, &email, &username, &password); err != nil {
		return nil, err
	}

	user := &schema.User{
		ID:        id,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		Email:     email,
		Username:  username,
		Password:  password,
	}

	return user, nil
}

func scanMap(s rowScanner) (*schema.Map, error) {
	var (
		id           int64
		createdAt    int
		updatedAt    int
		name         string
		description  string
		downloadCode string
		Type         string
		userID       int64
		views        int
	)

	if err := s.Scan(&id, &createdAt, &updatedAt, &name, &description, &downloadCode, &Type, &userID, &views); err != nil {
		return nil, err
	}

	m := &schema.Map{
		ID:           id,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
		Name:         name,
		Description:  description,
		DownloadCode: downloadCode,
		Type:         Type,
		UserID:       userID,
		Views:        views,
	}

	return m, nil
}

func scanPhoto(s rowScanner) (string, error) {
	var (
		url string
	)

	if err := s.Scan(&url); err != nil {
		return "", err
	}

	return url, nil
}

func QueryRow(stmt *sql.Stmt, args ...interface{}) *sql.Row {
	r := stmt.QueryRow(args...)

	return r
}

func Query(stmt *sql.Stmt, args ...interface{}) (*sql.Rows, error) {
	r, err := stmt.Query(args...)

	if err != nil {
		log.Error(err)
		return nil, err
	}

	return r, err
}
