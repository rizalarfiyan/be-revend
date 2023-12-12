-- name: GetAllDevice :many
SELECT * FROM device;

-- name: CountAllDevice :one
SELECT count(*) FROM device;
