# Pinia Stores

Application state is managed via two Pinia stores. Both are auto-imported in Vue SFCs.

## Imports

```typescript
// In Vue SFCs (auto-imported, no explicit import needed):
const authStore = useAuthStore()
const navStore = useNavigationStore()

// In plain TypeScript files:
import { useAuthStore } from '@/stores/auth'
import { useNavigationStore, type NavItem, type NavLevel } from '@/stores/navigation'
```

---

## useAuthStore

```typescript
const store = useAuthStore()
```
{ .api }

Manages authentication state. Bridges the Wails auth backend bindings with reactive Vue state.

### State

| Property | Type | Description |
|----------|------|-------------|
| `status` | `Ref<'unauthenticated' \| 'authenticated'>` | Current auth status. |
| `passwordSet` | `Ref<boolean>` | Whether an app password has been configured in the backend. |
| `loading` | `Ref<boolean>` | `true` while an async action is running. |
| `error` | `Ref<string \| null>` | Last error message (i18n-translated). Cleared on next action call. |

### Computed

| Property | Type | Description |
|----------|------|-------------|
| `isAuthenticated` | `ComputedRef<boolean>` | Shorthand for `status.value === 'authenticated'`. |

### Actions

#### init

```typescript
init(): Promise<void>
```
{ .api }

Fetches auth status from the backend (`GetAuthStatus`). Sets `passwordSet`. Should be called once on app mount. Does not change `status` (authentication still requires `verifyPassword`).

---

#### setPassword

```typescript
setPassword(password: string): Promise<boolean>
```
{ .api }

Sets the application password for the first time. Returns `true` on success and sets `status` to `'authenticated'`. Returns `false` if a password is already set or an error occurs; sets `store.error` with a translated message.

**Parameters:**
- `password` — New password string.

---

#### verifyPassword

```typescript
verifyPassword(password: string): Promise<boolean>
```
{ .api }

Verifies the entered password against the stored hash. Returns `true` and sets `status` to `'authenticated'` on success. Returns `false` with `store.error` set on failure.

**Parameters:**
- `password` — Password string to verify.

---

#### logout

```typescript
logout(): void
```
{ .api }

Sets `status` to `'unauthenticated'` and clears `error`. Does not clear `passwordSet`. Does not call any backend endpoint.

---

### Auth Type

```typescript
type AuthStatus = 'unauthenticated' | 'authenticated'
```
{ .api }

### Usage Example

```typescript
const auth = useAuthStore()

onMounted(async () => {
  await auth.init()
  if (!auth.passwordSet) {
    // First run — show set-password form
  } else {
    // Returning user — show verify-password form
  }
})

// Setting password (first run)
const ok = await auth.setPassword(enteredPassword)
if (!ok) console.error(auth.error)

// Verifying password (login)
const ok = await auth.verifyPassword(enteredPassword)

// Checking auth in watchers/guards
watchEffect(() => {
  if (!auth.isAuthenticated) router.push('/login')
})
```

---

## useNavigationStore

```typescript
const store = useNavigationStore()
```
{ .api }

Manages hierarchical navigation. Maintains a stack of `NavLevel` objects representing the user's current drill-down path through notebooks.

### Exported Types

#### NavItem

```typescript
type NavItem =
  | { kind: 'notebook'; data: dto.NotebookListItem }
  | { kind: 'note'; data: dto.NoteListItem }
```
{ .api }

A discriminated union representing either a notebook or a note in the navigation tree.

#### NavLevel

```typescript
interface NavLevel {
  parentId: string | null  // null = root level
  title: string
  items: NavItem[]
}
```
{ .api }

Represents one level in the navigation hierarchy. `parentId` is the ID of the notebook whose contents are shown; `null` means root.

---

### State

| Property | Type | Description |
|----------|------|-------------|
| `stack` | `Ref<NavLevel[]>` | Navigation stack. Index 0 = root level. Last item = deepest open level. |
| `selectedNoteId` | `Ref<string \| null>` | ID of the currently selected (editing) note. |
| `loading` | `Ref<boolean>` | `true` while loading a level from the backend. |

### Computed

| Property | Type | Description |
|----------|------|-------------|
| `leftPanel` | `ComputedRef<NavLevel \| null>` | Second-to-last stack entry (or first if stack length is 1). Drives the left pane. |
| `rightPanel` | `ComputedRef<NavLevel \| null>` | Last stack entry; `null` when stack has only one level. Drives the right pane. |
| `canGoBack` | `ComputedRef<boolean>` | `true` when `stack.length > 1`. |
| `activeNotebookId` | `ComputedRef<string \| null>` | The `parentId` of `rightPanel` (the deepest open notebook). `null` when at root. |

---

### Actions

#### init

