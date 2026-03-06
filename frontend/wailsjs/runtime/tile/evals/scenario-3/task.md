# Single-Fire Event Handler

A module for registering event listeners that fire at most once and are then automatically removed.

## Capabilities

### One-Time Event Subscription

Registers a callback for a named event such that the callback fires the first time the event is emitted and is then automatically deregistered. Subsequent emissions of the same event do not trigger the callback. The registration returns a cancel function that can be called before the event fires to remove the listener early.

- Calling `once("init", handler)` registers `handler` so that it is called on the first `"init"` event and then automatically removed [@test](./test.test.ts)
- After `handler` fires once on `"init"`, subsequent `"init"` events no longer invoke `handler` [@test](./test2.test.ts)
- Calling the cancel function returned by `once("init", handler)` before the event fires removes the listener without it ever being invoked [@test](./test3.test.ts)

## Implementation

[@generates](./src/index.ts)

## API

```typescript { #api }
export function once(eventName: string, handler: (...data: any[]) => void): () => void;
```

## Dependencies { .dependencies }

### @wailsapp/runtime 2.0.0 { .dependency }

Wails JavaScript runtime library providing the bridge function for registering limited-invocation event listeners from the frontend.
