# Event Broadcaster

A module that emits named events with optional data payloads to the Go backend event system.

## Capabilities

### Event Emission with Variadic Data

Emits a named event to the Wails event system so that any Go-side (or frontend-side) listeners on that event name receive it. Zero or more additional data arguments may be supplied after the event name; they are passed through to all registered listeners.

- Calling `broadcast("userAction")` emits the `"userAction"` event with no additional data [@test](./test.test.ts)
- Calling `broadcast("dataReady", { id: 1 }, "extra")` emits `"dataReady"` with the object `{ id: 1 }` and the string `"extra"` as data arguments [@test](./test2.test.ts)
- Calling `broadcast("tick", 42)` emits `"tick"` with the number `42` as the data argument [@test](./test3.test.ts)

## Implementation

[@generates](./src/index.ts)

## API

```typescript { #api }
export function broadcast(eventName: string, ...data: any[]): void;
```

## Dependencies { .dependencies }

### @wailsapp/runtime 2.0.0 { .dependency }

Wails JavaScript runtime library providing the bridge function for emitting events from the frontend to the Go backend and other listeners.
