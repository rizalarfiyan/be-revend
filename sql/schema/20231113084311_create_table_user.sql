-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
    id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
    first_name varchar(100) NOT NULL,
    last_name varchar(100),
    phone_number varchar(20) UNIQUE NOT NULL,
    google_id varchar(30) UNIQUE NOT NULL,
    role role NOT NULL DEFAULT ('guest'),
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp
);
CREATE INDEX idx_users ON users(id, phone_number, google_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX idx_users;
DROP TABLE users;
-- +goose StatementEnd
