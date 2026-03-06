# Vue Components

All components are auto-registered globally via `unplugin-vue-components`. They can be used in any Vue SFC template without importing.

---

## NoteEditor

Rich text editor component powered by TipTap. Handles note loading, autosave, title editing, and image upload.

### Props

```typescript
defineProps<{ noteId: string }>()
```
{ .api }

| Prop | Type | Required | Description |
|------|------|----------|-------------|
| `noteId` | `string` | Yes | UUID of the note to display and edit. Changes trigger a reload of the note content. |

### Emits

```typescript
defineEmits<{ titleChanged: [id: string, title: string] }>()
```
{ .api }

| Event | Payload | Description |
|-------|---------|-------------|
| `titleChanged` | `id: string, title: string` | Emitted after the user renames the note title. Carries the note ID and new title so parent can update navigation state. |

### Features

- Loads note by `noteId` via `GetNote()` on mount and on prop change
- Autosaves content on every keystroke with a debounce via `UpdateNoteContent()`
- Inline title editing with `required` validation; saves via `UpdateNoteTitle()`
- Full formatting toolbar: bold, italic, underline, strikethrough, text color, text alignment, headings (H1–H6), bullet lists, ordered lists, task lists, inline code, code blocks (syntax highlighted via lowlight), tables, images, hyperlinks
- Image upload: paste or drag-and-drop converts to byte array and calls `SaveImage()`; images rendered via `GetImagePath()`
- Image download: calls `DownloadImage(id)`
- Link handling: opens external links via `BrowserOpenURL()` (native OS browser)
- Character count display
- TipTap extensions used: Bold, Italic, Underline, Strike, TextStyle, Color, TextAlign, Heading, Paragraph, Document, Text, History, Link, Image, Table, TableRow, TableCell, TableHeader, TaskList, TaskItem, CodeBlockLowlight, CharacterCount

### Usage

```html
<!-- In GVNotes.vue or any parent -->
<NoteEditor
  v-if="nav.selectedNoteId"
  :noteId="nav.selectedNoteId"
  @titleChanged="(id, title) => nav.renameItem(id, 'note', title)"
/>
```

---

## NavPanel

Navigation panel showing hierarchical list of notebooks and notes with drag-and-drop reordering and context menu actions.

### Props

```typescript
defineProps<{
  level: NavLevel
  activeNoteId?: string | null
  activeNotebookId?: string | null
}>()
```
{ .api }

| Prop | Type | Required | Description |
|------|------|----------|-------------|
| `level` | `NavLevel` | Yes | The navigation level to display (from `useNavigationStore`). |
| `activeNoteId` | `string \| null` | No | ID of the currently selected note; highlighted in the list. |
| `activeNotebookId` | `string \| null` | No | ID of the notebook currently open in the right pane; highlighted. |

### Emits

```typescript
defineEmits<{
  notebookClick: [item: NavItem & { kind: 'notebook' }]
  noteClick: [item: NavItem & { kind: 'note' }]
  reorder: [level: NavLevel, kind: 'notebook' | 'note', newItems: NavItem[]]
  create: [kind: 'notebook' | 'note']
  delete: [item: NavItem]
  move: [item: NavItem]
  rename: [level: NavLevel, newTitle: string]
}>()
```
{ .api }

| Event | Payload | Description |
|-------|---------|-------------|
| `notebookClick` | `NavItem & { kind: 'notebook' }` | Emitted when the user clicks a notebook item. |
| `noteClick` | `NavItem & { kind: 'note' }` | Emitted when the user clicks a note item. |
| `reorder` | `level: NavLevel, kind, newItems: NavItem[]` | Emitted after drag-and-drop reorder. Pass to `nav.reorderItems()`. |
| `create` | `kind: 'notebook' \| 'note'` | Emitted when the user triggers create in context menu. |
| `delete` | `NavItem` | Emitted when the user triggers delete in context menu. |
| `move` | `NavItem` | Emitted when the user triggers move in context menu. |
| `rename` | `level: NavLevel, newTitle: string` | Emitted when the user renames a notebook header. |

### Features

- Displays notebooks and notes separated by category
- Client-side search/filter by title
- Drag-and-drop reordering via SortableJS
- Context menu (right-click or long-press) with create/delete/move/rename actions

---

## NavPanelHeader

Header for a navigation panel showing the level title with create and rename actions.

### Props

```typescript
defineProps<{ level: NavLevel }>()
```
{ .api }

| Prop | Type | Required | Description |
|------|------|----------|-------------|
| `level` | `NavLevel` | Yes | The navigation level whose title is displayed. |

### Emits

```typescript
defineEmits<{
  create: [kind: 'notebook' | 'note']
  rename: [newTitle: string]
}>()
```
{ .api }

| Event | Payload | Description |
|-------|---------|-------------|
| `create` | `kind: 'notebook' \| 'note'` | User selected create from the dropdown. |
| `rename` | `newTitle: string` | User submitted a new title for the level (for notebooks). |

---

## CreateItemModal

Modal dialog for creating a new notebook or note with name validation.

### Props

```typescript
defineProps<{
  visible: boolean
  kind: 'notebook' | 'note'
  level: NavLevel | null
}>()
```
{ .api }

| Prop | Type | Required | Description |
|------|------|----------|-------------|
| `visible` | `boolean` | Yes | Controls visibility. Use `v-model:visible`. |
| `kind` | `'notebook' \| 'note'` | Yes | Whether to create a notebook or a note. |
| `level` | `NavLevel \| null` | Yes | The navigation level where the item is created. |

### Emits

```typescript
defineEmits<{ 'update:visible': [value: boolean] }>()
```
{ .api }

