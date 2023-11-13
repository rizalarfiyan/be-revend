-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TYPE role AS ENUM (
    'admin',
    'guest'
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TYPE role;
-- +goose StatementEnd
