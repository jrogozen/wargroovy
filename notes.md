## wargroovy

### handlers

#### /auth
    [X] - POST - /auth/login

#### /user
    [x] - POST /user
    [X] - GET /user/{userId}

#### /campaign
    [x] - POST /campaign
    [x] - POST /campign/map
    [x] - GET /campaign/{campaignId}
    [] - GET /campaign/list // list of compaigns sortedBy?
    [] - PUT /campaign/{campaignId}
    [] - PUT /campaign/{campaignId}/map/{mapId} // necessary?


### relationships

users have campaigns
campaigns have one user

campaigns have maps
map has one campaign


### tables

#### campaign
    - id int
    - name string
    - description string
    - thumb_photo_url string
    - large_photo_url string
    - single_map_campaign boolean

#### map
    - id int
    - name string
    - description string
    - thumb_photo_url string
    - large_photo_url string
    - views int
    - download_code string

#### campaign_maps
    - id int
    - campaign_id int
    - map_id int

#### user_campaigns
    - id int
    - user_id int
    - campaign_id


### notes

on frontend, user is asked if they'd like to save a campaign or a single map
    - single map
        - create campaign with blank info + single_map_campaign TRUE
        - create map and attach to campaign
    - campaign
        - create campaign with info + single_map_campaign false
        - prompt for map info

all maps must be attached to a campaign. UI is different for single vs multi-map campaigns

can think about separating http handlers and database actions into separate packages