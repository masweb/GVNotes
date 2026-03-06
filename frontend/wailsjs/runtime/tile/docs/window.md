# Window Management

Control the application window's position, size, state, title, visibility, background colour, always-on-top behaviour, and theme.

## Import

```typescript
import {
  WindowReload,
  WindowReloadApp,
  WindowSetAlwaysOnTop,
  WindowCenter,
  WindowSetTitle,
  WindowFullscreen,
  WindowUnfullscreen,
  WindowIsFullscreen,
  WindowSetSize,
  WindowGetSize,
  WindowSetMaxSize,
  WindowSetMinSize,
  WindowSetPosition,
  WindowGetPosition,
  WindowHide,
  WindowShow,
  WindowMaximise,
  WindowToggleMaximise,
  WindowUnmaximise,
  WindowIsMaximised,
  WindowMinimise,
  WindowUnminimise,
  WindowIsMinimised,
  WindowIsNormal,
  WindowSetBackgroundColour,
  WindowSetSystemDefaultTheme,
  WindowSetLightTheme,
  WindowSetDarkTheme,
} from "../wailsjs/runtime/runtime";
```

## Types

```typescript { .api }
interface Size {
  w: number;
  h: number;
}

interface Position {
  x: number;
  y: number;
}
```

## Capabilities

### Reload

```typescript { .api }
/**
 * Forces a reload by the main application as well as any connected browsers.
 */
function WindowReload(): void;

/**
 * Reloads the application frontend only (not connected browsers).
 */
function WindowReloadApp(): void;
```

### Title

```typescript { .api }
/**
 * Sets the text displayed in the window title bar.
 * @param title - The new window title
 */
function WindowSetTitle(title: string): void;
```

### Position

```typescript { .api }
/**
 * Centers the window on the monitor it is currently displayed on.
 */
function WindowCenter(): void;

/**
 * Sets the window position relative to the monitor the window is currently on.
 * @param x - Horizontal offset in pixels from the left of the monitor
 * @param y - Vertical offset in pixels from the top of the monitor
 */
function WindowSetPosition(x: number, y: number): void;

/**
 * Gets the window position relative to the monitor the window is currently on.
 * @returns Promise resolving to { x, y } position in pixels
 */
function WindowGetPosition(): Promise<Position>;
```

### Size

```typescript { .api }
/**
 * Sets the window width and height in pixels.
 * @param width - Window width in pixels
 * @param height - Window height in pixels
 */
function WindowSetSize(width: number, height: number): void;

/**
 * Gets the current window width and height in pixels.
 * @returns Promise resolving to { w, h } dimensions
 */
function WindowGetSize(): Promise<Size>;

/**
 * Sets the maximum window size. Resizes the window if it is currently larger.
 * Pass 0, 0 to disable the constraint.
 * @param width - Maximum width in pixels (0 = no limit)
 * @param height - Maximum height in pixels (0 = no limit)
 */
function WindowSetMaxSize(width: number, height: number): void;

/**
 * Sets the minimum window size. Resizes the window if it is currently smaller.
 * Pass 0, 0 to disable the constraint.
 * @param width - Minimum width in pixels (0 = no limit)
 * @param height - Minimum height in pixels (0 = no limit)
 */
function WindowSetMinSize(width: number, height: number): void;
```

### Fullscreen

```typescript { .api }
/**
 * Makes the window fullscreen.
 */
function WindowFullscreen(): void;

/**
 * Restores the window to the dimensions and position it had before going fullscreen.
 */
function WindowUnfullscreen(): void;

/**
 * Returns true if the window is currently in fullscreen mode.
 */
function WindowIsFullscreen(): Promise<boolean>;
```

### Maximise

```typescript { .api }
/**
 * Maximises the window to fill the screen.
 */
function WindowMaximise(): void;

/**
 * Toggles between maximised and unmaximised states.
 */
function WindowToggleMaximise(): void;

/**
 * Restores the window to the dimensions and position it had before maximising.
 */
function WindowUnmaximise(): void;

/**
 * Returns true if the window is currently maximised.
 */
function WindowIsMaximised(): Promise<boolean>;
```

### Minimise

```typescript { .api }
/**
 * Minimises the window to the taskbar/dock.
 */
function WindowMinimise(): void;

/**
 * Restores the window to the dimensions and position it had before minimising.
 */
function WindowUnminimise(): void;

/**
 * Returns true if the window is currently minimised.
 */
function WindowIsMinimised(): Promise<boolean>;

/**
 * Returns true if the window is in its normal state
 * (not maximised, not minimised, not fullscreen).
 */
function WindowIsNormal(): Promise<boolean>;
```

### Visibility

```typescript { .api }
/**
 * Hides the window without quitting the application.
 */
function WindowHide(): void;

/**
 * Shows the window if it is currently hidden.
 */
function WindowShow(): void;

/**
 * Sets the window as always-on-top when b is true, removes the constraint when false.
 * @param b - true to pin the window above all others, false to restore normal stacking
 */
function WindowSetAlwaysOnTop(b: boolean): void;
```

### Background Colour

```typescript { .api }
/**
 * Sets the background colour of the window. This colour shows through all transparent pixels.
 * @param R - Red channel, 0–255
 * @param G - Green channel, 0–255
 * @param B - Blue channel, 0–255
 * @param A - Alpha channel, 0–255 (0 = fully transparent, 255 = fully opaque)
 */
function WindowSetBackgroundColour(R: number, G: number, B: number, A: number): void;
```

### Theme (Windows only)

```typescript { .api }
/**
 * Sets the window theme to the system default (follows OS dark/light setting).
 * Windows only.
 */
function WindowSetSystemDefaultTheme(): void;

/**
 * Sets the window to light theme.
 * Windows only.
 */
function WindowSetLightTheme(): void;

/**
 * Sets the window to dark theme.
 * Windows only.
 */
function WindowSetDarkTheme(): void;
```

## Usage Examples

```typescript
import {
  WindowCenter, WindowSetSize, WindowSetMinSize, WindowSetMaxSize,
  WindowGetSize, WindowGetPosition, WindowSetPosition,
  WindowFullscreen, WindowUnfullscreen, WindowIsFullscreen,
  WindowMaximise, WindowIsMaximised, WindowMinimise,
  WindowSetTitle, WindowSetAlwaysOnTop, WindowSetBackgroundColour,
  WindowHide, WindowShow,
} from "../wailsjs/runtime/runtime";

// Initial window setup
WindowSetTitle("My App");
WindowCenter();
WindowSetSize(1280, 800);
WindowSetMinSize(800, 600);
WindowSetMaxSize(1920, 1080); // Constrain to 1080p max (0,0 to remove constraint)

// Read current dimensions and position
const size = await WindowGetSize();     // { w: 1280, h: 800 }
const pos = await WindowGetPosition();  // { x: ..., y: ... }

// Toggle fullscreen
const isFS = await WindowIsFullscreen();
if (isFS) {
  WindowUnfullscreen();
} else {
  WindowFullscreen();
}

// Maximise if not already maximised
if (!(await WindowIsMaximised())) {
  WindowMaximise();
}

// Semi-transparent background (shows through transparent CSS areas)
WindowSetBackgroundColour(0, 0, 0, 128);

// Pin window on top of all others
WindowSetAlwaysOnTop(true);

// Hide window temporarily (process keeps running)
WindowHide();
// ... later
WindowShow();
```
