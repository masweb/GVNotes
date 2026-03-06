# Events

Bidirectional event bus for communication between the JavaScript frontend and the Go backend. JavaScript can both emit events to Go and listen for events emitted by Go.

## Import

```typescript
import {
  EventsOn,
  EventsOnce,
  EventsOnMultiple,
  EventsOff,
  EventsOffAll,
  EventsEmit,
} from "../wailsjs/runtime/runtime";
```

## Capabilities

### Subscribe to Events

```typescript { .api }
/**
 * Subscribes to an event by name, invoking callback on every occurrence.
 * @param eventName - The event name to listen for
 * @param callback - Function called with any data passed by the emitter
 * @returns An unlisten function — call it to remove this subscription
 */
function EventsOn(eventName: string, callback: (...data: any) => void): () => void;

/**
 * Subscribes to an event and triggers callback exactly once, then auto-unsubscribes.
 * @param eventName - The event name to listen for
 * @param callback - Function called with any data passed by the emitter
 * @returns An unlisten function — call it to cancel before the event fires
 */
function EventsOnce(eventName: string, callback: (...data: any) => void): () => void;

/**
 * Subscribes to an event and triggers callback at most `maxCallbacks` times, then auto-unsubscribes.
 * @param eventName - The event name to listen for
 * @param callback - Function called with any data passed by the emitter
 * @param maxCallbacks - Maximum number of times the callback will be invoked
 * @returns An unlisten function — call it to cancel early
 */
function EventsOnMultiple(
  eventName: string,
  callback: (...data: any) => void,
  maxCallbacks: number
): () => void;
```

### Unsubscribe from Events

```typescript { .api }
/**
 * Removes all listeners for the specified event name(s).
 * @param eventName - The primary event name to unsubscribe
 * @param additionalEventNames - Optional additional event names to unsubscribe in one call
 */
function EventsOff(eventName: string, ...additionalEventNames: string[]): void;

/**
 * Removes all event listeners registered via EventsOn, EventsOnce, or EventsOnMultiple.
 */
function EventsOffAll(): void;
```

### Emit Events

```typescript { .api }
/**
 * Emits an event with optional data payloads. Triggers any matching Go and JavaScript listeners.
 * @param eventName - The event name to emit
 * @param data - Optional data arguments passed to all listeners
 */
function EventsEmit(eventName: string, ...data: any): void;
```

## Usage Examples

```typescript
import { EventsOn, EventsOnce, EventsOnMultiple, EventsOff, EventsOffAll, EventsEmit } from "../wailsjs/runtime/runtime";

// Persistent listener — fires on every occurrence
const unlisten = EventsOn("file-saved", (path: string) => {
  console.log("File saved at:", path);
});
// Manually remove the listener later
unlisten();

// Listen exactly once
EventsOnce("app-ready", () => {
  initializeUI();
});

// Listen up to 5 times
const unlistenProgress = EventsOnMultiple("progress", (percent: number) => {
  updateProgressBar(percent);
}, 5);

// Emit an event with data payload to Go
EventsEmit("user-action", { type: "click", target: "submit-btn" });

// Emit with multiple data arguments
EventsEmit("coords-changed", 100, 200);

// Remove multiple events at once
EventsOff("event-a", "event-b", "event-c");

// Remove all listeners (e.g., on component unmount)
EventsOffAll();
```
