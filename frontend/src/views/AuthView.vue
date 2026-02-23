<script lang="ts" setup>
const { t } = useI18n()
const auth = useAuthStore()

// Modo: 'set' si no hay contraseña aún, 'verify' si ya existe
const mode = computed(() => (auth.passwordSet ? 'verify' : 'set'))

const formOpts = { validateOnMount: false, validateOnModelUpdate: false, validateOnBlur: false }

// Formulario "establecer contraseña"
const { handleSubmit: handleSet, resetForm: resetSet } = useForm(formOpts)
const { value: newPassword, errorMessage: newPasswordError } = useField<string>('newPassword', 'required|min:8', {
 validateOnValueUpdate: false
})
const { value: confirmPassword, errorMessage: confirmPasswordError } = useField<string>(
 'confirmPassword',
 (val: string) => {
  if (!val) return t('validation.required')
  if (val !== newPassword.value) return t('validation.confirmed')
  return true
 },
 { validateOnValueUpdate: false }
)

// Formulario "verificar contraseña"
const { handleSubmit: handleVerify } = useForm(formOpts)
const { value: password, errorMessage: passwordError } = useField<string>('password', 'required', {
 validateOnValueUpdate: false
})

const submitSet = handleSet(async () => {
 await auth.setPassword(newPassword.value)
 if (auth.isAuthenticated) resetSet()
})

const submitVerify = handleVerify(async () => {
 await auth.verifyPassword(password.value)
})

const firstInput = ref<HTMLInputElement | null>(null)
onMounted(() => firstInput.value?.focus())
</script>

<template>
 <div class="auth-wrap d-flex align-items-center justify-content-center vh-100" style="--wails-draggable: drag">
  <div class="auth-card card p-2" style="width: 360px">
   <div class="card-body">
    <h5 class="mb-4 text-center fw-semibold">
     {{ mode === 'set' ? t('auth.title_set') : t('auth.title_verify') }}
    </h5>

    <!-- Establecer contraseña -->
    <form v-if="mode === 'set'" @submit.prevent="submitSet">
     <div class="mb-3">
      <label class="form-label">{{ t('auth.password_label') }}</label>
      <input
       ref="firstInput"
       v-model="newPassword"
       type="password"
       class="form-control"
       :class="{ 'is-invalid': newPasswordError }"
       autocomplete="new-password"
       :placeholder="t('auth.password_placeholder_new')"
      />
      <div v-if="newPasswordError" class="invalid-feedback">{{ newPasswordError }}</div>
     </div>

     <div class="mb-4">
      <label class="form-label">{{ t('auth.confirm_label') }}</label>
      <input
       v-model="confirmPassword"
       type="password"
       class="form-control"
       :class="{ 'is-invalid': confirmPasswordError }"
       autocomplete="new-password"
       :placeholder="t('auth.confirm_placeholder')"
      />
      <div v-if="confirmPasswordError" class="invalid-feedback">{{ confirmPasswordError }}</div>
     </div>

     <div v-if="auth.error" class="alert alert-danger py-2 mb-3">{{ auth.error }}</div>

     <button type="submit" class="btn btn-primary w-100" :disabled="auth.loading">
      <span v-if="auth.loading" class="spinner-border spinner-border-sm me-2" />
      {{ t('auth.btn_set') }}
     </button>
    </form>

    <!-- Verificar contraseña -->
    <form v-else @submit.prevent="submitVerify">
     <div class="mb-4">
      <label class="form-label">{{ t('auth.password_label') }}</label>
      <input
       ref="firstInput"
       v-model="password"
       type="password"
       class="form-control"
       :class="{ 'is-invalid': passwordError }"
       autocomplete="current-password"
       :placeholder="t('auth.password_placeholder_current')"
      />
      <div v-if="passwordError" class="invalid-feedback">{{ passwordError }}</div>
     </div>

     <div v-if="auth.error" class="alert alert-danger py-2 mb-3">{{ auth.error }}</div>

     <button type="submit" class="btn btn-primary w-100" :disabled="auth.loading">
      <span v-if="auth.loading" class="spinner-border spinner-border-sm me-2" />
      {{ t('auth.btn_enter') }}
     </button>
    </form>
   </div>
  </div>
 </div>
</template>
