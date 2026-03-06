# Drag-and-Drop Item Reordering with Position Persistence

Implement drag-and-drop reordering of navigation items using vue-draggable-plus, and implement the corresponding store action that updates reactive state and persists new positions to the backend.

## Capabilities

### Store reorder action

- After reordering notes `[A, B, C]` to `[C, A, B]`, the resulting persisted positions are `0`, `1000`, and `2000` respectively [@test](../test/reorder_positions.test.ts)
- The reorder action keeps notebooks before notes in the items array regardless of the order items are passed in [@test](../test/reorder_kind_separation.test.ts)
- The backend position-update function is called exactly once per item in the reordered list [@test](../test/reorder_backend_calls.test.ts)

### Vue template integration

- The sortable list component uses a model binding that reflects the current items and emits the updated array on change [@test](../test/reorder_draggable_model.test.ts)
- The change handler receives the affected level, the item kind, and the new sorted array [@test](../test/reorder_handler_signature.test.ts)

## Implementation

[@generates](./src/stores/navigation.ts)

## API

```typescript { #api }
/**
 * Reorders items of a single kind within a navigation level.
 * Updates the reactive stack entry and persists each new position
 * to the backend using an index * 1000 gap-based strategy.
 *
 * @param level    The NavLevel whose items are being reordered.
 * @param kind     Whether notebooks or notes are being reordered.
 * @param newItems The items of the given kind in their new order.
 */
declare function reorderItems(
  level: NavLevel,
  kind: 'notebook' | 'note',
  newItems: NavItem[],
): Promise<void>
```

## Dependencies { .dependencies }

### frontend 0.0.0 { .dependency }

Vue 3 + TypeScript frontend application for GVNotes. Provides the drag-and-drop reorder pattern using vue-draggable-plus VueDraggable component with kind-separation logic and gap-based (index × 1000) position persistence.
