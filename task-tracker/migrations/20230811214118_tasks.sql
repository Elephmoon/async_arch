-- +goose Up
-- +goose StatementBegin
CREATE TABLE tasks (
    id uuid not null,
    name text not null,
    description text not null,
    user_id uuid not null,
    price bigint not null default 0,
    fee bigint not null default 0,
    created_at timestamptz not null default now(),
    closed_at timestamptz
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
