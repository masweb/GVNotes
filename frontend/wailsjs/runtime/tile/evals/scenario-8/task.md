# Drop Zone Manager

A module that enables and disables drag-and-drop file handling for the application window.

## Capabilities

### Register File Drop Handler

Attaches a listener that is invoked whenever files are dragged and dropped onto the window. The handler receives the pixel coordinates of the drop point and an array of file path strings. An optional flag restricts the handler to only fire when the drop lands on an element that has been marked as a drop target (via the `--wails-drop-target` CSS custom property).

- Calling `enableDrop(handler, false)` registers `handler` to fire on any file drop anywhere in the window, with the correct `(x, y, paths)` arguments [@test](./test.test.ts)
- Calling `enableDrop(handler, true)` restricts the handler to only fire when the drop finishes on a designated drop-target element [@test](./test2.test.ts)

### Remove File Drop Handler

Removes all drag-and-drop listeners and handlers previously registered on the window.

- Calling `disableDrop()` removes all drag-and-drop listeners and handlers [@test](./test3.test.ts)

## Implementation

[@generates](./src/index.ts)

## API

```typescript { #api }
export type DropHandler = (x: number, y: number, paths: string[]) => void;

export function enableDrop(handler: DropHandler, dropTargetOnly: boolean): void;
export function disableDrop(): void;
```

## Dependencies { .dependencies }

### @wailsapp/runtime 2.0.0 { .dependency }

Wails JavaScript runtime library providing bridge functions for registering and removing native file drag-and-drop event handlers.
