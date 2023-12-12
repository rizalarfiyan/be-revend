-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS history (
    id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
    user_id uuid,
    device_id uuid,
    success int NOT NULL DEFAULT 0,
    failed int NOT NULL DEFAULT 0,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (device_id) REFERENCES device (id)
);
CREATE INDEX idx_history ON history(id, user_id, device_id, success, failed);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX idx_history;
DROP TABLE history;
-- +goose StatementEnd
