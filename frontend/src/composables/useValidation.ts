import { defineRule } from 'vee-validate'

defineRule('required', (value: string) => {
  if (!value || !value.trim()) return 'Este campo es obligatorio'
  return true
})

defineRule('min', (value: string, [min]: [number]) => {
  if (value.length < min) return `Mínimo ${min} caracteres`
  return true
})

defineRule('confirmed', (value: string, [target]: [string]) => {
  if (value !== target) return 'Las contraseñas no coinciden'
  return true
})
