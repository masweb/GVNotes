# Move Item in Navigation Store

Implement a `moveItem` action inside a hierarchical navigation store that relocates a notebook or note to a different parent, operating entirely on in-memory reactive state.

## Capabilities

### Source removal

- After calling `moveItem`, the moved item is no longer present in the source level's items array [@test](../test/move_removes_from_source.test.ts)

### Destination insertion

- When the destination level is currently loaded in the navigation stack, the moved item is appended to that level's items array [@test](../test/move_inserts_at_destination.test.ts)
- When the destination level is not loaded in the stack, the item is removed from the source but not inserted anywhere [@test](../test/move_destination_not_loaded.test.ts)

### Stack collapse for open notebooks

- When a notebook being moved is currently open in the navigation stack (a level whose parent identifier equals the notebook's id exists), the stack is collapsed to remove that level and all levels beyond it [@test](../test/move_collapses_open_notebook.test.ts)
- When a note is moved, the stack is never collapsed regardless of navigation depth [@test](../test/move_note_no_collapse.test.ts)
- When a notebook is moved but it is not currently open in the stack, the stack is not collapsed [@test](../test/move_notebook_not_open_no_collapse.test.ts)

## Implementation

[@generates](./src/stores/navigation.ts)

## API

```typescript { #api }
/**
 * Move an item from one parent to another within the navigation stack.
 * Operates only on in-memory reactive state; callers must handle backend
 * persistence separately.
 *
 * @param sourceLevelParentId  parentId of the NavLevel that currently holds the item
 * @param item                 The NavItem to move
 * @param destinationParentId  parentId of the target NavLevel (null = root)
 */
declare function moveItem(
  sourceLevelParentId: string | null,
  item: NavItem,
  destinationParentId: string | null,
): void
```

## Dependencies { .dependencies }

### frontend 0.0.0 { .dependency }

Vue 3 + TypeScript frontend application for GVNotes. Provides the moveItem navigation store action pattern with three-step logic: source removal, conditional destination insertion, and stack collapse when moving an open notebook.
