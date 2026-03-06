# Authentication Pinia Store

Implement a Pinia store that manages authentication state for a desktop application with optional password protection.

## Capabilities

### State and derived state

- The initial `isAuthenticated` value is `false` [@test](../test/auth_initial_state.test.ts)
- The initial `passwordSet` value is `false` [@test](../test/auth_initial_password_set.test.ts)
- `isAuthenticated` becomes `true` after a successful authentication action [@test](../test/auth_is_authenticated.test.ts)

### Actions

- `init()` sets `passwordSet` to `true` when the backend reports a password is configured, and sets an error string on backend failure [@test](../test/auth_init.test.ts)
- `setPassword(password)` returns `true`, sets `passwordSet` to `true`, and sets `isAuthenticated` to `true` on a successful backend response [@test](../test/auth_set_password_success.test.ts)
- `setPassword(password)` returns `false` and stores an error message when the backend response reports failure [@test](../test/auth_set_password_fail.test.ts)
- `verifyPassword(password)` returns `true` and sets `isAuthenticated` to `true` on a successful backend response [@test](../test/auth_verify_success.test.ts)
- `verifyPassword(password)` returns `false` and stores an error message when credentials are wrong [@test](../test/auth_verify_fail.test.ts)
- `logout()` resets `isAuthenticated` to `false` and clears any error message [@test](../test/auth_logout.test.ts)

## Implementation

[@generates](./src/stores/auth.ts)

## API

```typescript { #api }
import { ComputedRef, Ref } from 'vue'

export declare const useAuthStore: () => {
  status: Ref<'unauthenticated' | 'authenticated'>
  passwordSet: Ref<boolean>
  loading: Ref<boolean>
  error: Ref<string | null>
  isAuthenticated: ComputedRef<boolean>
  init: () => Promise<void>
  setPassword: (password: string) => Promise<boolean>
  verifyPassword: (password: string) => Promise<boolean>
  logout: () => void
}
```

## Dependencies { .dependencies }

### frontend 0.0.0 { .dependency }

Vue 3 + TypeScript frontend application for GVNotes. Provides the Pinia authentication store pattern with async backend integration, loading/error state management, and an isAuthenticated computed property.
