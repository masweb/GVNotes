-- name: ListNotebooks :many
SELECT id, title, position, created_at, updated_at
FROM notebooks
WHERE parent_id IS ?
ORDER BY position ASC;

-- name: GetNotebook :one
SELECT id, parent_id, title, position, created_at, updated_at
FROM notebooks
WHERE id = ?;

-- name: CreateNotebook :one
INSERT INTO notebooks (id, parent_id, title, position)
VALUES (?, ?, ?, ?)
RETURNING *;

-- name: UpdateNotebookTitle :one
UPDATE notebooks
SET title = ?, updated_at = strftime('%Y-%m-%dT%H:%M:%fZ', 'now')
WHERE id = ?
RETURNING *;

-- name: UpdateNotebookPosition :exec
UPDATE notebooks
SET position = ?, updated_at = strftime('%Y-%m-%dT%H:%M:%fZ', 'now')
WHERE id = ?;

-- name: DeleteNotebook :exec
DELETE FROM notebooks WHERE id = ?;

-- name: MoveNotebook :one
UPDATE notebooks
SET parent_id  = ?,
    position   = ?,
    updated_at = strftime('%Y-%m-%dT%H:%M:%fZ', 'now')
WHERE id = ?
RETURNING id, parent_id, title, position, created_at, updated_at;
