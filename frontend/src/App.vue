<script lang="ts" setup>
import AuthView from '@/views/AuthView.vue'
import GVNotes from '@/views/GVNotes.vue'

const auth = useAuthStore()
onMounted(() => auth.init())

// Splitpanes lanza un error interno al desmontarse con panes activos (bug conocido).
// Lo suprimimos aquí para evitar el warning en consola.
onErrorCaptured((err) => {
  if (err instanceof TypeError && err.message.includes('index')) return false
})
</script>

<template>
  <GVNotes v-if="auth.isAuthenticated" />
  <AuthView v-else-if="!auth.loading" />
</template>
