# Logging

Send log messages from the JavaScript frontend to the Go backend at various severity levels. Log output appears in the Go application's logging output.

## Import

```typescript
import {
  LogPrint,
  LogTrace,
  LogDebug,
  LogInfo,
  LogWarning,
  LogError,
  LogFatal,
} from "../wailsjs/runtime/runtime";
```

## Capabilities

### Log Functions

```typescript { .api }
/**
 * Logs the message as a raw message with no level prefix.
 * @param message - The message to log
 */
function LogPrint(message: string): void;

/**
 * Logs the message at the TRACE level.
 * @param message - The message to log
 */
function LogTrace(message: string): void;

/**
 * Logs the message at the DEBUG level.
 * @param message - The message to log
 */
function LogDebug(message: string): void;

/**
 * Logs the message at the INFO level.
 * @param message - The message to log
 */
function LogInfo(message: string): void;

/**
 * Logs the message at the WARNING level.
 * @param message - The message to log
 */
function LogWarning(message: string): void;

/**
 * Logs the message at the ERROR level.
 * @param message - The message to log
 */
function LogError(message: string): void;

/**
 * Logs the message at the FATAL level and quits the application.
 * @param message - The message to log before quitting
 */
function LogFatal(message: string): void;
```

## Usage Examples

```typescript
import { LogPrint, LogTrace, LogDebug, LogInfo, LogWarning, LogError, LogFatal } from "../wailsjs/runtime/runtime";

LogTrace("Entering render function");
LogDebug("Component mounted, props: " + JSON.stringify(props));
LogInfo("User authenticated successfully");
LogWarning("Deprecated API called: use newMethod() instead");
LogError("Failed to load resource: " + err.message);
LogPrint("Raw diagnostic output");

// LogFatal logs the message and then immediately quits the application
LogFatal("Critical state corruption — cannot continue");
```
