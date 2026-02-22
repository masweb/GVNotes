<script lang="ts" setup>
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
      <CModalTitle>{{ kind === 'notebook' ? 'Eliminar libreta' : 'Eliminar nota' }}</CModalTitle>
    </CModalHeader>
    <CModalBody>
      <p class="mb-1">
        ¿Eliminar <strong>{{ title }}</strong>?
      </p>
      <p v-if="kind === 'notebook'" class="text-warning small mb-0">
        Se eliminarán también todas las libretas y notas que contiene.
      </p>
    </CModalBody>
    <CModalFooter>
      <button type="button" class="btn btn-secondary btn-sm" @click="close">Cancelar</button>
      <button type="button" class="btn btn-danger btn-sm" :disabled="loading" @click="onConfirm">
        {{ loading ? 'Eliminando…' : 'Eliminar' }}
      </button>
    </CModalFooter>
  </CModal>
</template>
