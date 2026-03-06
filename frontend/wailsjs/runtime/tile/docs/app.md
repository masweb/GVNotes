# Application Lifecycle

Control the application process and retrieve runtime environment information such as the current platform, CPU architecture, and build type.

## Import

```typescript
import { Environment, Quit, Hide, Show } from "../wailsjs/runtime/runtime";
```

## Types

```typescript { .api }
interface EnvironmentInfo {
  /** Build type: "production" or "development" */
  buildType: string;
  /** Operating system: "darwin", "windows", or "linux" */
  platform: string;
  /** CPU architecture: "amd64", "arm64", etc. */
  arch: string;
}
```

## Capabilities

### Environment

```typescript { .api }
/**
 * Returns information about the current runtime environment.
 * @returns EnvironmentInfo with buildType, platform, and arch
 */
function Environment(): Promise<EnvironmentInfo>;
```

### Application Control

```typescript { .api }
/**
 * Quits the application immediately.
 */
function Quit(): void;

/**
 * Hides the application.
 * On macOS this hides the dock icon and all application windows.
 */
function Hide(): void;

/**
 * Shows the application after it has been hidden via Hide().
 */
function Show(): void;
```

## Usage Examples

```typescript
import { Environment, Quit, Hide, Show } from "../wailsjs/runtime/runtime";

// Detect platform for conditional logic
const env = await Environment();
if (env.platform === "darwin") {
  // macOS-specific behaviour
}
if (env.platform === "windows") {
  // Windows-specific behaviour
}
if (env.buildType === "development") {
  LogDebug("Running in dev mode");
}
console.log("Architecture:", env.arch); // e.g. "amd64"

// Graceful quit triggered from the frontend
Quit();

// Hide/show the app (macOS: removes from dock while keeping process alive)
Hide();
// ... later
Show();
```
