# Composables

Reusable Vue composables for theme management, locale switching, navigation actions, and form validation.

## useTheme

```typescript
import { useTheme } from '@/composables/useTheme'
// Also auto-imported in Vue SFCs.
```

### useTheme()

```typescript
useTheme(): { currentTheme: Ref<string>; setTheme: (theme: string) => void; toggleTheme: () => void }
```
{ .api }

Provides reactive light/dark theme control. State is **module-level** (singleton): all calls to `useTheme()` share the same `currentTheme` ref. The theme is persisted to `localStorage` under key `'coreui-docs-theme'` and automatically applied as `data-coreui-theme` attribute on `document.documentElement`.

On first load (no stored preference), the composable follows the OS `prefers-color-scheme` media query. If the user has set a preference, that preference overrides the OS setting.

**Returns:**

| Property | Type | Description |
|----------|------|-------------|
| `currentTheme` | `Ref<string>` | Reactive theme name: `'light'` or `'dark'`. |
| `setTheme` | `(theme: string) => void` | Sets theme, updates DOM, persists to localStorage. |
| `toggleTheme` | `() => void` | Switches between `'light'` and `'dark'`. |

**Example:**
```typescript
const { currentTheme, setTheme, toggleTheme } = useTheme()

setTheme('dark')       // force dark
toggleTheme()          // flip between light/dark
console.log(currentTheme.value) // 'light'

// Bind in template:
// <CButton @click="toggleTheme">{{ currentTheme }}</CButton>
```

---

## useLocale

```typescript
import { useLocale, availableLocales } from '@/composables/useLocale'
import type { LocaleOption } from '@/composables/useLocale'
// Also auto-imported in Vue SFCs.
```

### LocaleOption

```typescript
interface LocaleOption {
  code: string   // e.g. 'en', 'es'
  label: string  // e.g. 'English', 'Español'
}
```
{ .api }

### availableLocales

```typescript
const availableLocales: LocaleOption[]
// Value: [{ code: 'es', label: 'Español' }, { code: 'en', label: 'English' }]
```
{ .api }

Module-level constant listing all supported locales.

### useLocale()

```typescript
useLocale(): {
  currentLocale: Ref<string>
  availableLocales: LocaleOption[]
  setLocale: (code: string) => void
}
```
{ .api }

Provides locale switching for vue-i18n. Persists the selected locale to `localStorage` under key `'lang'`.

**Returns:**

| Property | Type | Description |
|----------|------|-------------|
| `currentLocale` | `Ref<string>` | The vue-i18n global locale ref. Mutating it changes the app language immediately. |
| `availableLocales` | `LocaleOption[]` | Array of supported locale options. |
| `setLocale` | `(code: string) => void` | Sets locale on the vue-i18n instance and persists to localStorage. |

**Example:**
```typescript
const { currentLocale, availableLocales, setLocale } = useLocale()

setLocale('en')   // switch to English
console.log(currentLocale.value) // 'en'

// Render locale selector:
// <select v-model="currentLocale">
//   <option v-for="loc in availableLocales" :key="loc.code" :value="loc.code">
//     {{ loc.label }}
//   </option>
// </select>
```

---

## useNavActions

```typescript
import { useNavActions } from '@/composables/useNavActions'
// Also auto-imported in Vue SFCs.
```

### useNavActions()

```typescript
useNavActions(): {
  createNotebook: (level: NavLevel, title: string) => Promise<void>
  createNote: (level: NavLevel, title: string) => Promise<void>
  deleteItem: (level: NavLevel, item: NavItem) => Promise<void>
}
```
{ .api }

Coordinates backend API calls with `useNavigationStore` state updates. Each action calls the Wails backend and then updates the navigation stack reactively.

**Returns:**

#### createNotebook

```typescript
createNotebook(level: NavLevel, title: string): Promise<void>
```
{ .api }

Creates a notebook via the backend (`CreateNotebook`), adds it to the navigation stack at the given level, and then calls `nav.openNotebook()` to drill into the new notebook.

**Parameters:**
- `level` — The `NavLevel` in which to create the notebook (determines `parentId`).
- `title` — Notebook title.

---

#### createNote

```typescript
createNote(level: NavLevel, title: string): Promise<void>
```
{ .api }

Creates a note via the backend (`CreateNote`), adds it to the navigation stack at the given level, and selects it for editing (`nav.selectNote()`).

**Parameters:**
- `level` — The `NavLevel` in which to create the note.
- `title` — Note title.

---

#### deleteItem

```typescript
deleteItem(level: NavLevel, item: NavItem): Promise<void>
```
{ .api }

Deletes either a notebook (`DeleteNotebook`) or a note (`DeleteNote`) via the backend based on `item.kind`, then calls `nav.removeItem()` to update the stack.

**Parameters:**
- `level` — The `NavLevel` containing the item.
- `item` — The `NavItem` to delete (discriminated union: `kind === 'notebook'` or `kind === 'note'`).

---

**Example:**
```typescript
const { createNotebook, createNote, deleteItem } = useNavActions()
const nav = useNavigationStore()

// Create notebook at current level
await createNotebook(nav.leftPanel!, 'My New Notebook')

// Create note in active notebook
await createNote(nav.rightPanel ?? nav.leftPanel!, 'Meeting Notes')

// Delete an item
await deleteItem(level, selectedItem)
```

---

## useValidation

```typescript
// src/composables/useValidation.ts
// This file is executed as a side-effect in main.ts — no explicit import needed.
// Rules are globally registered with vee-validate's defineRule().
```

### Global Validation Rules

Registered with vee-validate's `defineRule()`. Use these rule names as strings in `useField()` or `useForm()` validation schemas. All error messages are i18n-translated.

#### required

```typescript
defineRule('required', (value: string) => string | true)
```
{ .api }

Validates that the value is non-empty and non-whitespace.

Usage: `useField('name', 'required')`

---

#### min

```typescript
defineRule('min', (value: string, [min]: [number]) => string | true)
```
{ .api }

Validates that the string length is >= the specified minimum.

Usage: `useField('password', 'min:8')`

---

#### confirmed

```typescript
defineRule('confirmed', (value: string, [target]: [string]) => string | true)
```
{ .api }

Validates that the value equals the `target` string. Used for password confirmation.

Usage: `useField('confirmPassword', 'confirmed:@password')`

---

**Form usage with vee-validate:**
```typescript
import { useForm, useField } from 'vee-validate'

const { handleSubmit } = useForm({ validateOnMount: false })
const { value: password, errorMessage: passwordError } = useField(
  'password', 'min:8', { validateOnValueUpdate: false }
)
const { value: confirm, errorMessage: confirmError } = useField(
  'confirmPassword', 'confirmed:@password', { validateOnValueUpdate: false }
)

const onSubmit = handleSubmit(async (values) => {
  // validation passes here
  await authStore.setPassword(values.password)
})
```
