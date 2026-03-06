# Wails Backend API

All functions in this module are auto-generated Wails bindings that make RPC calls to the Go backend. They return Promises and are available in the browser/WebView context only.

## Import

```typescript
import {
  GetAuthStatus, SetPassword, VerifyPassword,
  ListNotebooks, GetNotebook, CreateNotebook, UpdateNotebookTitle,
  UpdateNotebookPosition, MoveNotebook, DeleteNotebook,
  ListNotes, GetNote, CreateNote, UpdateNoteTitle,
  UpdateNoteContent, UpdateNotePosition, MoveNote, DeleteNote,
  SaveImage, GetImage, GetImagePath, ListImagesByNote, DeleteImage, DownloadImage,
} from '../wailsjs/go/main/App'
import { dto } from '../wailsjs/go/models'
```

---

## Authentication

### GetAuthStatus

```typescript
GetAuthStatus(): Promise<{ isPasswordSet: boolean }>
```
{ .api }

Returns the current auth state. Call this on app startup to determine whether a password has been set. Response is an untyped `any` but reliably contains `{ isPasswordSet: boolean }`.

**Example:**
```typescript
const status = await GetAuthStatus()
if (status.isPasswordSet) {
  // show password entry
} else {
  // show set-password form
}
```

---

### SetPassword

```typescript
SetPassword(password: string): Promise<{ success: boolean; error?: string }>
```
{ .api }

Sets the application password (Argon2id hash stored in backend). Can only be called once — returns an error if already set.

**Parameters:**
- `password` — The password string to set.

**Returns:** `{ success: boolean, error?: string }` — `success: false` if password already set (`ErrConflict` from backend) or if password hashing fails.

**Example:**
```typescript
const res = await SetPassword('mySecretPass')
if (!res.success) {
  console.error(res.error)
}
```

---

### VerifyPassword

```typescript
VerifyPassword(password: string): Promise<{ success: boolean }>
```
{ .api }

Verifies the password against the stored Argon2id hash using constant-time comparison.

**Parameters:**
- `password` — The password to verify.

**Returns:** `{ success: boolean }` — `success: false` if password does not match.

---

## Notebooks

Notebooks are hierarchically nested containers. `parentId = null/""` means root level.

### ListNotebooks

```typescript
ListNotebooks(parentId: string): Promise<dto.NotebookListItem[]>
```
{ .api }

Lists all notebooks under the given parent. Pass `""` (empty string) for root-level notebooks. Returns list items without content (id, title, position, timestamps).

**Parameters:**
- `parentId` — Parent notebook ID, or `""` for root.

**Example:**
```typescript
const rootNotebooks = await ListNotebooks('')
const childNotebooks = await ListNotebooks('notebook-uuid')
```

---

### GetNotebook

```typescript
GetNotebook(id: string): Promise<dto.NotebookDetail>
```
{ .api }

Retrieves full details of a single notebook.

**Parameters:**
- `id` — Notebook UUID.

**Returns:** `dto.NotebookDetail` with `{ id, parentId?, title, position, createdAt, updatedAt }`.

---

### CreateNotebook

```typescript
CreateNotebook(req: dto.CreateNotebookRequest): Promise<dto.NotebookDetail>
```
{ .api }

Creates a new notebook. Omit `parentId` to create at root level.

**Parameters:**
- `req.parentId` — (optional) Parent notebook ID. Omit for root.
- `req.title` — Notebook title string.

**Returns:** The created `dto.NotebookDetail`.

**Example:**
```typescript
const nb = await CreateNotebook({ title: 'My Notebook' })            // root
const sub = await CreateNotebook({ parentId: nb.id, title: 'Sub' }) // nested
```

---

### UpdateNotebookTitle

```typescript
UpdateNotebookTitle(id: string, req: dto.UpdateNotebookTitleRequest): Promise<dto.NotebookDetail>
```
{ .api }

Renames a notebook.

**Parameters:**
- `id` — Notebook UUID.
- `req.title` — New title.

**Returns:** Updated `dto.NotebookDetail`.

---

### UpdateNotebookPosition

```typescript
UpdateNotebookPosition(id: string, req: dto.UpdatePositionRequest): Promise<void>
```
{ .api }

Updates the sort position of a notebook within its parent. Position values use multiples of 1000 (e.g., 0, 1000, 2000) to allow reordering without renumbering all items.

**Parameters:**
- `id` — Notebook UUID.
- `req.position` — New integer position value.

---

### MoveNotebook

```typescript
MoveNotebook(id: string, req: dto.MoveNotebookRequest): Promise<dto.NotebookDetail>
```
{ .api }

Moves a notebook to a different parent. Omit `parentId` to move to root. Cannot move a notebook into one of its own descendants.

**Parameters:**
- `id` — Notebook UUID to move.
- `req.parentId` — (optional) New parent ID. Omit to move to root.

**Returns:** Updated `dto.NotebookDetail`.

---

### DeleteNotebook

```typescript
DeleteNotebook(id: string): Promise<void>
```
{ .api }

Deletes a notebook and all its children recursively (notebooks, notes, and images) via database CASCADE.

**Parameters:**
- `id` — Notebook UUID to delete.

---

## Notes

Notes contain TipTap JSON content. `notebookId = null/""` means root level.

### ListNotes

```typescript
ListNotes(notebookId: string): Promise<dto.NoteListItem[]>
```
{ .api }

Lists all notes under the given notebook. Pass `""` for root-level notes. Returns list items without content.

**Parameters:**
- `notebookId` — Notebook UUID, or `""` for root.

---

### GetNote

```typescript
GetNote(id: string): Promise<dto.NoteDetail>
```
{ .api }

