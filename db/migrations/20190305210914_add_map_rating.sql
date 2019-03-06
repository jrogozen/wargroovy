
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

create table if not exists map_ratings (
    user_id integer references users(id) not null,
    map_id integer references maps(id) not null,
    rating integer not null,
    primary key (map_id, user_id)
);


-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

drop table if exists map_ratings;