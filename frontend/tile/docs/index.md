# GVNotes Frontend (frontend)

## Package Information

- **Package**: `frontend` (Vue 3 + TypeScript, Wails v2 desktop app)
- **Version**: 0.0.0
- **Package Manager**: pnpm
- **Language**: TypeScript + Vue 3 SFCs
- **Runtime**: Wails v2 (Go + native WebView, no Electron)

## Overview

The `frontend` package is the Vue 3 + TypeScript user interface for GVNotes, a fast, offline-first desktop note-taking app. It communicates with a Go backend exclusively via Wails-generated JavaScript bindings. The application provides hierarchical notebook/note management with drag-and-drop, a TipTap-based rich text editor, optional app password authentication, and full light/dark theme support.

## Architecture

```
App Bootstrap (src/main.ts)
├── Vue app with Pinia + vue-i18n + CoreUI
├── App.vue  ←── root; renders AuthView or GVNotes based on auth state
│
├── Wails Backend Bindings (wailsjs/go/main/App)
│   ├── Auth: GetAuthStatus / SetPassword / VerifyPassword
│   ├── Notebooks: List / Get / Create / Update / Move / Delete
│   ├── Notes:     List / Get / Create / Update / Move / Delete
│   └── Images:    Save / Get / GetPath / List / Delete / Download
│
├── Pinia Stores (src/stores/)
│   ├── useAuthStore    — auth status, password set/verify, logout
│   └── useNavigationStore — nav stack, level loading, item CRUD
│
├── Composables (src/composables/)
│   ├── useTheme        — light/dark theme, localStorage persistence
│   ├── useLocale       — i18n locale switching, localStorage persistence
│   ├── useNavActions   — create/delete notebook & note via backend + store
│   └── useValidation   — global vee-validate rules (required, min, confirmed)
│
└── Views / Components (auto-registered globally)
    ├── AuthView        — set/verify password screen
    ├── GVNotes         — main split-pane view
    ├── NavPanel        — hierarchical tree panel
    ├── NoteEditor      — TipTap rich text editor
    └── (modals: Create, Delete, Move, Settings)
```

## Core Imports

```typescript
// Wails backend bindings
import { GetAuthStatus, SetPassword, VerifyPassword } from '../wailsjs/go/main/App'
import { ListNotebooks, GetNotebook, CreateNotebook, UpdateNotebookTitle,
         UpdateNotebookPosition, MoveNotebook, DeleteNotebook } from '../wailsjs/go/main/App'
import { ListNotes, GetNote, CreateNote, UpdateNoteTitle,
         UpdateNoteContent, UpdateNotePosition, MoveNote, DeleteNote } from '../wailsjs/go/main/App'
import { SaveImage, GetImage, GetImagePath, ListImagesByNote,
         DeleteImage, DownloadImage } from '../wailsjs/go/main/App'
import { dto } from '../wailsjs/go/models'

// Pinia stores (auto-imported in Vue SFCs)
import { useAuthStore } from '@/stores/auth'
import { useNavigationStore } from '@/stores/navigation'
import type { NavItem, NavLevel } from '@/stores/navigation'

// Composables (auto-imported in Vue SFCs)
import { useTheme } from '@/composables/useTheme'
import { useLocale, availableLocales } from '@/composables/useLocale'
import type { LocaleOption } from '@/composables/useLocale'
import { useNavActions } from '@/composables/useNavActions'
```

## Basic Usage

```typescript
// Authentication flow
const authStore = useAuthStore()
await authStore.init()                        // fetch auth status from backend
if (!authStore.passwordSet) {
  await authStore.setPassword('my-password')  // first run: set password
} else {
  await authStore.verifyPassword('my-password') // subsequent runs
}
if (authStore.isAuthenticated) { /* proceed */ }

// Navigation
const navStore = useNavigationStore()
await navStore.init()                          // load root level
await navStore.openNotebook(notebookListItem)  // drill into notebook
navStore.selectNote(noteId)                    // select note for editing

// Theme
const { currentTheme, setTheme, toggleTheme } = useTheme()
setTheme('dark')   // or 'light'

// Locale
const { currentLocale, setLocale, availableLocales } = useLocale()
setLocale('en')    // switch to English

// Direct backend call
const notebooks = await ListNotebooks('')      // root level (pass "" not null)
const note = await GetNote('note-uuid')
await UpdateNoteContent('note-uuid', { content: '{"type":"doc",...}' })
```

## Capabilities

### Backend Bindings — Wails Go API
Direct RPC calls to the Go backend. All functions return Promises.

