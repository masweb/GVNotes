# Theme Management Composable

Implement a Vue 3 composable that manages an application's light/dark theme.

## Capabilities

### Theme initialization

- When no preference is stored in browser storage and the OS reports dark mode, the initial theme is `'dark'` [@test](../test/theme_init_dark.test.ts)
- When no preference is stored and the OS reports light mode, the initial theme is `'light'` [@test](../test/theme_init_light.test.ts)

### Theme application

- Calling the set-theme function with `'dark'` updates the reactive theme value to `'dark'` and persists it to browser storage [@test](../test/theme_set_dark.test.ts)
- Calling the toggle-theme function from `'light'` switches the reactive value to `'dark'` [@test](../test/theme_toggle.test.ts)

### System preference reactivity

- When no explicit preference is stored, a change event on the OS dark-mode media query updates the theme automatically [@test](../test/theme_system_change.test.ts)
- When an explicit preference is stored, OS theme change events do not override it [@test](../test/theme_system_change_blocked.test.ts)

## Implementation

[@generates](./src/composables/useTheme.ts)

## API

```typescript { #api }
import { Ref } from 'vue'

export declare function useTheme(): {
  /** Reactive reference to the current theme name ('light' | 'dark'). */
  currentTheme: Ref<string>
  /** Set an explicit theme and persist it to storage. */
  setTheme: (theme: string) => void
  /** Toggle between 'light' and 'dark'. */
  toggleTheme: () => void
}
```

## Dependencies { .dependencies }

### frontend 0.0.0 { .dependency }

Vue 3 + TypeScript frontend application for GVNotes. Provides the composable pattern for theme management with localStorage persistence and CoreUI theme attribute integration.
