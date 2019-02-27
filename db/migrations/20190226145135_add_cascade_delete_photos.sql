
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
alter table map_photos
drop constraint if exists map_photos_map_id_fkey;

alter table map_photos
add constraint map_photos_map_id_fkey
foreign key (map_id)
references maps (id)
on delete cascade;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

alter table map_photos
drop constraint if exists map_photos_map_id_fkey;

alter table map_photos
add constraint map_photos_map_id_fkey
foreign key (map_id)
references maps (id);