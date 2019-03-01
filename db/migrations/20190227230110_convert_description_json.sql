
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
alter table maps
drop if exists description;

alter table maps
add column description jsonb;


-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
alter table maps
drop if exists description;

alter table maps
add column description text;

