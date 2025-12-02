-- name: CreateCommand :one
INSERT INTO commands (name, command, description, parameters)
VALUES (?, ?, ?, ?)
RETURNING *;

-- name: GetCommandByID :one
SELECT * FROM commands
WHERE id = ?;

-- name: GetCommandByName :one
SELECT * FROM commands
WHERE name = ?;

-- name: GetCommandByCommand :one
SELECT * FROM commands
WHERE command = ?;

-- name: ListCommands :many
SELECT * FROM commands
ORDER BY created_at DESC;

-- name: UpdateCommand :one
UPDATE commands
SET name = ?, command = ?, parameters = ?, description = ?
WHERE id = ?
RETURNING *;

-- name: UpdateCommandByName :one
UPDATE commands
SET name = ?, command = ?, parameters = ?, description = ?
WHERE name = ?
RETURNING *;

-- name: DeleteCommandByID :exec
DELETE FROM commands
WHERE id = ?;

-- name: DeleteCommandByName :exec
DELETE FROM commands
WHERE name = ?;
