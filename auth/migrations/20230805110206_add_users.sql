-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id uuid not null primary key,
    public_id uuid not null primary key,
    name text not null unique,
    password text not null,
    role text not null,
    created_at timestamptz not null default now(),
    deleted_at timestamptz
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
