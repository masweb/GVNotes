import { GetAuthStatus, SetPassword, VerifyPassword } from '../../wailsjs/go/main/App'
import { i18n } from '@/i18n/i18n'

const t = i18n.global.t

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
      error.value = t('auth.error_connect')
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
        error.value = res.error ?? t('auth.error_set_password')
        return false
      }
      passwordSet.value = true
      status.value = 'authenticated'
      return true
    } catch (e: any) {
      error.value = e?.message ?? t('auth.error_set_password')
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
        error.value = t('auth.error_wrong_password')
        return false
      }
      status.value = 'authenticated'
      return true
    } catch (e: any) {
      error.value = e?.message ?? t('auth.error_verify_password')
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
