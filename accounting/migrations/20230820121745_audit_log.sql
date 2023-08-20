-- +goose Up
-- +goose StatementBegin
CREATE TABLE audit_log (
    id uuid not null default gen_random_uuid() primary key,
    event text not null,
    account_number text not null,
    task_id uuid,
    change_amount bigint not null,
    idempotency_key text not null unique,
    created_at timestamptz not null default now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
