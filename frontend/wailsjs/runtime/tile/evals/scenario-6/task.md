# Clipboard Manager

A module for reading from and writing to the system clipboard.

## Capabilities

### Read Clipboard Text

Retrieves the current text content stored on the system clipboard.

- Calling `readText()` returns a promise that resolves to the current clipboard text string [@test](./test.test.ts)

### Write Clipboard Text

Stores a string on the system clipboard and reports whether the operation succeeded.

- Calling `writeText("hello world")` writes the string `"hello world"` to the clipboard and returns a promise that resolves to a boolean indicating success [@test](./test2.test.ts)
- Calling `writeText("")` writes an empty string to the clipboard and returns a promise resolving to a boolean [@test](./test3.test.ts)

## Implementation

[@generates](./src/index.ts)

## API

```typescript { #api }
export function readText(): Promise<string>;
export function writeText(text: string): Promise<boolean>;
```

## Dependencies { .dependencies }

### @wailsapp/runtime 2.0.0 { .dependency }

Wails JavaScript runtime library providing bridge functions for reading from and writing to the native system clipboard.
