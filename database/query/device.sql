-- name: GetAllDevice :many
SELECT * FROM device;

-- name: GetAllNameDevice :many
SELECT id, name FROM device;

-- name: CountAllDevice :one
SELECT count(*) FROM device;

-- name: CreateDevice :exec
INSERT INTO device (token, name, location)
VALUES ($1, $2, $3);

-- name: UpdateDevice :exec
UPDATE device
SET name = $1, location = $2
WHERE id = $3;

-- name: ToggleDeleteDevice :exec
UPDATE device SET
deleted_by = CASE WHEN deleted_by IS NULL THEN sqlc.narg('deleted_by')::UUID ELSE NULL END,
deleted_at = CASE WHEN deleted_at IS NULL THEN CURRENT_TIMESTAMP ELSE NULL
END WHERE id = sqlc.narg('id');
