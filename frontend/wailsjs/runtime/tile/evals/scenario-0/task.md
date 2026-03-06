# Severity Log Dispatcher

A module that dispatches log messages to the Go backend at the appropriate severity level.

## Capabilities

### Severity-Based Log Routing

Routes a log message to the corresponding backend logging function based on the provided severity level string. Supported levels are: `print`, `trace`, `debug`, `info`, `warning`, `error`, and `fatal`. Messages sent at the `fatal` level will cause the application to quit after the message is delivered.

- Calling `dispatch("debug", "starting up")` delegates the message to the debug-level backend logger with the string `"starting up"` [@test](./test.test.ts)
- Calling `dispatch("error", "connection lost")` delegates the message to the error-level backend logger with the string `"connection lost"` [@test](./test2.test.ts)
- Calling `dispatch("fatal", "unrecoverable state")` delegates to the fatal-level backend logger with `"unrecoverable state"` [@test](./test3.test.ts)
- Calling `dispatch("warning", "low memory")` delegates the message to the warning-level backend logger with `"low memory"` [@test](./test4.test.ts)

## Implementation

[@generates](./src/index.ts)

## API

```typescript { #api }
export type LogLevel = "print" | "trace" | "debug" | "info" | "warning" | "error" | "fatal";

export function dispatch(level: LogLevel, message: string): void;
```

## Dependencies { .dependencies }

### @wailsapp/runtime 2.0.0 { .dependency }

Wails JavaScript runtime library providing bridge functions for communicating with the Go backend, including backend logging at multiple severity levels.
