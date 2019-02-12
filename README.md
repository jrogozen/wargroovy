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
- [] set version of 3rd party packages
- [] set up error monitoring (datadog, ...)
- [] set up google sdk locally
- [] set up gcloud project, enable google app engine (GAE)
- [] deploy to GAE (docker (?) vs...)
- [] connect GAE and domain
- [] set up SSL
- [] set up postgres on google cloud
- [] set up env variable injection for GAE
- [] continous integration (circle CI, drone, ...)
- [] CI github hook (master push -> deploy)

### MVP todo
- [x] either remove chi/jwtauth or jwt/go
- [x] split out handler code into multiple files per package (create, update, etc...)
- [x] simplify nested conditionals in web response code (return a response?)
- [x] set up level based logging. maybe [logrus](https://github.com/Sirupsen/logrus)
- [x] jwt middleware should return json instead of 40x + text
- [] finish protecting user actions with jwt

### future todo
- [] replace gorm with raw sql queries
- [] discord w/ deploy / error integrations
- [] public issue repo on github
- [] ~~set up db migrations. maybe [sql-migrate](https://github.com/rubenv/sql-migrate)~~ temp use gorm.AutoMigrate