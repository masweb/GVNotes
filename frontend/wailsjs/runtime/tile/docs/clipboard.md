# Clipboard

Get and set text content on the system clipboard.

## Import

```typescript
import { ClipboardGetText, ClipboardSetText } from "../wailsjs/runtime/runtime";
```

## Capabilities

### Get and Set Clipboard Text

```typescript { .api }
/**
 * Returns the current text content stored on the system clipboard.
 * @returns Promise resolving to the clipboard text string (empty string if clipboard is empty or contains non-text)
 */
function ClipboardGetText(): Promise<string>;

/**
 * Sets text content on the system clipboard.
 * @param text - The text to write to the clipboard
 * @returns Promise resolving to true on success, false on failure
 */
function ClipboardSetText(text: string): Promise<boolean>;
```

## Usage Examples

```typescript
import { ClipboardGetText, ClipboardSetText } from "../wailsjs/runtime/runtime";

// Read current clipboard content
const text = await ClipboardGetText();
console.log("Clipboard contains:", text);

// Write text to clipboard
const ok = await ClipboardSetText("Hello, clipboard!");
if (ok) {
  LogInfo("Copied to clipboard");
} else {
  LogWarning("Failed to copy to clipboard");
}
```
