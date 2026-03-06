# Event Subscription Manager

A module for managing persistent subscriptions to named events emitted by the Go backend.

## Capabilities

### Persistent Event Subscription

Registers a handler function for a named event so that every time that event fires from the Go backend, the handler is called with any accompanying data. The subscription persists until explicitly removed. The subscribe operation returns a cancel function that removes the listener when called.

- Calling `subscribe("ready", handler)` registers `handler` on the `"ready"` event and returns a callable that cancels the subscription [@test](./test.test.ts)
- Invoking the cancel function returned by `subscribe("ready", handler)` removes the listener for `"ready"` [@test](./test2.test.ts)

### Multi-Name Event Unsubscription

Removes listeners for one or more named events in a single call.

- Calling `unsubscribe("update", "refresh")` removes all registered listeners for both `"update"` and `"refresh"` in one operation [@test](./test3.test.ts)

## Implementation

[@generates](./src/index.ts)

## API

```typescript { #api }
export function subscribe(eventName: string, handler: (...data: any[]) => void): () => void;
export function unsubscribe(eventName: string, ...additionalNames: string[]): void;
```

## Dependencies { .dependencies }

### @wailsapp/runtime 2.0.0 { .dependency }

Wails JavaScript runtime library providing bridge functions for registering and removing Go-backend event listeners from the frontend.
