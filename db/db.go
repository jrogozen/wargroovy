package db

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jrogozen/wargroovy/schema"
	"github.com/lib/pq"
	"github.com/rs/xid"
	log "github.com/sirupsen/logrus"
	"strings"
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
	update     *sql.Stmt
	// delete     *sql.Stmt
}

type mapSQL struct {
	insert       *sql.Stmt
	insertPhotos *sql.Stmt
	get          *sql.Stmt
	getBySlug    *sql.Stmt
	getPhotos    *sql.Stmt
	deletePhoto  *sql.Stmt
	listBy       *sql.Stmt
	update       *sql.Stmt
	delete       *sql.Stmt
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
	var userUpdate *sql.Stmt

	var mapInsert *sql.Stmt
	var mapPhotosInsert *sql.Stmt
	var mapGet *sql.Stmt
	var mapGetBySlug *sql.Stmt
	var mapListBy *sql.Stmt
	var mapUpdate *sql.Stmt
	var mapPhotoDelete *sql.Stmt
	var mapDelete *sql.Stmt

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

	if mapListBy, err = db.Conn.Prepare(listByMapStatement); err != nil {
		return nil, fmt.Errorf("psql: listBy map statement: %v", err)
	}

	if mapUpdate, err = db.Conn.Prepare(updateMapStatement); err != nil {
		return nil, fmt.Errorf("psql: update map: %v", err)
	}

	if mapGetBySlug, err = db.Conn.Prepare(getMapBySlugStatement); err != nil {
		return nil, fmt.Errorf("psql: prepare map get by slug: %v", err)
	}

	if mapPhotoDelete, err = db.Conn.Prepare(deleteMapPhotoStatement); err != nil {
		return nil, fmt.Errorf("psql: prepare map delete photo: %v", err)
	}

	if mapDelete, err = db.Conn.Prepare(deleteMapStatement); err != nil {
		return nil, fmt.Errorf("psql: prepare delete map: %v", err)
	}

	if userUpdate, err = db.Conn.Prepare(updateUserStatement); err != nil {
		return nil, fmt.Errorf("psql: prepare update user: %v", err)
	}

	db.users.insert = userInsert
	db.users.get = userGet
	db.users.getByLogin = userGetByLogin
	db.users.update = userUpdate

	db.maps.insert = mapInsert
	db.maps.insertPhotos = mapPhotosInsert
	db.maps.get = mapGet
	db.maps.getBySlug = mapGetBySlug
	db.maps.listBy = mapListBy
	db.maps.update = mapUpdate
	db.maps.deletePhoto = mapPhotoDelete
	db.maps.delete = mapDelete

	return db, nil
}

const insertUserStatement = `
	INSERT INTO users (
		created_at, updated_at, email, username, password
	) VALUES ($1, $2, $3, $4, $5) RETURNING id`

func (db *PsqlDB) AddUser(u *schema.CreateUser) (id int64, err error) {
	now := time.Now().UnixNano()

	var insertedID int64

	if u.Email == "" || u.Username == "" {
		err = QueryRow(db.users.insert, now, now, nil, nil, u.Password).
			Scan(&insertedID)
	} else {
		err = QueryRow(db.users.insert, now, now, u.Email, u.Username, u.Password).
			Scan(&insertedID)
	}

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

func (db *PsqlDB) GetUser(id int64) (*schema.User, error) {
	user, err := scanUser(db.users.get.QueryRow(id))

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("psql: could not find user with id %d", id)
	}

	if err != nil {
		return nil, fmt.Errorf("psql: could not get user: %v", err)
	}

	return user, nil
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

const updateUserStatement = `UPDATE users
	SET updated_at = $1, email = $2, username = $3, password = $4
	WHERE id = $5
	RETURNING id`

func (db *PsqlDB) UpdateUser(u *schema.User) (int64, error) {
	if u.ID == 0 {
		return 0, errors.New("psql: cannot update user with unassigned ID")
	}

	now := time.Now().UnixNano()

	var updatedID int64

	err := QueryRow(db.users.update, now, u.Email, u.Username, u.Password, u.ID).
		Scan(&updatedID)

	if err != nil {
		return 0, err
	}

	return updatedID, nil
}

const insertMapStatement = `INSERT into maps (
		created_at, updated_at, name, description, download_code, type, views, slug, user_id
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`
const insertMapPhotosStatement = `INSERT into map_photos (
		map_id, url
	) VALUES ($1, $2) RETURNING id`

func (db *PsqlDB) AddMap(m *schema.Map) (int64, error) {
	now := time.Now().UnixNano()

	// nicely formatted unique url
	slug := strings.Replace(m.Name, " ", "-", -1) + "-" + xid.New().String()

	var insertedID int64

	err := QueryRow(db.maps.insert,
		now, now, m.Name, m.Description, m.DownloadCode, m.Type, 0, slug, m.UserID).
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

const getMapStatement = `SELECT m.id, m.created_at, m.updated_at, m.name, m.description, m.download_code, m.type, m.user_id, m.views, m.slug, photos
FROM maps m
left join (
	select map_id, string_agg(url, ',') AS photos
	from map_photos
	GROUP BY map_id
) p ON m.id = map_id
WHERE id = $1
`

func (db *PsqlDB) GetMap(id int64) (*schema.Map, error) {
	m, err := scanMap(db.maps.get.QueryRow(id))

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("psql: could not find map with id %d", id)
	}

	if err != nil {
		return nil, fmt.Errorf("psql: could not get map: %v", err)
	}

	return m, nil
}

const getMapBySlugStatement = `SELECT m.id, m.created_at, m.updated_at, m.name, m.description, m.download_code, m.type, m.user_id, m.views, m.slug, photos
FROM maps m
left join (
	select map_id, string_agg(url, ',') AS photos
	from map_photos
	GROUP BY map_id
) p ON m.id = map_id
WHERE slug = $1
`

func (db *PsqlDB) GetMapBySlug(slug string) (*schema.Map, error) {
	m, err := scanMap(db.maps.getBySlug.QueryRow(slug))

	log.Info("looking for map")
	log.Info(slug)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("psql: could not find map with slug %s", slug)
	}

	if err != nil {
		return nil, fmt.Errorf("psql: could not get map: %v", err)
	}

	return m, nil
}

