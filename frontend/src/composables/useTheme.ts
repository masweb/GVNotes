import { ref } from 'vue'

const STORAGE_KEY = 'coreui-docs-theme'

const getInitialTheme = (): string => {
  const stored = localStorage.getItem(STORAGE_KEY)
  if (stored) return stored
  return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light'
}

const currentTheme = ref(getInitialTheme())

const applyTheme = (theme: string) => {
  document.documentElement.setAttribute('data-coreui-theme', theme)
  localStorage.setItem(STORAGE_KEY, theme)
  currentTheme.value = theme
}

applyTheme(currentTheme.value)

window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', (e) => {
  if (!localStorage.getItem(STORAGE_KEY)) {
    applyTheme(e.matches ? 'dark' : 'light')
  }
})

export const useTheme = () => ({
  currentTheme,
  setTheme: applyTheme,
  toggleTheme: () => applyTheme(currentTheme.value === 'light' ? 'dark' : 'light'),
})
