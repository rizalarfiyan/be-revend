-- name: GetAllHistory :many
SELECT h.*, u.first_name, u.last_name, d.name as device_name FROM history h
JOIN users u ON u.id = h.user_id
JOIN device d ON d.id = h.device_id;

-- name: GetAllHistoryStatistic :many
SELECT MIN(created_at)::date as date, SUM(success) AS success, SUM(failed) AS failed
FROM history;

-- name: GetAllHistoryTopPerformance :many
SELECT h.user_id AS user_id, u.first_name, u.last_name, u.phone_number, SUM(h.success) AS success, SUM(h.failed) AS failed
FROM history h
JOIN users u ON u.id = h.user_id
WHERE h.created_at BETWEEN sqlc.arg('start_date') AND sqlc.arg('end_date')
GROUP BY h.user_id, u.first_name, u.last_name, u.phone_number
ORDER BY success DESC
LIMIT sqlc.arg('limit');

-- name: CountAllHistory :one
SELECT count(h.*) FROM history h
JOIN users u ON u.id = h.user_id
JOIN device d ON d.id = h.device_id;

-- name: CreateHistory :exec
INSERT INTO history (user_id, device_id, success, failed)
VALUES ($1, $2, $3, $4);
