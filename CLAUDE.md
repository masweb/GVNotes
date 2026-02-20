# gvnotes — Claude Reference

## Communication
Always respond to the user **in Spanish (Castilian)**.

---

## Stack
Wails v2 · Go 1.24 · SQLite (mattn/go-sqlite3 CGO) · sqlc · golang-migrate · adrg/xdg · Vue 3 + TS

## Layout
```
main.go          # embeds frontend/dist + db/migrations
app.go           # App facade: startup/shutdown + DI composition
internal/
  config/
    database.go  # DBPath(), ImagesDir()
    db.go        # OpenDB() — opens SQLite, runs migrations
  controllers/   # auth.go  notebook.go  note.go  image.go
  services/      # auth.go  notebook.go  note.go  image.go
  repositories/  # only when extra logic beyond sqlc is needed
  dto/           # auth.go  notebook.go  note.go  image.go
  errors/        # ErrNotFound ErrConflict ErrForbidden ErrUnauthenticated ErrInvalidInput
db/
  migrations/    # YYYYMMDDHHMMSS_name.up/down.sql
  queries/       # settings.sql  notebooks.sql  notes.sql  images.sql
  generated/     # DO NOT EDIT — run: cd db && sqlc generate
  sqlc.yaml      # paths relative to db/; emit_interface emit_json_tags emit_pointers_for_null_types
```

## DI flow (app.go startup)
`OpenDB → db.New(sqlDB) → NewXxxService(querier) → NewXxxController(svc)` — no framework.
`App` binds to Wails and delegates to controllers. Controllers initialized in `startup()`, not `NewApp()` (Bind runs before startup).

## DB schema
- `settings(key PK, value)` — generic k/v; used for `auth.password_hash`
- `notebooks(id, parent_id nullable FK self-ref CASCADE, title, position, created_at, updated_at)`
- `notes(id, notebook_id nullable FK CASCADE, title, content JSON TipTap, position, created_at, updated_at)`
- `images(id, note_id FK CASCADE, filename, mime_type, created_at)`
- `NULL parent_id / notebook_id` = root level
- `position` is per-type (notebooks among notebooks, notes among notes) for drag & drop
- List queries return only `id, title, position, timestamps` — content only in GetNote

## Data dirs (via adrg/xdg)
- DB: `<XDG_DATA_HOME>/gvnotes/gvnotes.db`
- Images: `<XDG_DATA_HOME>/gvnotes/images/<uuid>.<ext>`
- macOS XDG_DATA_HOME = `~/Library/Application Support`

## Auth
- Single app password, no users
- Stored as `base64(salt):base64(argon2id hash)` in `settings` under key `auth.password_hash`
- argon2id params: time=1, memory=64MB, threads=4, keyLen=32, saltLen=16 (OWASP minimums)
- `IsPasswordSet` → checks row exists; `SetPassword` → errors if already set (ErrConflict); `VerifyPassword` → constant-time compare
- Wails methods: `GetAuthStatus()`, `SetPassword(password)`, `VerifyPassword(password)`

## Wails methods exposed
Auth: `GetAuthStatus` `SetPassword` `VerifyPassword`
Notebooks: `ListNotebooks(parentID)` `GetNotebook` `CreateNotebook` `UpdateNotebookTitle` `UpdateNotebookPosition` `DeleteNotebook`
Notes: `ListNotes(notebookID)` `GetNote` `CreateNote` `UpdateNoteTitle` `UpdateNoteContent` `UpdateNotePosition` `DeleteNote`
Images: `SaveImage(noteID, mimeType, []byte)` `GetImage` `GetImagePath(filename)` `ListImagesByNote` `DeleteImage`
— `ListNotebooks("")` / `ListNotes("")` → root level (converts empty string to nil)

## Rules
- `sqlc generate` must be run from `db/` directory
- Controllers never receive `*sql.DB`
- Never compare error strings — use `errors.Is` against typed errors
- Repositories only when sqlc querier is insufficient
- Images: file written to disk first, DB insert after; file removed on DB failure (best-effort)
