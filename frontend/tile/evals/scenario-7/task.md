# Navigation CRUD Actions Composable

Implement a composable that provides create and delete operations for notebooks and notes, integrating async backend calls with immediate reactive navigation store updates.

## Capabilities

### Create notebook

- After creating a notebook, a typed item is inserted into the matching navigation level (found by its parent identifier) before any navigation action is called [@test](../test/nav_actions_create_notebook_insert.test.ts)
- After inserting the item, the store's open-notebook action is called with the new notebook's data [@test](../test/nav_actions_create_notebook_navigate.test.ts)

### Create note

- After creating a note, a typed item is inserted into the matching navigation level [@test](../test/nav_actions_create_note_insert.test.ts)
- After inserting the item, the store's select-note action is called with the new note's identifier [@test](../test/nav_actions_create_note_select.test.ts)

### Delete item

- When the item kind is `'notebook'`, the notebook backend deletion function is called with the item's identifier [@test](../test/nav_actions_delete_notebook.test.ts)
- When the item kind is `'note'`, the note backend deletion function is called with the item's identifier [@test](../test/nav_actions_delete_note.test.ts)
- After a successful deletion, the store's remove-item action is called [@test](../test/nav_actions_delete_remove.test.ts)

## Implementation

[@generates](./src/composables/useNavActions.ts)

## API

```typescript { #api }
import type { NavItem, NavLevel } from '@/stores/navigation'

export declare function useNavActions(): {
  /**
   * Create a notebook inside the given level, insert it into the store,
   * and navigate into it.
   */
  createNotebook: (level: NavLevel, title: string) => Promise<void>
  /**
   * Create a note inside the given level, insert it into the store,
   * and select it.
   */
  createNote: (level: NavLevel, title: string) => Promise<void>
  /**
   * Delete a notebook or note, removing it from the backend and the store.
   */
  deleteItem: (level: NavLevel, item: NavItem) => Promise<void>
}
```

## Dependencies { .dependencies }

### frontend 0.0.0 { .dependency }

Vue 3 + TypeScript frontend application for GVNotes. Provides the useNavActions composable pattern that couples Wails backend calls (CreateNotebook, CreateNote, DeleteNotebook, DeleteNote) with immediate Pinia store mutations for instant UI feedback.
