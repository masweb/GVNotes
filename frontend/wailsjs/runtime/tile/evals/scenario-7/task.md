# Screen Inspector

A module that retrieves live information about all connected monitors.

## Capabilities

### Fetch All Screen Details

Returns a list of all screens detected by the windowing system at the time of the call. Each entry in the list describes one monitor and includes whether it currently contains the application window, whether it is the primary display, and its pixel dimensions.

- Calling `getScreens()` returns a promise that resolves to an array of screen information objects [@test](./test.test.ts)
- Each object in the resolved array has `isCurrent` (boolean), `isPrimary` (boolean), `width` (number), and `height` (number) properties [@test](./test2.test.ts)
- The data returned by `getScreens()` reflects the live state of the connected displays at the moment of the call [@test](./test3.test.ts)

## Implementation

[@generates](./src/index.ts)

## API

```typescript { #api }
export interface ScreenInfo {
  isCurrent: boolean;
  isPrimary: boolean;
  width: number;
  height: number;
}

export function getScreens(): Promise<ScreenInfo[]>;
```

## Dependencies { .dependencies }

### @wailsapp/runtime 2.0.0 { .dependency }

Wails JavaScript runtime library providing the bridge function for querying all connected monitors from the native windowing system.
