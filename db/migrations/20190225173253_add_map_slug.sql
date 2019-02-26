
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

ALTER TABLE maps
ADD COLUMN IF NOT EXISTS slug text;


-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

ALTER TABLE maps
DROP COLUMN IF EXISTS slug;