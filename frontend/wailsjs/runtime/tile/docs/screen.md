# Screen

Retrieve information about all connected monitors, including their dimensions, and which screen is current and primary.

## Import

```typescript
import { ScreenGetAll } from "../wailsjs/runtime/runtime";
```

## Types

```typescript { .api }
interface Screen {
  /** true if this is the screen the application window is currently displayed on */
  isCurrent: boolean;
  /** true if this is the primary/main system screen */
  isPrimary: boolean;
  /** Screen width in pixels */
  width: number;
  /** Screen height in pixels */
  height: number;
}
```

## Capabilities

### Get All Screens

```typescript { .api }
/**
 * Returns an array of all connected screens with their dimensions and status.
 * Call this function anew each time you need fresh data — the underlying windowing
 * system is queried on every call.
 * @returns Promise resolving to an array of Screen objects
 */
function ScreenGetAll(): Promise<Screen[]>;
```

## Usage Examples

```typescript
import { ScreenGetAll } from "../wailsjs/runtime/runtime";

const screens = await ScreenGetAll();

// Find the primary monitor
const primary = screens.find((s) => s.isPrimary);
console.log("Primary screen:", primary?.width, "x", primary?.height);

// Find the screen where the window currently resides
const current = screens.find((s) => s.isCurrent);
console.log("Current screen:", current?.width, "x", current?.height);

// Log all screens
screens.forEach((s, i) => {
  console.log(`Screen ${i}: ${s.width}x${s.height} primary=${s.isPrimary} current=${s.isCurrent}`);
});
```
