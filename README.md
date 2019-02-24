## wargroovy

### handlers

#### /auth
- [X] - POST - /auth/login

#### /user
- [x] - POST /user
- [X] - GET /user/{userId}

#### /campaign
- [x] - POST /campaign
- [x] - POST /campign/map
- [x] - GET /campaign/{campaignId}
- [X] - GET /campaign/list - list of compaigns with queryParam
- [X] - PUT /campaign/{campaignId}
- [X] - PUT /campaign/{campaignId}/map/{mapId}


### relationships
users have campaigns
campaigns have one user

campaigns have maps
map has one campaign

### notes
on frontend, user is asked if they'd like to save a campaign or a single map
    - single map
        - create campaign with blank info + single_map_campaign TRUE
        - create map and attach to campaign
    - campaign
        - create campaign with info + single_map_campaign false
        - prompt for map info

all maps must be attached to a campaign. UI is different for single vs multi-map campaigns

### deployment todo
- [x] set version of 3rd party packages
- [ ] set up error monitoring (datadog, ...)
- [x] set up google sdk locally
- [x] set up gcloud project, enable google app engine (GAE)
- [x] deploy to GAE
- [ ] ~~connect GAE and domain~~ probably wait on this. frontend and backend prob separate services
- [x] ~~set up SSL~~
- [x] set up postgres on google cloud
- [x] set up env variable injection for GAE
- [ ] ~~continous integration (circle CI, drone, ...)~~
- [x] CI github hook (master push -> deploy)
- [x] best way to access cloud sql db. `gcloud sql connect wargroovy-production -u postgres`

### MVP todo
- [x] either remove chi/jwtauth or jwt/go
- [x] split out handler code into multiple files per package (create, update, etc...)
- [x] simplify nested conditionals in web response code (return a response?)
- [x] set up level based logging. maybe [logrus](https://github.com/Sirupsen/logrus)
- [x] jwt middleware should return json instead of 40x + text
- [x] finish protecting user actions with jwt
- [ ] clean up models/controller interactions. models should return fully completed response :)
- [ ] search db for campaigns based on title (relating to url slug)
- [x] combine maps + campaigns
- [x] add to campaign: download code, type
- [ ] encode descriptions, make work with draft-js data structure
- [x] do not return user email
- [x] add to user: name
- [ ] FILTER posts data
- [x] change photos to separate table
- [ ] upload photos api
- [ ] can create user without username AND email (but fail if only one supplied not two)

### future todo
- [ ] replace gorm with raw sql queries [example](https://github.com/GoogleCloudPlatform/golang-samples/blob/master/appengine/go11x/cloudsql/cloudsql.go)
- [ ] discord w/ deploy / error integrations
- [ ] public issue repo on github
- [ ] ~~set up db migrations. maybe [sql-migrate](https://github.com/rubenv/sql-migrate)~~ temp use gorm.AutoMigrate
- [ ] image pipeline to create thumbnails and upload images to gloud bucket
- [ ] ML model to flag images that are not associated with wargroove
- [ ] add map comments
- [ ] add map ratings

### helpful commands

#### connect to production sql
`gcloud sql connect wargroovy-production -u postgres`

#### deploy app
`gcloud app deploy`

### local dev
`cd web && gin -p 3000 -a 8080 -t ../ -d . run main.go`

runs the proxy on port 3000 (main app on 8080)