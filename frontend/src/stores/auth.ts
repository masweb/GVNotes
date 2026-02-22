import { GetAuthStatus, SetPassword, VerifyPassword } from '../../wailsjs/go/main/App'

type AuthStatus = 'unauthenticated' | 'authenticated'

export const useAuthStore = defineStore('auth', () => {
  const status = ref<AuthStatus>('unauthenticated')
  const passwordSet = ref(false)
  const loading = ref(false)
  const error = ref<string | null>(null)

  const isAuthenticated = computed(() => status.value === 'authenticated')

  const init = async () => {
    loading.value = true
    error.value = null
    try {
      const result = await GetAuthStatus()
      passwordSet.value = result.isPasswordSet
    } catch (e) {
      error.value = 'Error al conectar con la aplicación'
    } finally {
      loading.value = false
    }
  }

  const setPassword = async (password: string): Promise<boolean> => {
    loading.value = true
    error.value = null
    try {
      const res = await SetPassword(password)
      if (!res.success) {
        error.value = res.error ?? 'Error al establecer la contraseña'
        return false
      }
      passwordSet.value = true
      status.value = 'authenticated'
      return true
    } catch (e: any) {
      error.value = e?.message ?? 'Error al establecer la contraseña'
      return false
    } finally {
      loading.value = false
    }
  }

  const verifyPassword = async (password: string): Promise<boolean> => {
    loading.value = true
    error.value = null
    try {
      const res = await VerifyPassword(password)
      if (!res.success) {
        error.value = 'Contraseña incorrecta'
        return false
      }
      status.value = 'authenticated'
      return true
    } catch (e: any) {
      error.value = e?.message ?? 'Error al verificar la contraseña'
      return false
    } finally {
      loading.value = false
    }
  }

  const logout = () => {
    status.value = 'unauthenticated'
    error.value = null
  }

  return { status, passwordSet, loading, error, isAuthenticated, init, setPassword, verifyPassword, logout }
})
