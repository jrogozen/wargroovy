
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
create table if not exists map_tags (
    map_id integer references maps(id) on delete cascade,
    tag_name text not null,
    primary key (map_id, tag_name)
);


-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

drop table if exists map_tags cascade;