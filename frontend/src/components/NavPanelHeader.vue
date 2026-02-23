<script lang="ts" setup>
const { t } = useI18n()

import { IconPlus, IconNote, IconNotebook } from '@tabler/icons-vue'
import type { NavLevel } from '@/stores/navigation'

const props = defineProps<{
  level: NavLevel
}>()

const emit = defineEmits<{
  create: [kind: 'notebook' | 'note']
  rename: [newTitle: string]
}>()

// Dropdown crear
const open = ref(false)
const btnEl = useTemplateRef<HTMLElement>('btnEl')
const menuStyle = ref<Record<string, string>>({})

const toggle = () => {
  if (!open.value && btnEl.value) {
    const rect = btnEl.value.getBoundingClientRect()
    menuStyle.value = {
      position: 'fixed',
      top: `${rect.bottom + 4}px`,
      right: `${window.innerWidth - rect.right}px`,
    }
  }
  open.value = !open.value
}

const select = (kind: 'notebook' | 'note') => {
  open.value = false
  emit('create', kind)
}

const onDocClick = (e: MouseEvent) => {
  if (btnEl.value && !btnEl.value.closest('.dropdown-wrapper')?.contains(e.target as Node)) {
    open.value = false
  }
}

onMounted(() => document.addEventListener('click', onDocClick))
onUnmounted(() => document.removeEventListener('click', onDocClick))

// Renombrar título
const editingTitle = ref(false)
const titleInput = ref<HTMLInputElement | null>(null)
const { handleSubmit, resetForm } = useForm({ validateOnMount: false })
const { value: titleValue, errorMessage: titleError } = useField<string>('title', 'required', { validateOnValueUpdate: false })

const startEdit = () => {
  if (!props.level.parentId) return
  titleValue.value = props.level.title
  editingTitle.value = true
  nextTick(() => titleInput.value?.select())
}

const cancel = () => {
  editingTitle.value = false
  resetForm()
}

const submit = handleSubmit(async (values) => {
  if (values.title === props.level.title) { cancel(); return }
  emit('rename', values.title)
  cancel()
})
</script>

<template>
  <div class="nav-panel__title px-3 py-2 border-bottom flex-shrink-0 d-flex align-items-center justify-content-between gap-2">
    <form v-if="editingTitle" class="flex-grow-1 me-1" @submit.prevent="submit" @keydown.esc="cancel">
      <input
        ref="titleInput"
        v-model="titleValue"
        type="text"
        class="form-control form-control-sm fw-semibold border-0 border-bottom rounded-0 px-0"
        :class="{ 'is-invalid': titleError }"
        @blur="submit"
      />
    </form>
    <span
      v-else
      class="fw-semibold text-truncate flex-grow-1"
      :class="{ 'cursor-pointer': !!level.parentId }"
      :title="level.parentId ? t('nav.rename_hint') : undefined"
      @click="startEdit"
    >{{ level.title }}</span>

    <div class="dropdown-wrapper flex-shrink-0">
      <button
        ref="btnEl"
        class="btn btn-sm p-1 d-flex align-items-center"
        type="button"
        @click.stop="toggle"
      >
        <IconPlus :size="16" />
      </button>
      <ul class="dropdown-menu" :class="{ show: open }" :style="menuStyle">
        <li><button class="dropdown-item d-flex align-items-center gap-2" type="button" @click="select('note')"><IconNote :size="22" stroke-width="1" />{{ t('nav.new_note') }}</button></li>
        <li><button class="dropdown-item d-flex align-items-center gap-2" type="button" @click="select('notebook')"><IconNotebook :size="22" stroke-width="1" />{{ t('nav.new_notebook') }}</button></li>
      </ul>
    </div>
  </div>
</template>
