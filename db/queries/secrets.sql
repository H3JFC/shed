-- name: CreateSecret :one
INSERT INTO secrets (key, value)
VALUES (?, ?)
RETURNING *;

-- name: GetSecretByID :one
SELECT * FROM secrets
WHERE id = ?;

-- name: GetSecretByKey :one
SELECT * FROM secrets
WHERE key = ?;

-- name: ListSecrets :many
SELECT * FROM secrets
ORDER BY created_at DESC;

-- name: UpdateSecret :one
UPDATE secrets
SET key = ?, value = ?
WHERE id = ?
RETURNING *;

-- name: UpdateSecretByKey :one
UPDATE secrets
SET value = ?
WHERE key = ?
RETURNING *;

-- name: DeleteSecretByID :exec
DELETE FROM secrets
WHERE id = ?;

-- name: DeleteSecretByKey :exec
DELETE FROM secrets
WHERE key = ?;
