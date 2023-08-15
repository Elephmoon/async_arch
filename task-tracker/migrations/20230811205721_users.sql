-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id uuid not null primary key,
    name text not null unique,
    role text not null,
    created_at timestamptz not null,
    synced_at timestamptz not null,
    deleted_at timestamptz
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
