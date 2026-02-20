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
      await SetPassword(password)
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
      const ok = await VerifyPassword(password)
      if (ok) {
        status.value = 'authenticated'
      } else {
        error.value = 'Contraseña incorrecta'
      }
      return ok
    } catch (e: any) {
      error.value = e?.message ?? 'Error al verificar la contraseña'
      return false
    } finally {
      loading.value = false
    }
  }

  return { status, passwordSet, loading, error, isAuthenticated, init, setPassword, verifyPassword }
})
