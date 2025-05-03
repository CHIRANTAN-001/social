-- +goose Up
-- +goose StatementBegin
ALTER TABLE followers
    RENAME COLUMN follower_id to followed_id;
ALTER TABLE followers
    RENAME COLUMN user_id to follower_id;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- +goose StatementEnd
