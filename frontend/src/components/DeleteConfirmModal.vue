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
      <h5 class="modal-title">{{ kind === 'notebook' ? t('notebook.delete_title') : t('note.delete_title') }}</h5>
    </CModalHeader>
    <div class="modal-body">
      <p class="mb-1">
        {{ t('common.confirm_delete') }} <strong>{{ title }}</strong>?
      </p>
      <p v-if="kind === 'notebook'" class="text-warning small mb-0">
        {{ t('notebook.delete_warning') }}
      </p>
    </div>
    <div class="modal-footer">
      <button type="button" class="btn btn-secondary btn-sm" @click="close">{{ t('common.cancel') }}</button>
      <button type="button" class="btn btn-danger btn-sm" :disabled="loading" @click="onConfirm">
        {{ loading ? t('common.deleting') : t('common.delete') }}
      </button>
    </div>
  </CModal>
</template>
