
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
-- Table Definition ----------------------------------------------

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    created_at BIGINT,
    updated_at BIGINT,
    email text UNIQUE,
    username text UNIQUE,
    password text
);

-- Table Definition ----------------------------------------------

CREATE TABLE IF NOT EXISTS maps (
    id SERIAL PRIMARY KEY,
    created_at BIGINT,
    updated_at BIGINT,
    name text,
    description text,
    download_code text,
    type text,
    user_id integer REFERENCES users(id),
    views integer DEFAULT 0
);

-- Table Definition ----------------------------------------------

CREATE TABLE IF NOT EXISTS map_photos (
    id SERIAL PRIMARY KEY,
    map_id integer REFERENCES maps(id),
    url text
);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE map_photos, maps, users;
