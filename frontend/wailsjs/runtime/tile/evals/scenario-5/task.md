# Fullscreen Controller

A module that controls and queries the fullscreen state of the application window.

## Capabilities

### Enter and Exit Fullscreen

Switches the window into fullscreen mode or restores it from fullscreen to its previous dimensions and position.

- Calling `enterFullscreen()` switches the window to fullscreen mode [@test](./test.test.ts)
- Calling `exitFullscreen()` restores the window from fullscreen to the dimensions and position it had before going fullscreen [@test](./test2.test.ts)

### Query Fullscreen State

Returns the current fullscreen state of the window.

- Calling `isFullscreen()` returns a promise that resolves to `true` when the window is in fullscreen mode and `false` otherwise [@test](./test3.test.ts)

## Implementation

[@generates](./src/index.ts)

## API

```typescript { #api }
export function enterFullscreen(): void;
export function exitFullscreen(): void;
export function isFullscreen(): Promise<boolean>;
```

## Dependencies { .dependencies }

### @wailsapp/runtime 2.0.0 { .dependency }

Wails JavaScript runtime library providing bridge functions for controlling and querying the native application window's fullscreen state.
