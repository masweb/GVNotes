# GVNotes

A fast, private, offline-first note-taking app for desktop. No accounts, no cloud, no subscriptions — just your notes, stored locally.

---

## Features

### Organization
- **Hierarchical notebooks** — create nested notebooks to organize your notes any way you like
- **Drag & drop reordering** — reorder notes and notebooks within any level by dragging
- **Move notes and notebooks** — move any item to a different notebook or back to root using a tree picker modal
- **Inline rename** — click any notebook title to rename it in place
- **Quick delete** — remove notes and notebooks with a confirmation modal (cascades to all contents)

### Editor
- **Rich text editing** powered by TipTap — bold, italic, underline, strikethrough, inline code, text color
- **Headings** — H1 through H6 with a toolbar dropdown
- **Lists** — bullet lists, ordered lists, and nested task lists with clickable checkboxes
- **Tables** — insert and resize tables with add/remove row and column controls
- **Code blocks** — syntax-highlighted code blocks with a language selector (JavaScript, TypeScript, Go, Python, Rust, SQL, and more)
- **Images** — paste or drag images directly into notes; resize by dragging corners; download with one click
- **Links** — insert, edit, and open hyperlinks
- **Auto-save** — content is saved automatically as you type

### Navigation & UX
- **Dual-panel navigation** — browse a parent notebook and a child notebook side by side
- **Search** — filter notes and notebooks instantly within any panel
- **Keyboard shortcuts** — `⌘N` to create a new note, `⌘NN` to create a new notebook
- **Dark / Light theme** — toggle between themes; respects system preference on first launch
- **Multilingual** — interface available in English and Spanish

### Privacy & Security
- **100% offline** — no network requests, no telemetry, no external services
- **Local storage** — all data stored in a SQLite database in your system's standard data directory
- **App password** — optional single password (Argon2id) to lock the app

---

## Platform Support

| Platform | Status |
|----------|--------|
| macOS    | ✅ Supported |
| Windows  | 🔜 Coming soon |
| Linux    | 🔜 Coming soon |

---

## Technical Overview

### Architecture

GVNotes is a desktop application built with **Wails v2**, combining a **Go** backend with a **Vue 3** frontend rendered in a native WebView. There is no Electron, no Node.js runtime, and no browser process — the binary is self-contained.

```
main.go              # Entry point — embeds frontend/dist and db/migrations
app.go               # App facade: Wails bindings + dependency injection
internal/
  config/            # Database path (XDG), image directory
  controllers/       # HTTP-style handlers: auth, notebook, note, image
  services/          # Business logic interfaces + implementations
  repositories/      # Extra query logic beyond sqlc (when needed)
  dto/               # Request/response data transfer objects
  errors/            # Typed sentinel errors (ErrNotFound, ErrConflict…)
db/
  migrations/        # golang-migrate SQL files (up/down)
  queries/           # sqlc input: raw SQL queries
  generated/         # sqlc output: type-safe Go query functions (do not edit)
  sqlc.yaml          # sqlc configuration
frontend/
  src/
    components/      # Vue SFC components
    views/           # Top-level views (GVNotes.vue)
    stores/          # Pinia stores (navigation)
    composables/     # Reusable Vue composition functions
    locales/         # i18n JSON files (en, es)
    css/             # Global SCSS (main.scss)
```

### Backend stack

| Technology | Role |
|------------|------|
| Go 1.24 | Application runtime |
| Wails v2 | Native desktop shell + JS bridge |
| SQLite (mattn/go-sqlite3, CGO) | Embedded database |
| sqlc | Compile-time type-safe SQL → Go |
| golang-migrate | Schema migrations (embedded in binary) |
| adrg/xdg | Cross-platform data directory resolution |
| Argon2id | Password hashing (OWASP parameters) |

**Data locations (macOS):**
- Database: `~/Library/Application Support/gvnotes/gvnotes.db`
- Images: `~/Library/Application Support/gvnotes/images/<uuid>.<ext>`

**DI flow:** `OpenDB → db.New → NewXxxService(querier) → NewXxxController(svc)` — no framework, plain constructor injection wired in `app.startup()`.

### Frontend stack

| Technology | Role |
|------------|------|
| Vue 3 + TypeScript | UI framework |
| Pinia | Reactive state management |
| TipTap v2 | Rich text editor (ProseMirror-based) |
| lowlight + highlight.js | Syntax highlighting in code blocks |
| CoreUI + Bootstrap 5 | Component library and utility classes |
| Vue I18n | Internationalization |
| Sortable.js | Drag & drop reordering |
| Splitpanes | Resizable split panels |
| Vite | Build tool and dev server |
| vee-validate | Form validation |

### Database schema

```sql
settings   (key PK, value)                          -- key/value store (e.g. password hash)
notebooks  (id, parent_id FK self-ref, title, position, created_at, updated_at)
notes      (id, notebook_id FK, title, content JSON, position, created_at, updated_at)
images     (id, note_id FK, filename, mime_type, created_at)
```

- `NULL parent_id` / `NULL notebook_id` = root level
- `position` is an integer per-type within a parent, multiplied by 1000 to allow gap-based reordering
- Note content is stored as TipTap JSON
- Images are stored as files on disk; only metadata in the DB. File is written first, DB insert after (file removed on failure)

### Development

```bash
# Prerequisites: Go 1.24+, Node.js, pnpm, Wails CLI
# macOS also requires Xcode Command Line Tools (CGO for SQLite)

wails dev          # Start with hot reload (Vite + Go)
wails build        # Build production binary

# Regenerate DB query code after editing db/queries/*.sql
cd db && sqlc generate

# Seed the database with sample data
go run seed.go          # Add sample notebooks and notes
go run seed.go --reset  # Reset and reseed
```
