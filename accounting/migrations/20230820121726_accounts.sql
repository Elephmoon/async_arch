-- +goose Up
-- +goose StatementBegin
CREATE TABLE accounts (
    number text not null primary key,
    public_id uuid not null,
    user_id uuid not null,
    user_public_id uuid not null,
    balance bigint not null default 0,
    created_at timestamptz not null default now(),
    unique (user_id, public_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
