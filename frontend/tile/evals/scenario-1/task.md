# Vue I18n Locale Management

Set up Vue I18n for a multi-language application and implement a locale management composable.

## Capabilities

### i18n instance creation

- The i18n instance uses Spanish (`'es'`) as the default locale when browser storage has no saved language [@test](../test/i18n_default_locale.test.ts)
- The i18n instance uses English (`'en'`) as the fallback locale [@test](../test/i18n_fallback.test.ts)
- Global injection is enabled so translation functions are available without explicit import in components [@test](../test/i18n_global_injection.test.ts)

### Locale composable

- Calling the locale-change function with `'en'` updates the active locale to `'en'` [@test](../test/locale_set_en.test.ts)
- Calling the locale-change function persists the chosen code to browser storage under the key `'lang'` [@test](../test/locale_persist.test.ts)
- The available-locales list contains an entry with code `'es'` and label `'Español'`, and an entry with code `'en'` and label `'English'` [@test](../test/locale_available.test.ts)

## Implementation

[@generates](./src/i18n/i18n.ts)

## API

```typescript { #api }
import { I18n } from 'vue-i18n'
import { Ref } from 'vue'

export interface LocaleOption {
  code: string
  label: string
}

/** Singleton i18n instance to install into the Vue app. */
export declare const i18n: I18n

export declare function useLocale(): {
  /** Reactive reference to the current locale code. */
  currentLocale: Ref<string>
  /** Static list of all supported locales. */
  availableLocales: LocaleOption[]
  /** Change the active locale and persist it to storage. */
  setLocale: (code: string) => void
}
```

## Dependencies { .dependencies }

### frontend 0.0.0 { .dependency }

Vue 3 + TypeScript frontend application for GVNotes. Provides the i18n setup pattern with createI18n, localStorage-based locale persistence, and a useLocale composable wrapping useI18n.
