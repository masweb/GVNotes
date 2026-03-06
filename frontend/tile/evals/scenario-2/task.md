# Form Validation with vee-validate Custom Rules

Register three globally available validation rules using vee-validate. All error messages must be retrieved from the Vue I18n global translation function.

## Capabilities

### Required rule

- The required rule returns an error message string for an empty string `''` [@test](../test/validation_required_empty.test.ts)
- The required rule returns an error message string for a whitespace-only string `'   '` [@test](../test/validation_required_whitespace.test.ts)
- The required rule returns `true` for a non-empty string `'hello'` [@test](../test/validation_required_valid.test.ts)

### Min-length rule

- The min rule returns an error message when the value `'ab'` is tested against minimum length `3` [@test](../test/validation_min_short.test.ts)
- The min rule returns `true` when the value `'abc'` is tested against minimum length `3` [@test](../test/validation_min_valid.test.ts)

### Confirmed rule

- The confirmed rule returns an error message when `'abc'` is compared to target `'xyz'` [@test](../test/validation_confirmed_mismatch.test.ts)
- The confirmed rule returns `true` when both values are `'secret'` [@test](../test/validation_confirmed_match.test.ts)

## Implementation

[@generates](./src/composables/useValidation.ts)

## API

```typescript { #api }
/**
 * Call this function once at application startup to register the three
 * global vee-validate rules: 'required', 'min', and 'confirmed'.
 * After registration they are available by name in all Field/Form instances.
 */
export declare function setupValidationRules(): void
```

## Dependencies { .dependencies }

### frontend 0.0.0 { .dependency }

Vue 3 + TypeScript frontend application for GVNotes. Provides the vee-validate rule-registration pattern using defineRule with i18n-sourced error messages for required, min, and confirmed rules.
