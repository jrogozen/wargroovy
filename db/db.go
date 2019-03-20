package db

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jrogozen/wargroovy/schema"
	// u "wargroovy/utils"
	"github.com/lib/pq"
	"github.com/rs/xid"
	log "github.com/sirupsen/logrus"
	"regexp"
	"strings"
	"time"
)

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
	insert            *sql.Stmt
	insertPhotos      *sql.Stmt
	get               *sql.Stmt
	getBySlug         *sql.Stmt
	getByDownloadCode *sql.Stmt
	getPhotos         *sql.Stmt
	deletePhoto       *sql.Stmt
	listBy            *sql.Stmt
	update            *sql.Stmt
	delete            *sql.Stmt
	view              *sql.Stmt
	rate              *sql.Stmt
	getUserRating     *sql.Stmt
	insertTag         *sql.Stmt
}

type PsqlDB struct {
	Conn *sql.DB

	users userSQL
	maps  mapSQL
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

	return db, nil
}

func PrepareDB(db *PsqlDB) (*PsqlDB, error) {
	var userInsert *sql.Stmt
	var userGet *sql.Stmt
	var userGetByLogin *sql.Stmt
	var userUpdate *sql.Stmt

	var mapInsert *sql.Stmt
	var mapPhotosInsert *sql.Stmt
	var mapGet *sql.Stmt
	var mapGetBySlug *sql.Stmt
	var mapUpdate *sql.Stmt
	var mapPhotoDelete *sql.Stmt
	var mapDelete *sql.Stmt
	var mapView *sql.Stmt
	var mapRate *sql.Stmt
	var mapUserRating *sql.Stmt
	var mapInsertTag *sql.Stmt
	var mapGetByDownloadCdode *sql.Stmt

	var err error

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

	if mapView, err = db.Conn.Prepare(incrementMapViewStatement); err != nil {
		return nil, fmt.Errorf("psql: prepare map view: %v", err)
	}

	if mapRate, err = db.Conn.Prepare(rateMapStatement); err != nil {
		return nil, fmt.Errorf("psql: prepare rate map error: %v", err)
	}

	if mapUserRating, err = db.Conn.Prepare(getUserRatingForMapStatement); err != nil {
		return nil, fmt.Errorf("psql: prepare map user rating: %v", err)
	}

	if mapInsertTag, err = db.Conn.Prepare(insertMapTagStatement); err != nil {
		return nil, fmt.Errorf("psql: prepare insert map tag: %v", err)
	}

	if mapGetByDownloadCdode, err = db.Conn.Prepare(getMapByDownloadCodeStatement); err != nil {
		return nil, fmt.Errorf("psql: prepare get map by download code: %v", err)
	}

	db.users.insert = userInsert
	db.users.get = userGet
	db.users.getByLogin = userGetByLogin
	db.users.update = userUpdate

	db.maps.insert = mapInsert
	db.maps.insertPhotos = mapPhotosInsert
	db.maps.get = mapGet
	db.maps.getBySlug = mapGetBySlug
	db.maps.getByDownloadCode = mapGetByDownloadCdode
	db.maps.update = mapUpdate
	db.maps.deletePhoto = mapPhotoDelete
	db.maps.delete = mapDelete
	db.maps.view = mapView
	db.maps.rate = mapRate
	db.maps.getUserRating = mapUserRating
	db.maps.insertTag = mapInsertTag

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
	reg, _ := regexp.Compile("[^a-zA-Z0-9]+")

	processedName := reg.ReplaceAllString(m.Name, "-")[0:10]

	// need to generate a new slug
	slug := strings.ToLower(processedName) + "-" + xid.New().String()

	var insertedID int64
	var err error

	//TODO: more idomatic way of doing this
	if len(m.Description) == 0 {
		err = QueryRow(db.maps.insert, now, now, m.Name, nil, m.DownloadCode, m.Type, 0, slug, m.UserID).
			Scan(&insertedID)
	} else {
		err = QueryRow(db.maps.insert, now, now, m.Name, m.Description, m.DownloadCode, m.Type, 0, slug, m.UserID).
			Scan(&insertedID)
	}

	if err != nil {
		return 0, fmt.Errorf("psql: could not create map: %v", err)
	}

	if len(m.Photos) > 0 {
		var insertedPhotoID int64

		for _, url := range m.Photos {
			err = QueryRow(db.maps.insertPhotos, insertedID, url).
				Scan(&insertedPhotoID)

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
	if len(m.Tags) > 0 {
		for _, tag := range m.Tags {
			tagToInsert := strings.ToLower(tag)

			var insertedTag string

			err = QueryRow(db.maps.insertTag, insertedID, tagToInsert).
				Scan(&insertedTag)

			if err != nil {
				log.WithFields(log.Fields{
					"error": err,
					"tag":   tagToInsert,
					"mapID": insertedID,
				}).Error("Could not save tag")
			} else {
				log.WithFields(log.Fields{
					"tag":   tagToInsert,
					"mapID": insertedID,
				}).Info("Tag saved")
			}
		}
	}

	return insertedID, nil
}

const getMapStatement = `SELECT m.id, m.created_at, m.updated_at, m.name, m.description, m.download_code, m.type, m.user_id, m.views, m.slug, photos, u.username, rating, tags
FROM maps m
left join (
	select map_id, string_agg(url, ',') AS photos
	from map_photos
	GROUP BY map_id
) p ON m.id = map_id
left join (
	select map_id, round(cast(sum(rating) as decimal) / (cast(count(rating) as decimal) * 2), 2) * 100 as rating
	from map_ratings
	group by map_id
) r on m.id = r.map_id
left join (
	select map_id, string_agg(tag_name, ',') as tags
	from map_tags
	group by map_id
) t on m.id = t.map_id
left join users u
on m.user_id = u.id
WHERE m.id = $1
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

const getMapBySlugStatement = `SELECT m.id, m.created_at, m.updated_at, m.name, m.description, m.download_code, m.type, m.user_id, m.views, m.slug, photos, u.username, rating, tags
FROM maps m
left join (
	select map_id, string_agg(url, ',') AS photos
	from map_photos
	GROUP BY map_id
) p ON m.id = map_id
left join (
	select map_id, round(cast(sum(rating) as decimal) / (cast(count(rating) as decimal) * 2), 2) * 100 as rating
	from map_ratings
	group by map_id
) r on m.id = r.map_id
left join (
	select map_id, string_agg(tag_name, ',') as tags
	from map_tags
	group by map_id
) t on m.id = t.map_id
left join users u
on m.user_id = u.id
WHERE slug = $1
`

func (db *PsqlDB) GetMapBySlug(slug string) (*schema.Map, error) {
	m, err := scanMap(db.maps.getBySlug.QueryRow(slug))

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("psql: could not find map with slug %s", slug)
	}

	if err != nil {
		return nil, fmt.Errorf("psql: could not get map: %v", err)
	}

	return m, nil
}

const getMapByDownloadCodeStatement = `SELECT m.id, m.created_at, m.updated_at, m.name, m.description, m.download_code, m.type, m.user_id, m.views, m.slug, photos, u.username, rating, tags
FROM maps m
left join (
	select map_id, string_agg(url, ',') AS photos
	from map_photos
	GROUP BY map_id
) p ON m.id = map_id
left join (
	select map_id, round(cast(sum(rating) as decimal) / (cast(count(rating) as decimal) * 2), 2) * 100 as rating
	from map_ratings
	group by map_id
) r on m.id = r.map_id
left join (
	select map_id, string_agg(tag_name, ',') as tags
	from map_tags
	group by map_id
) t on m.id = t.map_id
left join users u
on m.user_id = u.id
WHERE download_code = $1
`

func (db *PsqlDB) GetMapByDownloadCode(code string) (*schema.Map, error) {
	m, err := scanMap(db.maps.getByDownloadCode.QueryRow(code))

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("psql: could not find map with download code %s", code)
	}

	if err != nil {
		return nil, fmt.Errorf("psql: could not get map: %v", err)
	}

	return m, nil
}

const listByMapStatement = `select m.id, m.created_at, m.updated_at, m.name, m.description, m.download_code, m.type, m.user_id, m.views, m.slug, photos, u.username, rating, tags
from maps m
left join (
	select map_id, string_agg(url, ',') AS photos
	from map_photos
	GROUP BY map_id
) p ON m.id = map_id
left join (
	select map_id, round(cast(sum(rating) as decimal) / (cast(count(rating) as decimal) * 2), 2) * 100 as rating
	from map_ratings
	group by map_id
) r on m.id = r.map_id
left join (
	select map_id, string_agg(tag_name, ',') as tags
	from map_tags
	group by map_id
) t on m.id = t.map_id
left join users u
on m.user_id = u.id
%s
%s
order by %s
limit %d
offset %d
`

func (db *PsqlDB) ListByMap(options *schema.SortOptions) ([]*schema.Map, error) {
	// order by is dynamic and cannot be prepared
	qtext := fmt.Sprintf(listByMapStatement, options.Type, options.Tags, options.OrderBy, options.Limit, options.Offset)

	log.WithFields(log.Fields{
		"type":    options.Type,
		"orderBy": options.OrderBy,
		"limit":   options.Limit,
		"offset":  options.Offset,
		"tags":    options.Tags,
	}).Info("listing map with options")

	log.Info(qtext)
	rows, err := db.Conn.Query(qtext)

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

const incrementMapViewStatement = `UPDATE maps
SET views = views + 1
WHERE ID = $1
RETURNING id
`

func (db *PsqlDB) IncrementMapView(mapID int64) (int64, error) {
	if mapID == 0 {
		return 0, errors.New("psql: cannot update map with unassigned map ID")
	}

	var updatedID int64

	err := QueryRow(db.maps.view, mapID).
		Scan(&updatedID)

	if err != nil {
		return 0, err
	}

	return updatedID, nil
}

const rateMapStatement = `insert into map_ratings (
	map_id, user_id, rating
) values ($1, $2, $3) returning rating`

func (db *PsqlDB) RateMap(mapID int64, userID int64, rating int64) (int64, error) {
	if mapID == 0 {
		return 0, errors.New("psql: cannot update map with unassigned map ID")
	}

	if userID == 0 {
		return 0, errors.New("psql: cannot rate map with unassigned user ID")
	}

	var insertedRating int64

	// 2 is thumbs up, 1 is thumbs down
	if rating >= 2 {
		rating = 2
	} else {
		rating = 1
	}

	err := QueryRow(db.maps.rate, mapID, userID, rating).
		Scan(&insertedRating)

	if err != nil {
		return 0, err
	}

	return insertedRating, nil
}

const getUserRatingForMapStatement = `select rating from map_ratings where user_id = $1 and map_id = $2`

func (db *PsqlDB) GetMapUserRating(mapID int64, userID int64) (int64, error) {
	if mapID == 0 {
		return 0, errors.New("psql: cannot update map with unassigned map ID")
	}

	if userID == 0 {
		return 0, errors.New("psql: cannot rate map with unassigned user ID")
	}

	var rating int64

	err := QueryRow(db.maps.getUserRating, userID, mapID).
		Scan(&rating)

	if err != nil {
		return 0, err
	}

	return rating, nil
}

type rowScanner interface {
	Scan(dest ...interface{}) error
}

const insertMapTagStatement = `insert into map_tags (
	map_id, tag_name
) values ($1, $2) returning tag_name`

func (db *PsqlDB) AddMapTag(mapID int64, name string) (string, error) {
	if name == "" {
		return "", errors.New("psql: cannot create tag with blank name")
	}

	var insertedTag string

	err := QueryRow(db.maps.insertTag, mapID, name).
		Scan(&insertedTag)

	if err != nil {
		return "", err
	}

	return insertedTag, nil
}

const getMapTagsStatement = `select tag_name, COUNT(*)
from map_tags
group by tag_name
order by %s
limit %d
`

func (db *PsqlDB) GetMapListTags(orderBy string, limit int) ([]*schema.TagCount, error) {
	// order by is dynamic and cannot be prepared
	qtext := fmt.Sprintf(getMapTagsStatement, orderBy, limit)

	rows, err := db.Conn.Query(qtext)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var tags []*schema.TagCount

	for rows.Next() {
		t, err := scanTag(rows)

		if err != nil {
			return nil, fmt.Errorf("psql: could not read row: %v", err)
		}

		tags = append(tags, t)
	}

	return tags, nil
}

func scanUser(s rowScanner) (*schema.User, error) {
	var (
		id        int64
		createdAt int
		updatedAt int
		email     sql.NullString
		username  sql.NullString
		password  sql.NullString
	)

	if err := s.Scan(&id, &createdAt, &updatedAt, &email, &username, &password); err != nil {
		return nil, err
	}

	user := &schema.User{
		ID:        id,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		Email:     email.String,
		Username:  username.String,
		Password:  password.String,
	}

	return user, nil
}

func scanMap(s rowScanner) (*schema.Map, error) {
	var (
		id           int64
		createdAt    int
		updatedAt    int
		name         string
		description  schema.DescriptionMap
		downloadCode string
		Type         string
		userID       int64
		views        int
		slug         string
		photos       sql.NullString
		username     sql.NullString
		rating       sql.NullFloat64
		tags         sql.NullString
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
		&username,
		&rating,
		&tags,
	); err != nil {
		return nil, err
	}

	var photosArray []string
	var author string
	var ratingFloat float64
	var tagsArray []string

	if photos.Valid {
		photosArray = strings.Split(photos.String, ",")
	} else {
		photosArray = make([]string, 0)
	}

	if tags.Valid {
		tagsArray = strings.Split(tags.String, ",")
	} else {
		tagsArray = make([]string, 0)
	}

	if username.Valid {
		author = username.String
	} else {
		author = "anonymous"
	}

	if rating.Valid {
		ratingFloat = rating.Float64
	} else {
		ratingFloat = 0
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
		Author:       author,
		Rating:       ratingFloat,
		Tags:         tagsArray,
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

func scanTag(t rowScanner) (*schema.TagCount, error) {
	var (
		tagName string
		count   int
	)

	if err := t.Scan(
		&tagName,
		&count,
	); err != nil {
		return nil, err
	}

	tag := &schema.TagCount{
		Name:  tagName,
		Count: count,
	}

	return tag, nil
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
