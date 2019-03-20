## wargroovy

### handlers

#### /auth
- [X] - POST - /auth/login

#### /user
- [x] - POST /user
- [X] - GET /user (uses jwt for auth)

#### /map
- [x] - POST /map
- [x] - POST /map/{mapId}/rate
- [x] - GET /map/{mapId}
- [x] - GET /map/bySlug/{slug}
- [X] - GET /map/list - list of compaigns with queryParam
- [X] - PUT /map/{mapId}
- [X] - DELETE /map/{mapId}/photo

#### /photo
- [x] - POST /photo



### relationships
users have maps
maps have one user

maps have many photos
photos have one map


### deployment todo
- [x] set version of 3rd party packages
- [ ] set up error monitoring (datadog, ...)
- [x] set up google sdk locally
- [x] set up gcloud project, enable google app engine (GAE)
- [x] deploy to GAE
- [x] ~~set up SSL~~
- [x] set up postgres on google cloud
- [x] set up env variable injection for GAE
- [x] CI github hook (master push -> deploy)
- [x] best way to access cloud sql db. `gcloud sql connect wargroovy-production -u postgres`

### MVP todo
- [x] either remove chi/jwtauth or jwt/go
- [x] split out handler code into multiple files per package (create, update, etc...)
- [x] simplify nested conditionals in web response code (return a response?)
- [x] set up level based logging. maybe [logrus](https://github.com/Sirupsen/logrus)
- [x] jwt middleware should return json instead of 40x + text
- [x] finish protecting user actions with jwt
- [x] clean up models/controller interactions. models should return fully completed response :)
- [x] combine maps + campaigns
- [x] add to campaign: download code, type
- [x] do not return user email
- [x] add to user: name
- [x] change photos to separate table
- [x] upload photos route
- [x] change user/get to work based off of jwt. return not logged in without it. prevents people from scraping for users
- [x] add slug field to maps. generate this on the server for use in url
- [x] add finding map by slug
- [x] can create user without username AND email (but fail if only one supplied not two)
- [x] delete map photos
- [x] delete map
- [x] edit user
- [x] api return correct status codes (20x, 40x)
- [x] jwt set as cookie in response
- [x] increment map view
- [x] make work with draft-js data structure
- [x] filter by type for map list
- [x] return author for map info (userId -> username)
- [ ] log a billion times more

### ratings feature
- [x] add post rating api
- [x] get map apis should return a user's vote on a map
- [x] sort map list by rating
- [x] list map apis should return aggregate score per map
- [ ] fix users should be able to override previous vote

### tags feature
- [x] /map/tags get most popular tags
- [x] map_tags for map_id, tag_name
- [x] return tags per map, map list
- [x] map list filter by tags

### map_comments feature
comments can nest max 2 levels (main thread, reply)
- [ ] post comment
- [ ] edit comment
- [ ] mapById, mapBySlug includes comments

### future todo
- [x] replace gorm with raw sql queries [example](https://github.com/GoogleCloudPlatform/golang-samples/blob/master/appengine/go11x/cloudsql/cloudsql.go)
- [x] set up db migrations. using goose
- [ ] discord w/ deploy / error integrations
- [ ] public issue repo on github
- [ ] ML model to flag images that are not associated with wargroove
- [ ] add map comments
- [ ] jwt auth the photo upload api (requires username/email-less user creation)
- [ ] delete user

### punt
- ~~[ ] image pipeline to create thumbnails and upload images to gloud bucket~~
- [ ] ~~connect GAE and domain~~ probably wait on this. frontend and backend prob separate services
- [ ] ~~continous integration (circle CI, drone, ...)~~


### helpful commands

#### connect to production sql
`PGDATABASE=production gcloud sql connect wargroovy-production -u postgres`

#### deploy app
`gcloud app deploy`

#### deploy gcloud functions
`cd functions && gcloud functions deploy functions --entry-point RedditBot --runtime go111 --set-env-vars WARGROOVY_API=https://api.wargroovy.com,WARGROOVY_WEB_API=https://wargroovy.com,USER_TOKEN=`

### local dev
`cd web && gin -p 4000 -a 8080 -t ../ -d . run main.go`

runs the proxy on port 3000 (main app on 8080)