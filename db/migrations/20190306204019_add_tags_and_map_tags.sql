
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
create table if not exists tags (
    id serial primary key,
    name text unique
);

create table if not exists map_tags (
    id serial primary key,
    map_id integer references maps(id) on delete cascade,
    tag_id integer references tags(id) on delete cascade
);


-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

drop table if exists tags;
drop table if exists map_tags;