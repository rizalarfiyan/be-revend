-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS device (
    id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
    token varchar(64) NOT NULL,
    name varchar(50) NOT NULL,
    location varchar(150) NOT NULL,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp
);
CREATE INDEX idx_device ON device(id, token);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX idx_device;
DROP TABLE device;
-- +goose StatementEnd