| Event | Payload | Description |
|-------|---------|-------------|
| `update:visible` | `boolean` | Emitted to close the modal (`false`). Use with `v-model:visible`. |

### Features

Internally calls `useNavActions().createNotebook()` or `createNote()` on submit. No external create event is emitted — the store is updated directly.

### Usage

```html
<CreateItemModal
  v-model:visible="modalVisible"
  :kind="modalKind"
  :level="modalLevel"
/>
```

---

## DeleteConfirmModal

Simple delete confirmation dialog.

### Props

```typescript
defineProps<{
  visible: boolean
  kind: 'notebook' | 'note'
  title: string
}>()
```
{ .api }

| Prop | Type | Required | Description |
|------|------|----------|-------------|
| `visible` | `boolean` | Yes | Controls visibility. Use `v-model:visible`. |
| `kind` | `'notebook' \| 'note'` | Yes | Type of item being deleted (affects icon and i18n text). |
| `title` | `string` | Yes | Name of the item to show in the confirmation message. |

### Emits

```typescript
defineEmits<{
  'update:visible': [value: boolean]
  confirm: []
}>()
```
{ .api }

| Event | Payload | Description |
|-------|---------|-------------|
| `update:visible` | `boolean` | Emitted to close the modal. |
| `confirm` | none | Emitted when user clicks the confirm button. Parent must handle the actual deletion. |

### Usage

```html
<DeleteConfirmModal
  v-model:visible="deleteModalVisible"
  :kind="deleteTarget?.kind"
  :title="deleteTarget?.data.title"
  @confirm="onDeleteConfirm"
/>
```

---

## MoveItemModal

Modal dialog for selecting a destination notebook when moving a notebook or note. Displays a tree of notebooks (root + expandable children).

### Props

```typescript
defineProps<{
  visible: boolean
  item: NavItem | null
  level: NavLevel | null
}>()
```
{ .api }

| Prop | Type | Required | Description |
|------|------|----------|-------------|
| `visible` | `boolean` | Yes | Controls visibility. Use `v-model:visible`. |
| `item` | `NavItem \| null` | Yes | The item being moved. |
| `level` | `NavLevel \| null` | Yes | The source level (used to determine current parent). |

### Emits

```typescript
defineEmits<{ 'update:visible': [value: boolean] }>()
```
{ .api }

| Event | Payload | Description |
|-------|---------|-------------|
| `update:visible` | `boolean` | Emitted to close the modal. |

### Features

Internally calls `MoveNotebook()` or `MoveNote()` backend bindings and then `nav.moveItem()`. No external move event — store is updated directly.

### Usage

```html
<MoveItemModal
  v-model:visible="moveModalVisible"
  :item="moveTarget"
  :level="moveLevel"
/>
```

---

## AppSettingsModal

Settings dialog providing theme toggling, language selection, and logout.

### Props

```typescript
defineProps<{ visible: boolean }>()
```
{ .api }

| Prop | Type | Required | Description |
|------|------|----------|-------------|
| `visible` | `boolean` | Yes | Controls visibility. Use `v-model:visible`. |

### Emits

```typescript
defineEmits<{ 'update:visible': [value: boolean] }>()
```
{ .api }

| Event | Payload | Description |
|-------|---------|-------------|
| `update:visible` | `boolean` | Emitted to close the modal. |

### Features

- Theme toggle (light ↔ dark) via `useTheme()`
- Language selector (ES/EN) via `useLocale()`
- Lock/logout button: calls `nav.reset()` + `auth.logout()`

---

## NotebookTreeNode

Recursive tree node for displaying a notebook and its children in the MoveItemModal tree.

### Props

```typescript
defineProps<{
  notebook: dto.NotebookListItem
  selectedId: string | null
  disabledId: string | null
}>()
```
{ .api }

| Prop | Type | Required | Description |
|------|------|----------|-------------|
| `notebook` | `dto.NotebookListItem` | Yes | The notebook to render as a tree node. |
| `selectedId` | `string \| null` | Yes | ID of the currently selected destination; highlights this node. |
| `disabledId` | `string \| null` | Yes | ID of the item being moved; this node is non-selectable (can't move into itself). |

### Emits

```typescript
defineEmits<{ select: [id: string] }>()
```
{ .api }

| Event | Payload | Description |
|-------|---------|-------------|
| `select` | `id: string` | Emitted when user selects this notebook as the move destination. |

### Features

- Expand/collapse toggle
- Lazy-loads children via `ListNotebooks(notebook.id)` on first expand
- Propagates `select` events from child `NotebookTreeNode` components upward

---

## AuthView

Authentication screen shown before the main app. Handles both first-run (set password) and returning user (verify password) scenarios.

No props or emits. Reads state from `useAuthStore()` directly and renders the appropriate form.

### Features

- Set password form (first run): two fields (password + confirm), `min:8` and `confirmed` validation
- Verify password form (returning): single password field, `required` validation
- Uses vee-validate `useForm` + `useField` with `validateOnMount: false`
- All validation triggered only on submit (`handleSubmit`)

---

## GVNotes

Main application view. Renders the split-pane layout with left/right navigation panels and the note editor. No props or emits — reads and writes navigation store state directly.

### Features

- Two-panel layout via Splitpanes (left: `nav.leftPanel`, right: `nav.rightPanel` or editor)
- Back navigation button (visible when `nav.canGoBack`)
- Keyboard shortcuts: `Cmd/Ctrl+N` → create note, `Cmd/Ctrl+Shift+N` → create notebook
- Hosts all modal dialogs: CreateItemModal, DeleteConfirmModal, MoveItemModal, AppSettingsModal
- Settings button via AppSettingsModal
