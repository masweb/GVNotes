# @wailsapp/runtime

The Wails v2 JavaScript/TypeScript runtime library is the frontend bridge for Wails desktop applications. It wraps the native runtime injected by Wails into the WebView (`window.runtime`) and exposes it as a clean ES module API for bidirectional communication between JavaScript and Go.

## Package Information

- **Package Name**: `@wailsapp/runtime`
- **Package Type**: npm
- **Language**: JavaScript / TypeScript
- **Version**: 2.0.0
- **License**: MIT
- **Installation**: Embedded automatically in Wails projects at `frontend/wailsjs/runtime/`. Add as a local dependency:
  ```json
  { "@wailsapp/runtime": "./wailsjs/runtime" }
  ```

## Core Imports

```typescript
import {
  EventsOn, EventsOnce, EventsOnMultiple, EventsOff, EventsOffAll, EventsEmit,
  LogPrint, LogTrace, LogDebug, LogInfo, LogWarning, LogError, LogFatal,
  WindowSetTitle, WindowSetSize, WindowGetSize, WindowCenter,
  WindowFullscreen, WindowUnfullscreen, WindowIsFullscreen,
  WindowMaximise, WindowUnmaximise, WindowIsMaximised,
  WindowToggleMaximise, WindowMinimise, WindowUnminimise, WindowIsMinimised,
  WindowIsNormal, WindowHide, WindowShow, WindowReload, WindowReloadApp,
  WindowSetPosition, WindowGetPosition, WindowSetMaxSize, WindowSetMinSize,
  WindowSetAlwaysOnTop, WindowSetBackgroundColour,
  WindowSetSystemDefaultTheme, WindowSetLightTheme, WindowSetDarkTheme,
  ScreenGetAll, BrowserOpenURL,
  Environment, Quit, Hide, Show,
  ClipboardGetText, ClipboardSetText,
  OnFileDrop, OnFileDropOff, CanResolveFilePaths, ResolveFilePaths,
} from "../wailsjs/runtime/runtime";
```

## Basic Usage

```typescript
import { EventsOn, EventsEmit, WindowSetTitle, Quit, Environment } from "../wailsjs/runtime/runtime";

// Listen for events from Go
const unlisten = EventsOn("data-updated", (payload: any) => {
  console.log("Received:", payload);
});
// Remove listener when done
unlisten();

// Emit an event to Go
EventsEmit("frontend-ready", { version: "1.0" });

// Control the window
WindowSetTitle("My Wails App");

// Get runtime environment
const env = await Environment();
// { buildType: "production", platform: "darwin", arch: "amd64" }

// Quit the application
Quit();
```

## Capabilities

### Events

Bidirectional event bus between JavaScript frontend and Go backend. Subscribe, listen once, limit calls, unsubscribe, and emit events with optional data payloads.

```typescript { .api }
function EventsOn(eventName: string, callback: (...data: any) => void): () => void;
function EventsEmit(eventName: string, ...data: any): void;
```

[Events](./events.md)

### Window Management

Control the application window — position, size, state (fullscreen/maximise/minimise), title, visibility, background colour, always-on-top, and theme (Windows only).

```typescript { .api }
function WindowSetTitle(title: string): void;
function WindowSetSize(width: number, height: number): void;
function WindowGetSize(): Promise<Size>;
function WindowCenter(): void;
function WindowMaximise(): void;
function WindowIsMaximised(): Promise<boolean>;
```

[Window Management](./window.md)

### Logging

Send log messages from the JavaScript frontend to the Go backend at various severity levels (trace, debug, info, warning, error, fatal, raw).

```typescript { .api }
function LogInfo(message: string): void;
function LogError(message: string): void;
function LogFatal(message: string): void;
```

[Logging](./logging.md)

### Application Lifecycle

Control the application process and retrieve environment information such as platform, architecture, and build type.

```typescript { .api }
function Environment(): Promise<EnvironmentInfo>;
function Quit(): void;
function Hide(): void;
function Show(): void;
```

[Application Lifecycle](./app.md)

### Clipboard

Get and set text on the system clipboard.

```typescript { .api }
function ClipboardGetText(): Promise<string>;
function ClipboardSetText(text: string): Promise<boolean>;
```

[Clipboard](./clipboard.md)

### Screen

Retrieve information about all connected monitors including dimensions, current screen, and primary screen status.

```typescript { .api }
function ScreenGetAll(): Promise<Screen[]>;
```

[Screen](./screen.md)

### Browser

Open URLs in the system's default browser.

```typescript { .api }
function BrowserOpenURL(url: string): void;
```

[Browser](./browser.md)

### Drag and Drop

Handle file drag-and-drop events and resolve native file paths from dropped browser `File` objects.

```typescript { .api }
function OnFileDrop(callback: (x: number, y: number, paths: string[]) => void, useDropTarget: boolean): void;
function OnFileDropOff(): void;
function CanResolveFilePaths(): boolean;
function ResolveFilePaths(files: File[]): void;
```

[Drag and Drop](./drag-drop.md)

## Types

```typescript { .api }
interface Position {
  x: number;
  y: number;
}

interface Size {
  w: number;
  h: number;
}

interface Screen {
  isCurrent: boolean;
  isPrimary: boolean;
  width: number;
  height: number;
}

interface EnvironmentInfo {
  /** "production" or "development" */
  buildType: string;
  /** Operating system: "darwin", "windows", "linux" */
  platform: string;
  /** CPU architecture: "amd64", "arm64", etc. */
  arch: string;
}
```
