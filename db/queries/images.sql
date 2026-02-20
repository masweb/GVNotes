-- name: ListImagesByNote :many
SELECT id, note_id, filename, mime_type, created_at
FROM images
WHERE note_id = ?;

-- name: GetImage :one
SELECT id, note_id, filename, mime_type, created_at
FROM images
WHERE id = ?;

-- name: CreateImage :one
INSERT INTO images (id, note_id, filename, mime_type)
VALUES (?, ?, ?, ?)
RETURNING *;

-- name: DeleteImage :exec
DELETE FROM images WHERE id = ?;

-- name: DeleteImagesByNote :exec
DELETE FROM images WHERE note_id = ?;
