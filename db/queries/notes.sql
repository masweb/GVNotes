-- name: ListNotes :many
SELECT id, title, position, created_at, updated_at
FROM notes
WHERE notebook_id IS ?
ORDER BY position ASC;

-- name: GetNote :one
SELECT id, notebook_id, title, content, position, created_at, updated_at
FROM notes
WHERE id = ?;

-- name: CreateNote :one
INSERT INTO notes (id, notebook_id, title, content, position)
VALUES (?, ?, ?, ?, ?)
RETURNING *;

-- name: UpdateNoteTitle :one
UPDATE notes
SET title = ?, updated_at = strftime('%Y-%m-%dT%H:%M:%fZ', 'now')
WHERE id = ?
RETURNING *;

-- name: UpdateNoteContent :one
UPDATE notes
SET content = ?, updated_at = strftime('%Y-%m-%dT%H:%M:%fZ', 'now')
WHERE id = ?
RETURNING *;

-- name: UpdateNotePosition :exec
UPDATE notes
SET position = ?, updated_at = strftime('%Y-%m-%dT%H:%M:%fZ', 'now')
WHERE id = ?;

-- name: DeleteNote :exec
DELETE FROM notes WHERE id = ?;

-- name: MoveNote :one
UPDATE notes
SET notebook_id = ?,
    position    = ?,
    updated_at  = strftime('%Y-%m-%dT%H:%M:%fZ', 'now')
WHERE id = ?
RETURNING id, notebook_id, title, content, position, created_at, updated_at;