Retrieves a single note including its full TipTap JSON content.

**Parameters:**
- `id` — Note UUID.

**Returns:** `dto.NoteDetail` with `{ id, notebookId?, title, content, position, createdAt, updatedAt }`. `content` is a TipTap JSON string (e.g., `'{"type":"doc","content":[...]}'`).

---

### CreateNote

```typescript
CreateNote(req: dto.CreateNoteRequest): Promise<dto.NoteDetail>
```
{ .api }

Creates a new note. Omit `notebookId` to create at root level.

**Parameters:**
- `req.notebookId` — (optional) Parent notebook ID. Omit for root.
- `req.title` — Note title.

**Returns:** The created `dto.NoteDetail`.

---

### UpdateNoteTitle

```typescript
UpdateNoteTitle(id: string, req: dto.UpdateNoteTitleRequest): Promise<dto.NoteDetail>
```
{ .api }

Renames a note.

**Parameters:**
- `id` — Note UUID.
- `req.title` — New title.

**Returns:** Updated `dto.NoteDetail`.

---

### UpdateNoteContent

```typescript
UpdateNoteContent(id: string, req: dto.UpdateNoteContentRequest): Promise<dto.NoteDetail>
```
{ .api }

Saves note content. Content must be a TipTap JSON string. Called on every keystroke (autosave).

**Parameters:**
- `id` — Note UUID.
- `req.content` — TipTap JSON string.

**Returns:** Updated `dto.NoteDetail`.

**Example:**
```typescript
const editor = useEditor({ /* TipTap config */ })
// Autosave on content change:
watch(() => editor.value?.getJSON(), async (json) => {
  if (json && noteId) {
    await UpdateNoteContent(noteId, { content: JSON.stringify(json) })
  }
})
```

---

### UpdateNotePosition

```typescript
UpdateNotePosition(id: string, req: dto.UpdatePositionRequest): Promise<void>
```
{ .api }

Updates the sort position of a note within its notebook.

**Parameters:**
- `id` — Note UUID.
- `req.position` — New integer position value (use multiples of 1000).

---

### MoveNote

```typescript
MoveNote(id: string, req: dto.MoveNoteRequest): Promise<dto.NoteDetail>
```
{ .api }

Moves a note to a different notebook. Omit `notebookId` to move to root.

**Parameters:**
- `id` — Note UUID to move.
- `req.notebookId` — (optional) Destination notebook ID. Omit for root.

**Returns:** Updated `dto.NoteDetail`.

---

### DeleteNote

```typescript
DeleteNote(id: string): Promise<void>
```
{ .api }

Deletes a note and all its images (CASCADE).

**Parameters:**
- `id` — Note UUID to delete.

---

## Images

Images are stored on disk (XDG data dir) and referenced in the database. Image data is transferred as a `number[]` byte array (JavaScript `Array<number>`, not `Uint8Array`).

### SaveImage

```typescript
SaveImage(noteId: string, mimeType: string, data: Array<number>): Promise<dto.ImageItem>
```
{ .api }

Saves an image to disk and records it in the database. The image is written to disk first; the DB insert is rolled back if the file write fails.

**Parameters:**
- `noteId` — Note UUID to attach the image to.
- `mimeType` — MIME type string (e.g., `'image/png'`, `'image/jpeg'`).
- `data` — Image bytes as a plain `number[]` array.

**Returns:** `dto.ImageItem` with the saved image's metadata including `id`, `filename`, `mimeType`, `createdAt`.

**Example:**
```typescript
// Convert File to number[]
const arrayBuffer = await file.arrayBuffer()
const bytes = Array.from(new Uint8Array(arrayBuffer))
const image = await SaveImage(noteId, file.type, bytes)
```

---

### GetImage

```typescript
GetImage(id: string): Promise<dto.ImageItem>
```
{ .api }

Retrieves image metadata by ID.

**Parameters:**
- `id` — Image UUID.

**Returns:** `dto.ImageItem` with `{ id, noteId, filename, mimeType, createdAt }`.

---

### GetImagePath

```typescript
GetImagePath(filename: string): Promise<string>
```
{ .api }

Returns the filesystem path to an image file on disk.

**Parameters:**
- `filename` — The image filename (as stored in `dto.ImageItem.filename`).

**Returns:** Absolute filesystem path string. Used to construct `<img src>` URLs for display in the editor.

**Example:**
```typescript
const item = await GetImage(imageId)
const path = await GetImagePath(item.filename)
// path is like: /home/user/.local/share/gvnotes/images/uuid.png
```

---

### ListImagesByNote

```typescript
ListImagesByNote(noteId: string): Promise<Array<dto.ImageItem>>
```
{ .api }

Lists all images attached to a note.

**Parameters:**
- `noteId` — Note UUID.

**Returns:** Array of `dto.ImageItem`.

---

### DeleteImage

```typescript
DeleteImage(id: string): Promise<void>
```
{ .api }

Deletes an image from disk and from the database.

**Parameters:**
- `id` — Image UUID to delete.

---

### DownloadImage

```typescript
DownloadImage(id: string): Promise<void>
```
{ .api }

Triggers the native OS file save dialog to download the image to the user's machine.

**Parameters:**
- `id` — Image UUID to download.

---

## Wails Runtime

Selected functions from the Wails runtime (`wailsjs/runtime/runtime`). The app currently uses `BrowserOpenURL`.

### BrowserOpenURL

```typescript
import { BrowserOpenURL } from '../wailsjs/runtime/runtime'
BrowserOpenURL(url: string): void
```
{ .api }

Opens the given URL in the system's default browser. Used by the NoteEditor to open hyperlinks in external browser windows.

**Parameters:**
- `url` — URL string to open.
