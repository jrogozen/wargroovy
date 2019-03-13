module github.com/jrogozen/wargroovy

require (
	bitbucket.org/liamstask/goose v0.0.0-20150115234039-8488cc47d90c // indirect
	cloud.google.com/go v0.37.0 // indirect
	github.com/go-chi/chi v4.0.2+incompatible
	github.com/go-chi/jwtauth v4.0.0+incompatible // indirect
	github.com/go-chi/render v1.0.1
	github.com/go-sql-driver/mysql v1.4.1 // indirect
	github.com/gofrs/uuid v3.2.0+incompatible // indirect
	github.com/joho/godotenv v1.3.0 // indirect
	github.com/jrogozen/wargroovy/db v0.0.0 // indirect
	github.com/jrogozen/wargroovy/handlers/auth v0.0.0
	github.com/jrogozen/wargroovy/handlers/maps v0.0.0
	github.com/jrogozen/wargroovy/handlers/photos v0.0.0
	github.com/jrogozen/wargroovy/handlers/user v0.0.0
	github.com/jrogozen/wargroovy/internal/config v0.0.0
	github.com/jrogozen/wargroovy/schema v0.0.0 // indirect
	github.com/jrogozen/wargroovy/utils v0.0.0 // indirect
	github.com/kylelemons/go-gypsy v0.0.0-20160905020020-08cad365cd28 // indirect
	github.com/lib/pq v1.0.0 // indirect
	github.com/mattn/go-sqlite3 v1.10.0 // indirect
	github.com/rs/xid v1.2.1 // indirect
	github.com/sirupsen/logrus v1.4.0
	github.com/ziutek/mymysql v1.5.4 // indirect
	google.golang.org/appengine v1.4.0
)

replace (
	github.com/jrogozen/wargroovy/db => ./src/db
	github.com/jrogozen/wargroovy/handlers/auth => ./src/handlers/auth
	github.com/jrogozen/wargroovy/handlers/maps => ./src/handlers/maps
	github.com/jrogozen/wargroovy/handlers/photos => ./src/handlers/photos
	github.com/jrogozen/wargroovy/handlers/user => ./src/handlers/user
	github.com/jrogozen/wargroovy/internal/config => ./src/internal/config
	github.com/jrogozen/wargroovy/schema => ./src/schema
	github.com/jrogozen/wargroovy/utils => ./src/utils
)
