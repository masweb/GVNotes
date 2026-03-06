# Environment Inspector

A module that retrieves information about the current runtime environment of the Wails application.

## Capabilities

### Query Runtime Environment

Fetches environment details including the application's build type, the operating system platform, and the CPU architecture. This information is provided by the Go backend at runtime.

- Calling `getEnvironment()` returns a promise that resolves to an object with `buildType`, `platform`, and `arch` string properties [@test](./test.test.ts)
- The resolved `buildType` property contains either `"debug"` (development build) or `"production"` (release build) [@test](./test2.test.ts)
- The resolved `platform` property identifies the operating system (e.g., `"windows"`, `"linux"`, `"darwin"`) [@test](./test3.test.ts)

## Implementation

[@generates](./src/index.ts)

## API

```typescript { #api }
export interface RuntimeEnvironment {
  buildType: string;
  platform: string;
  arch: string;
}

export function getEnvironment(): Promise<RuntimeEnvironment>;
```

## Dependencies { .dependencies }

### @wailsapp/runtime 2.0.0 { .dependency }

Wails JavaScript runtime library providing the bridge function for querying build type, platform, and CPU architecture from the Go backend.
