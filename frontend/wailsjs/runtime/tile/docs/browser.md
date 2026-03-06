# Browser

Open URLs in the system's default browser from within a Wails desktop application.

## Import

```typescript
import { BrowserOpenURL } from "../wailsjs/runtime/runtime";
```

## Capabilities

### Open URL

```typescript { .api }
/**
 * Opens the given URL in the system default browser.
 * @param url - The URL to open (must be a fully-formed URL including scheme)
 */
function BrowserOpenURL(url: string): void;
```

## Usage Examples

```typescript
import { BrowserOpenURL } from "../wailsjs/runtime/runtime";

BrowserOpenURL("https://wails.io/docs");
BrowserOpenURL("https://github.com/wailsapp/wails");
```