```typescript
// Authentication
GetAuthStatus(): Promise<{ isPasswordSet: boolean }>
SetPassword(password: string): Promise<{ success: boolean; error?: string }>
VerifyPassword(password: string): Promise<{ success: boolean }>

// Notebooks
ListNotebooks(parentId: string): Promise<dto.NotebookListItem[]>  // "" = root
GetNotebook(id: string): Promise<dto.NotebookDetail>
CreateNotebook(req: dto.CreateNotebookRequest): Promise<dto.NotebookDetail>
UpdateNotebookTitle(id: string, req: dto.UpdateNotebookTitleRequest): Promise<dto.NotebookDetail>
UpdateNotebookPosition(id: string, req: dto.UpdatePositionRequest): Promise<void>
MoveNotebook(id: string, req: dto.MoveNotebookRequest): Promise<dto.NotebookDetail>
DeleteNotebook(id: string): Promise<void>  // cascades to children

// Notes
ListNotes(notebookId: string): Promise<dto.NoteListItem[]>  // "" = root
GetNote(id: string): Promise<dto.NoteDetail>  // includes TipTap JSON content
CreateNote(req: dto.CreateNoteRequest): Promise<dto.NoteDetail>
UpdateNoteTitle(id: string, req: dto.UpdateNoteTitleRequest): Promise<dto.NoteDetail>
UpdateNoteContent(id: string, req: dto.UpdateNoteContentRequest): Promise<dto.NoteDetail>
UpdateNotePosition(id: string, req: dto.UpdatePositionRequest): Promise<void>
MoveNote(id: string, req: dto.MoveNoteRequest): Promise<dto.NoteDetail>
DeleteNote(id: string): Promise<void>

// Images
SaveImage(noteId: string, mimeType: string, data: number[]): Promise<dto.ImageItem>
GetImage(id: string): Promise<dto.ImageItem>
GetImagePath(filename: string): Promise<string>
ListImagesByNote(noteId: string): Promise<Array<dto.ImageItem>>
DeleteImage(id: string): Promise<void>
DownloadImage(id: string): Promise<void>  // triggers OS save dialog

// Wails Runtime (import from '../wailsjs/runtime/runtime')
BrowserOpenURL(url: string): void  // open URL in system default browser
```

See [Wails Backend API](wails-backend.md) for full reference.

### Pinia Stores — Application State
Reactive state management for auth and navigation.

```typescript
// Auth store
const auth = useAuthStore()
auth.status        // 'authenticated' | 'unauthenticated'
auth.passwordSet   // boolean
auth.isAuthenticated  // computed boolean
await auth.init()

// Navigation store
const nav = useNavigationStore()
nav.stack          // NavLevel[] — current navigation path
nav.leftPanel      // ComputedRef<NavLevel | null>
nav.rightPanel     // ComputedRef<NavLevel | null>
nav.selectedNoteId // string | null
await nav.openNotebook(item)
await nav.reorderItems(level, 'notebook', newItems)
```

See [Pinia Stores](stores.md) for full reference.

### Composables — Reusable Logic

```typescript
// Theme management
const { currentTheme, setTheme, toggleTheme } = useTheme()

// Locale management
const { currentLocale, availableLocales, setLocale } = useLocale()

// Navigation actions (backend + store coordination)
const { createNotebook, createNote, deleteItem } = useNavActions()

// Validation rules (auto-registered, use in vee-validate)
// Rules: 'required', 'min:N', 'confirmed'
```

See [Composables](composables.md) for full reference.

### Data Types — DTOs
All backend request and response types live in the `dto` namespace.

```typescript
import { dto } from '../wailsjs/go/models'

// Key types
dto.NotebookDetail        // { id, parentId?, title, position, createdAt, updatedAt }
dto.NoteDetail            // { id, notebookId?, title, content, position, createdAt, updatedAt }
dto.ImageItem             // { id, noteId, filename, mimeType, createdAt }
dto.CreateNotebookRequest // { parentId?: string, title: string }
dto.CreateNoteRequest     // { notebookId?: string, title: string }
```

See [Data Types](types.md) for full reference.

### Vue Components — Auto-Registered UI
All components are globally available in templates without imports.

```html
<!-- Note editor (loads by noteId, autosaves, emits titleChanged) -->
<NoteEditor :noteId="nav.selectedNoteId" @titleChanged="(id, title) => nav.renameItem(id, 'note', title)" />

<!-- Navigation panels (left/right) -->
<NavPanel :level="nav.leftPanel!" :activeNoteId="nav.selectedNoteId"
  @notebookClick="item => nav.openNotebookFromLeft(item.data)"
  @noteClick="item => nav.selectNoteFromRoot(item.data.id)"
  @reorder="(lvl, kind, items) => nav.reorderItems(lvl, kind, items)"
  @create="kind => openCreateModal(kind)" />

<!-- Modals (v-model:visible pattern) -->
<CreateItemModal v-model:visible="createVisible" :kind="'note'" :level="currentLevel" />
<DeleteConfirmModal v-model:visible="deleteVisible" :kind="item.kind" :title="item.data.title" @confirm="doDelete" />
<MoveItemModal v-model:visible="moveVisible" :item="moveTarget" :level="currentLevel" />
<AppSettingsModal v-model:visible="settingsVisible" />
```

See [Vue Components](components.md) for full props/emits reference.
