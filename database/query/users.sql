-- name: GetUserByPhoneNumber :one
SELECT * FROM users
WHERE phone_number = $1 LIMIT 1;

-- name: GetUserByGoogleId :one
SELECT * FROM users
WHERE google_id = $1 LIMIT 1;
