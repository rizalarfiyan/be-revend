-- name: GetUserById :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByPhoneNumber :one
SELECT * FROM users
WHERE phone_number = $1 LIMIT 1;

-- name: GetUserByGoogleId :one
SELECT * FROM users
WHERE google_id = $1 LIMIT 1;

-- name: GetUserByIdentity :one
SELECT * FROM users
WHERE identity = $1 LIMIT 1;

-- name: GetUserByGoogleIdOrPhoneNumber :one
SELECT * FROM users
WHERE google_id = $1 OR phone_number = $2 LIMIT 1;

-- name: GetAllUsers :many
SELECT * FROM users;

-- name: GetAllNameUsers :many
SELECT id, first_name, last_name FROM users;

-- name: CountAllUsers :one
SELECT count(*) FROM users;

-- name: CreateUser :exec
INSERT INTO users (first_name, last_name, phone_number, google_id, identity)
VALUES ($1, $2, $3, $4, $5);

-- name: ToggleDeleteUser :exec
UPDATE users SET
deleted_by = CASE WHEN deleted_by IS NULL THEN sqlc.narg('deleted_by')::UUID ELSE NULL END,
deleted_at = CASE WHEN deleted_at IS NULL THEN CURRENT_TIMESTAMP ELSE NULL
END WHERE id = sqlc.narg('id');
