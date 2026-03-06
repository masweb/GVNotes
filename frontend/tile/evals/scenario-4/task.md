# Stack-Based Hierarchical Navigation Store

Implement a Pinia store for hierarchical notebook/note navigation using a stack of navigation levels.

A **navigation level** has: a nullable parent identifier, a display title, and an array of typed items (each item is either a notebook or a note, identified by a `kind` field).

## Capabilities

### Computed panel state

- After `init()` loads a single root level, `leftPanel` returns that level and `rightPanel` is `null` [@test](../test/nav_init_panels.test.ts)
- After `openNotebook()` pushes a second level, `rightPanel` is the new level and `leftPanel` is the root [@test](../test/nav_open_notebook_panels.test.ts)
- `canGoBack` is `false` at root depth and `true` when depth > 1 [@test](../test/nav_can_go_back.test.ts)
- `activeNotebookId` is `null` at root depth and equals the last level's parent identifier when depth > 1 [@test](../test/nav_active_notebook_id.test.ts)

### Navigation actions

- `openNotebook()` appends a new level to the stack and clears the selected note [@test](../test/nav_open_notebook_push.test.ts)
- `openNotebookFromLeft()` replaces the last stack entry instead of pushing a new one [@test](../test/nav_open_from_left_replace.test.ts)
- `goBack()` removes the last level from the stack and clears the selected note [@test](../test/nav_go_back.test.ts)
- `selectNoteFromRoot()` collapses the stack to its first entry before selecting a note [@test](../test/nav_select_from_root.test.ts)

## Implementation

[@generates](./src/stores/navigation.ts)

## API

```typescript { #api }
import { ComputedRef, Ref } from 'vue'

export interface NavItem {
  kind: 'notebook' | 'note'
  data: { id: string; title: string; position: number }
}

export interface NavLevel {
  parentId: string | null
  title: string
  items: NavItem[]
}

export declare const useNavigationStore: () => {
  stack: Ref<NavLevel[]>
  selectedNoteId: Ref<string | null>
  loading: Ref<boolean>
  leftPanel: ComputedRef<NavLevel | null>
  rightPanel: ComputedRef<NavLevel | null>
  canGoBack: ComputedRef<boolean>
  activeNotebookId: ComputedRef<string | null>
  init: () => Promise<void>
  openNotebook: (notebook: { id: string; title: string }) => Promise<void>
  openNotebookFromLeft: (notebook: { id: string; title: string }) => Promise<void>
  selectNote: (noteId: string) => void
  selectNoteFromRoot: (noteId: string) => void
  goBack: () => void
}
```

## Dependencies { .dependencies }

### frontend 0.0.0 { .dependency }

Vue 3 + TypeScript frontend application for GVNotes. Provides the stack-based navigation store pattern with leftPanel/rightPanel computed properties, push vs replace navigation semantics, and root-collapse behaviour.
