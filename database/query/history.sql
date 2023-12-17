-- name: GetAllHistory :many
SELECT h.*, u.first_name, u.last_name, d.name as device_name FROM history h
JOIN users u ON u.id = h.user_id
JOIN device d ON d.id = h.device_id;

-- name: CountAllHistory :one
SELECT count(h.*) FROM history h
JOIN users u ON u.id = h.user_id
JOIN device d ON d.id = h.device_id;

-- name: CreateHistory :exec
INSERT INTO history (user_id, device_id, success, failed)
VALUES ($1, $2, $3, $4);
