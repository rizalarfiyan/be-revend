-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
    id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
    first_name varchar(100) NOT NULL,
    last_name varchar(100),
    phone_number varchar(20) UNIQUE NOT NULL,
    google_id varchar(30),
    identity varchar(50) UNIQUE NOT NULL,
    role role NOT NULL DEFAULT ('guest'),
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp,
    deleted_by uuid,
    deleted_at timestamp,
    FOREIGN KEY (deleted_by) REFERENCES users (id)
);
CREATE UNIQUE INDEX users_google_id_key ON users (google_id) WHERE google_id IS NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