```typescript
init(): Promise<void>
```
{ .api }

Loads the root level from the backend (both notebooks and notes at root) and sets `stack` to a single-item array. Should be called after authentication.

---

#### openNotebook

```typescript
openNotebook(notebook: dto.NotebookListItem): Promise<void>
```
{ .api }

Drills into a notebook by pushing a new `NavLevel` onto the stack. Clears `selectedNoteId`. Used when clicking a notebook in the right panel or the only panel.

**Parameters:**
- `notebook` — A `dto.NotebookListItem` from the current level's items.

---

#### openNotebookFromLeft

```typescript
openNotebookFromLeft(notebook: dto.NotebookListItem): Promise<void>
```
{ .api }

Replaces the top-of-stack level with the notebook's contents. Used when clicking a notebook in the **left** panel to change what the right panel shows.

**Parameters:**
- `notebook` — A `dto.NotebookListItem` from the left panel's items.

---

#### selectNote

```typescript
selectNote(noteId: string): void
```
{ .api }

Sets `selectedNoteId` to the given note ID. Triggers the note editor to load that note.

---

#### selectNoteFromRoot

```typescript
selectNoteFromRoot(noteId: string): void
```
{ .api }

Collapses the stack to the root level and selects the note. Used for search or direct note navigation.

---

#### goBack

```typescript
goBack(): void
```
{ .api }

Pops the top level from the stack (goes one level up). Clears `selectedNoteId`. Does nothing if already at root.

---

#### reorderItems

```typescript
reorderItems(level: NavLevel, kind: 'notebook' | 'note', newItems: NavItem[]): Promise<void>
```
{ .api }

Updates the order of notebooks or notes in a level after a drag-and-drop operation. Mutates the stack immediately for optimistic UI, then persists new positions to the backend (`UpdateNotebookPosition` / `UpdateNotePosition` with `position = index * 1000`).

**Parameters:**
- `level` — The `NavLevel` whose items are being reordered.
- `kind` — `'notebook'` or `'note'`.
- `newItems` — The reordered items array (only items of `kind` are reordered; the other kind is preserved).

---

#### reset

```typescript
reset(): void
```
{ .api }

Clears `stack` and `selectedNoteId`. Used on logout.

---

#### renameNotebook

```typescript
renameNotebook(notebookId: string, newTitle: string): Promise<void>
```
{ .api }

Persists the notebook title to the backend and updates the stack: updates the `title` of the `NavLevel` whose `parentId` is this notebook and updates the item's title in its parent level.

**Parameters:**
- `notebookId` — Notebook UUID.
- `newTitle` — New title string.

---

#### renameItem

```typescript
renameItem(itemId: string, kind: 'notebook' | 'note', newTitle: string): void
```
{ .api }

Updates the title of a `NavItem` in the stack (local only, no backend call). Use `renameNotebook` for notebooks if you also need to persist the title.

---

#### removeItem

```typescript
removeItem(levelParentId: string | null, itemId: string, itemKind: 'notebook' | 'note'): void
```
{ .api }

Removes a `NavItem` from the specified level in the stack (local only). If `itemKind === 'notebook'` and that notebook is currently open in the stack, the stack is collapsed up to that point.

**Parameters:**
- `levelParentId` — `parentId` of the `NavLevel` containing the item.
- `itemId` — ID of the item to remove.
- `itemKind` — `'notebook'` or `'note'`.

---

#### moveItem

```typescript
moveItem(sourceLevelParentId: string | null, item: NavItem, destinationParentId: string | null): void
```
{ .api }

Moves a `NavItem` between levels in the stack (local only, no backend call). Removes it from the source level and inserts it at the destination level if that level is currently in the stack. If the moved item is a notebook that was open in the stack, the stack is collapsed to its parent.

**Parameters:**
- `sourceLevelParentId` — `parentId` of the source `NavLevel`.
- `item` — The `NavItem` to move.
- `destinationParentId` — `parentId` of the destination `NavLevel`, or `null` for root.

---

### Usage Example

```typescript
const nav = useNavigationStore()

// Initialize after auth
await nav.init()

// Open a notebook (right panel click or single panel)
await nav.openNotebook(notebookItem)

// Left panel click — change right panel content without drilling deeper
await nav.openNotebookFromLeft(notebookItem)

// Select note for editing
nav.selectNote(noteId)

// Go back one level
if (nav.canGoBack) nav.goBack()

// After drag-and-drop reorder
await nav.reorderItems(nav.leftPanel!, 'notebook', reorderedNotebooks)

// After successful delete via API
nav.removeItem(level.parentId, item.data.id, item.kind)

// After successful move via API
nav.moveItem(sourceLevel.parentId, item, destinationNotebookId)

// On logout
auth.logout()
nav.reset()
```
