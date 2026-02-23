<script lang="ts" setup>
const { t } = useI18n()

import { useForm, useField } from 'vee-validate'
import type { NavLevel } from '@/stores/navigation'

const props = defineProps<{
  visible: boolean
  kind: 'notebook' | 'note'
  level: NavLevel | null
}>()

const emit = defineEmits<{
  'update:visible': [value: boolean]
}>()

const { createNotebook, createNote } = useNavActions()

const title = computed(() =>
  props.kind === 'notebook' ? t('notebook.new_title') : t('note.new_title')
)

const { handleSubmit, resetForm } = useForm({
  validateOnMount: false,
})

const { value: name, errorMessage } = useField<string>('name', 'required', {
  validateOnValueUpdate: false,
})

const loading = ref(false)
const error = ref<string | null>(null)

const onSubmit = handleSubmit(async () => {
  if (!props.level) return
  loading.value = true
  error.value = null
  try {
    if (props.kind === 'notebook') {
      await createNotebook(props.level, name.value)
    } else {
      await createNote(props.level, name.value)
    }
    close()
  } catch (e: any) {
    error.value = e?.message ?? t('common.error_create')
  } finally {
    loading.value = false
  }
})

const inputEl = useTemplateRef<HTMLInputElement>('inputEl')

const close = () => {
  emit('update:visible', false)
}

watch(() => props.visible, async (v) => {
  if (v) {
    resetForm()
    error.value = null
    await nextTick()
    inputEl.value?.focus()
  }
})
</script>

<template>
  <CModal :visible="visible" alignment="top" @close="close">
    <CModalHeader>
      <CModalTitle>{{ title }}</CModalTitle>
    </CModalHeader>
    <form @submit.prevent="onSubmit">
      <CModalBody>
        <div class="mb-1">
          <input
            ref="inputEl"
            v-model="name"
            type="text"
            class="form-control"
            :class="{ 'is-invalid': errorMessage }"
            :placeholder="t('common.name_placeholder')"
            @keydown.enter.prevent="onSubmit"
          />
          <div v-if="errorMessage" class="invalid-feedback">{{ errorMessage }}</div>
        </div>
        <p v-if="error" class="text-danger small mt-2 mb-0">{{ error }}</p>
      </CModalBody>
      <CModalFooter>
        <button type="button" class="btn btn-secondary btn-sm" @click="close">{{ t('common.cancel') }}</button>
        <button type="submit" class="btn btn-primary btn-sm" :disabled="loading">
          {{ loading ? t('common.creating') : t('common.create') }}
        </button>
      </CModalFooter>
    </form>
  </CModal>
</template>
