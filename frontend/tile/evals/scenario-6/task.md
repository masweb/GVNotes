# Dual-Panel Splitpanes Layout

Build the main application layout using the Splitpanes library with an adaptive sidebar that shows one or two navigation panels depending on navigation depth.

## Capabilities

### Outer layout

- The outer Splitpanes occupies the full viewport height [@test](../test/layout_full_height.test.ts)
- The sidebar pane has an initial size of `25`, a minimum size of `10`, and a maximum size of `60` [@test](../test/layout_sidebar_constraints.test.ts)
- The main content pane has a minimum size of `40` [@test](../test/layout_content_min.test.ts)

### Adaptive sidebar

- When there is no right panel (root navigation depth), only a single navigation panel is rendered inside the sidebar [@test](../test/layout_single_panel.test.ts)
- When a right panel is present, a nested Splitpanes renders two navigation panels side-by-side inside the sidebar [@test](../test/layout_dual_panel.test.ts)
- Each inner pane in the nested split has minimum size `20` and maximum size `80` [@test](../test/layout_inner_constraints.test.ts)

### Theme

- A dark-theme CSS class is applied to the Splitpanes components when the current theme is `'dark'` [@test](../test/layout_dark_class.test.ts)

## Implementation

[@generates](./src/views/GVNotes.vue)

## API

```typescript { #api }
/**
 * GVNotes component props (none — reads from Pinia store and useTheme composable).
 *
 * Expected events emitted by each NavPanel:
 *   @notebook-click  — user clicked a notebook item
 *   @note-click      — user clicked a note item
 *   @reorder         — user reordered items via drag-and-drop
 *   @create          — user triggered item creation
 *   @delete          — user triggered item deletion
 *   @move            — user triggered item move
 *   @rename          — user renamed the panel's notebook
 */
```

## Dependencies { .dependencies }

### frontend 0.0.0 { .dependency }

Vue 3 + TypeScript frontend application for GVNotes. Provides the Splitpanes-based adaptive dual-panel layout pattern with dynamic panel switching, inner pane constraints, and conditional dark-theme class application.
