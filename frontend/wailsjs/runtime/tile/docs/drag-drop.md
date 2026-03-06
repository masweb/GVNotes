# Drag and Drop

Handle file drag-and-drop events and resolve native filesystem paths from dropped browser `File` objects.

## Import

```typescript
import {
  OnFileDrop,
  OnFileDropOff,
  CanResolveFilePaths,
  ResolveFilePaths,
} from "../wailsjs/runtime/runtime";
```

## Capabilities

### File Drop Listener

```typescript { .api }
/**
 * Sets up a listener for drag-and-drop file events.
 * The callback receives the drop coordinates and an array of native file path strings.
 *
 * @param callback - Called when a file drop completes.
 *   - x: horizontal coordinate of the drop position
 *   - y: vertical coordinate of the drop position
 *   - paths: array of native filesystem path strings for the dropped files
 * @param useDropTarget - When true, the callback only fires if the drop landed
 *   on an element styled with the `--wails-drop-target` CSS property.
 *   When false, any drop in the window triggers the callback.
 */
function OnFileDrop(
  callback: (x: number, y: number, paths: string[]) => void,
  useDropTarget: boolean
): void;

/**
 * Removes all drag-and-drop listeners and handlers registered via OnFileDrop.
 */
function OnFileDropOff(): void;
```

### Path Resolution

```typescript { .api }
/**
 * Returns true if the native file path resolver is available on this platform.
 * Check this before calling ResolveFilePaths.
 */
function CanResolveFilePaths(): boolean;

/**
 * Resolves native filesystem paths for an array of browser File objects.
 * Required on platforms where File.path is not natively accessible.
 * Call CanResolveFilePaths() first to verify availability.
 *
 * @param files - Array of browser File objects to resolve paths for
 */
function ResolveFilePaths(files: File[]): void;
```

## Usage Examples

```typescript
import { OnFileDrop, OnFileDropOff, CanResolveFilePaths, ResolveFilePaths } from "../wailsjs/runtime/runtime";

// Accept drops anywhere in the window
OnFileDrop((x, y, paths) => {
  console.log(`Dropped ${paths.length} file(s) at (${x}, ${y})`);
  for (const path of paths) {
    console.log(" -", path);
  }
}, false);

// Accept drops only on elements with --wails-drop-target CSS styling
OnFileDrop((x, y, paths) => {
  handleDroppedFiles(paths);
}, true);

// Remove the listener when no longer needed
OnFileDropOff();

// Resolve paths from browser File objects (e.g., from a file input or drag event)
if (CanResolveFilePaths()) {
  ResolveFilePaths(Array.from(inputElement.files));
}
```
