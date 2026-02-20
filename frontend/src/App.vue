<script lang="ts" setup>
import AuthView from '@/views/AuthView.vue'
import GVNotes from '@/views/GVNotes.vue'

const auth = useAuthStore()
onMounted(() => auth.init())
</script>

<template>
  <GVNotes v-if="auth.isAuthenticated" />
  <AuthView v-else-if="!auth.loading" />
</template>
