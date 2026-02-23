<script lang="ts" setup>
const { t } = useI18n()

const props = defineProps<{
  visible: boolean
  kind: 'notebook' | 'note'
  title: string
}>()

const emit = defineEmits<{
  'update:visible': [value: boolean]
  confirm: []
}>()

const loading = ref(false)

const onConfirm = async () => {
  loading.value = true
  emit('confirm')
}

const close = () => {
  emit('update:visible', false)
}

watch(() => props.visible, v => {
  if (v) loading.value = false
})
</script>

<template>
  <CModal :visible="visible" alignment="top" @close="close">
    <CModalHeader>
      <CModalTitle>{{ kind === 'notebook' ? t('notebook.delete_title') : t('note.delete_title') }}</CModalTitle>
    </CModalHeader>
    <CModalBody>
      <p class="mb-1">
        {{ t('common.confirm_delete') }} <strong>{{ title }}</strong>?
      </p>
      <p v-if="kind === 'notebook'" class="text-warning small mb-0">
        {{ t('notebook.delete_warning') }}
      </p>
    </CModalBody>
    <CModalFooter>
      <button type="button" class="btn btn-secondary btn-sm" @click="close">{{ t('common.cancel') }}</button>
      <button type="button" class="btn btn-danger btn-sm" :disabled="loading" @click="onConfirm">
        {{ loading ? t('common.deleting') : t('common.delete') }}
      </button>
    </CModalFooter>
  </CModal>
</template>
