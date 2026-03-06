# Data Types (DTOs)

All backend data transfer objects (DTOs) are in the `dto` namespace exported from `wailsjs/go/models.ts`.

## Import

```typescript
import { dto } from '../wailsjs/go/models'
// Or with type-only import:
import type { dto } from '../wailsjs/go/models'
```

---

## Request DTOs

### dto.CreateNotebookRequest

```typescript
class CreateNotebookRequest {
  parentId?: string   // Omit for root level
  title: string
  static createFrom(source: any): CreateNotebookRequest
}
```
{ .api }

Payload for `CreateNotebook`. `parentId` is optional — omit it (or pass `undefined`) to create a root-level notebook.

---

### dto.CreateNoteRequest

```typescript
class CreateNoteRequest {
  notebookId?: string  // Omit for root level
  title: string
  static createFrom(source: any): CreateNoteRequest
}
```
{ .api }

Payload for `CreateNote`. `notebookId` is optional — omit to create a root-level note.

---

### dto.MoveNotebookRequest

```typescript
class MoveNotebookRequest {
  parentId?: string   // Omit to move to root
  static createFrom(source: any): MoveNotebookRequest
}
```
{ .api }

Payload for `MoveNotebook`. `parentId` is optional — omit to move to root.

---

### dto.MoveNoteRequest

```typescript
class MoveNoteRequest {
  notebookId?: string  // Omit to move to root
  static createFrom(source: any): MoveNoteRequest
}
```
{ .api }

Payload for `MoveNote`. `notebookId` is optional — omit to move to root.

---

### dto.UpdateNotebookTitleRequest

```typescript
class UpdateNotebookTitleRequest {
  title: string
  static createFrom(source: any): UpdateNotebookTitleRequest
}
```
{ .api }

Payload for `UpdateNotebookTitle`.

---

### dto.UpdateNoteTitleRequest

```typescript
class UpdateNoteTitleRequest {
  title: string
  static createFrom(source: any): UpdateNoteTitleRequest
}
```
{ .api }

Payload for `UpdateNoteTitle`.

---

### dto.UpdateNoteContentRequest

```typescript
class UpdateNoteContentRequest {
  content: string   // TipTap JSON string
  static createFrom(source: any): UpdateNoteContentRequest
}
```
{ .api }

Payload for `UpdateNoteContent`. `content` must be a valid TipTap JSON string (serialized via `editor.getJSON()` and `JSON.stringify()`).

---

### dto.UpdatePositionRequest

```typescript
class UpdatePositionRequest {
  position: number   // Integer; use multiples of 1000
  static createFrom(source: any): UpdatePositionRequest
}
```
{ .api }

Payload for `UpdateNotebookPosition` and `UpdateNotePosition`. Positions use multiples of 1000 (0, 1000, 2000, …) to allow insertion without renumbering all items.

---

## Response DTOs

### dto.NotebookListItem

```typescript
class NotebookListItem {
  id: string
  title: string
  position: number
  createdAt: any   // Go time.Time serialized as JSON
  updatedAt: any   // Go time.Time serialized as JSON
  static createFrom(source: any): NotebookListItem
}
```
{ .api }

Lightweight notebook representation returned by `ListNotebooks`. Does not include child counts or content. Used as items in `NavItem` discriminated union when `kind === 'notebook'`.

---

### dto.NotebookDetail

```typescript
class NotebookDetail {
  id: string
  parentId?: string   // undefined/null for root-level notebooks
  title: string
  position: number
  createdAt: any
  updatedAt: any
  static createFrom(source: any): NotebookDetail
}
```
{ .api }

Full notebook representation returned by `GetNotebook`, `CreateNotebook`, `UpdateNotebookTitle`, and `MoveNotebook`.

---

### dto.NoteListItem

```typescript
class NoteListItem {
  id: string
  title: string
  position: number
  createdAt: any
  updatedAt: any
  static createFrom(source: any): NoteListItem
}
```
{ .api }

Lightweight note representation returned by `ListNotes`. Does not include content. Used as items in `NavItem` discriminated union when `kind === 'note'`.

---

### dto.NoteDetail

```typescript
class NoteDetail {
  id: string
  notebookId?: string   // undefined/null for root-level notes
  title: string
  content: string       // TipTap JSON string
  position: number
  createdAt: any
  updatedAt: any
  static createFrom(source: any): NoteDetail
}
```
{ .api }

Full note representation returned by `GetNote`, `CreateNote`, `UpdateNoteTitle`, `UpdateNoteContent`, and `MoveNote`. `content` is a TipTap JSON string; parse with `JSON.parse(detail.content)` to get TipTap's document structure.

---

### dto.ImageItem

```typescript
class ImageItem {
  id: string
  noteId: string
  filename: string   // UUID-based filename with extension, e.g. "550e8400-e29b-41d4-a716-446655440000.png"
  mimeType: string   // e.g. "image/png", "image/jpeg"
  createdAt: any
  static createFrom(source: any): ImageItem
}
```
{ .api }

Image metadata returned by `SaveImage`, `GetImage`, and `ListImagesByNote`. To display an image, use `GetImagePath(item.filename)` to get the filesystem path.

---

## createFrom Static Factory

All DTO classes expose a static `createFrom` method:

```typescript
static createFrom(source: any = {}): T
```
{ .api }

Constructs a typed DTO from a plain object or JSON string. This is the standard way to instantiate DTOs from raw API responses in custom code:

```typescript
const item = dto.NotebookListItem.createFrom({ id: 'abc', title: 'Notes', position: 0, createdAt: ..., updatedAt: ... })
```

---

## Timestamps

The `createdAt` and `updatedAt` fields on all DTOs are typed as `any` because they originate from Go's `time.Time`. When serialized through the Wails runtime they arrive as JavaScript `Date`-compatible ISO 8601 strings or numeric timestamps, depending on the serialization path. Treat them as opaque values for display; use `new Date(item.createdAt)` for formatting.

---

## Position Convention

`position` fields use integer values spaced by **1000** (0, 1000, 2000, …). When persisting a reorder, multiply each item's new array index by 1000:

```typescript
items.forEach((item, idx) => {
  UpdateNotebookPosition(item.data.id, { position: idx * 1000 })
})
```

This pattern leaves room for future insertions without renumbering all items.
