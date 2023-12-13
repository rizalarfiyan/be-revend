-- name: GetAllDevice :many
SELECT * FROM device;

-- name: GetAllNameDevice :many
SELECT id, name FROM device;

-- name: CountAllDevice :one
SELECT count(*) FROM device;
