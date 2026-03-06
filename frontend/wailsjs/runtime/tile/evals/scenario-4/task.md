# Window Dimensions Manager

A module for reading and setting the application window's pixel dimensions.

## Capabilities

### Retrieve Window Dimensions

Fetches the current width and height of the application window as an object with `w` (width) and `h` (height) number properties.

- Calling `getDimensions()` returns a promise that resolves to an object with numeric `w` and `h` properties representing the current window size [@test](./test.test.ts)

### Apply Window Dimensions

Sets the application window to the specified pixel width and height immediately.

- Calling `setDimensions(1280, 720)` sets the window width to `1280` and height to `720` [@test](./test2.test.ts)
- Calling `setDimensions(800, 600)` sets the window width to `800` and height to `600` [@test](./test3.test.ts)

## Implementation

[@generates](./src/index.ts)

## API

```typescript { #api }
export interface Dimensions {
  w: number;
  h: number;
}

export function getDimensions(): Promise<Dimensions>;
export function setDimensions(width: number, height: number): void;
```

## Dependencies { .dependencies }

### @wailsapp/runtime 2.0.0 { .dependency }

Wails JavaScript runtime library providing bridge functions for querying and updating the native application window size.
