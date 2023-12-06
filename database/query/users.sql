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

-- name: CreateUser :exec
INSERT INTO users (first_name, last_name, phone_number, google_id)
VALUES ($1, $2, $3, $4);