const listByMapStatement = `select m.id, m.created_at, m.updated_at, m.name, m.description, m.download_code, m.type, m.user_id, m.views, m.slug, photos
	from maps m
	left join (
		select map_id, string_agg(url, ',') AS photos
		from map_photos
		GROUP BY map_id
	) p ON m.id = map_id
	order by $1
	limit $2
	offset $3
`

func (db *PsqlDB) ListByMap(options *schema.SortOptions) ([]*schema.Map, error) {
	rows, err := db.maps.listBy.Query(options.OrderBy, options.Limit, options.Offset)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var maps []*schema.Map

	for rows.Next() {
		m, err := scanMap(rows)

		if err != nil {
			return nil, fmt.Errorf("psql: could not read row: %v", err)
		}

		log.WithFields(log.Fields{
			"map": m,
		}).Info("appending map")

		maps = append(maps, m)
	}

	return maps, nil
}

const updateMapStatement = `UPDATE maps
	SET updated_at = $1, name = $2, description = $3, download_code = $4, type = $5, slug = $6
	WHERE id = $7
	RETURNING id
`

func (db *PsqlDB) UpdateMap(m *schema.Map) (int64, error) {
	if m.ID == 0 {
		return 0, errors.New("psql: cannot update map with unassigned ID")
	}

	now := time.Now().UnixNano()

	var updatedID int64

	err := QueryRow(db.maps.update, now, m.Name, m.Description, m.DownloadCode, m.Type, m.Slug, m.ID).
		Scan(&updatedID)

	if err != nil {
		return 0, err
	}

	return updatedID, nil
}

const deleteMapPhotoStatement = "DELETE FROM map_photos WHERE url = $1 AND map_id = $2"

func (db *PsqlDB) DeleteMapPhoto(mapID int64, url string) (int64, error) {
	if mapID == 0 {
		return 0, errors.New("psql: cannot delete photos from unassigned map ID")
	}

	r, err := db.maps.deletePhoto.Exec(url, mapID)

	if err != nil {
		return 0, errors.New("psql: could not delete photo")
	}

	rowsAffected, err := r.RowsAffected()

	return rowsAffected, nil
}

const deleteMapStatement = "DELETE from maps WHERE id = $1"

func (db *PsqlDB) DeleteMap(mapID int64) (int64, error) {
	if mapID == 0 {
		return 0, errors.New("psql: cannot delete map with unassigned map ID")
	}

	r, err := db.maps.delete.Exec(mapID)

	if err != nil {
		return 0, errors.New("psql: could not delete map")
	}

	rowsAffected, err := r.RowsAffected()

	return rowsAffected, nil
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
		slug         string
		photos       sql.NullString
	)

	if err := s.Scan(
		&id,
		&createdAt,
		&updatedAt,
		&name,
		&description,
		&downloadCode,
		&Type,
		&userID,
		&views,
		&slug,
		&photos,
	); err != nil {
		return nil, err
	}

	var photosArray []string

	if photos.Valid {
		photosArray = strings.Split(photos.String, ",")
	} else {
		photosArray = make([]string, 0)
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
		Slug:         slug,
		Photos:       photosArray,
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
