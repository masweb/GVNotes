<script lang="ts" setup>
import { IconSun, IconMoon } from '@tabler/icons-vue'

const { t } = useI18n()
const { currentTheme, toggleTheme } = useTheme()

defineProps<{ visible: boolean }>()
const emit = defineEmits<{ 'update:visible': [value: boolean] }>()

const close = () => emit('update:visible', false)
</script>

<template>
  <CModal :visible="visible" alignment="top" @close="close">
    <CModalHeader>
      <h5 class="modal-title">{{ t('settings.title') }}</h5>
    </CModalHeader>
    <div class="modal-body">
      <div class="d-flex align-items-center justify-content-between">
        <span class="small">{{ t('settings.theme') }}</span>
        <button type="button" class="btn btn-sm d-flex align-items-center gap-1" @click="toggleTheme">
          <IconSun v-if="currentTheme === 'dark'" :size="18" stroke-width="1" />
          <IconMoon v-else :size="18" stroke-width="1" />
          <span class="small">{{ currentTheme === 'dark' ? t('settings.theme_light') : t('settings.theme_dark') }}</span>
        </button>
      </div>
    </div>
    <div class="modal-footer">
      <button type="button" class="btn btn-secondary btn-sm" @click="close">{{ t('common.close') }}</button>
    </div>
  </CModal>
</template>
