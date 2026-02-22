<script lang="ts" setup>
import { IconPlus, IconNote, IconNotebook } from '@tabler/icons-vue'
import type { NavLevel } from '@/stores/navigation'

defineProps<{
  level: NavLevel
}>()

const emit = defineEmits<{
  create: [kind: 'notebook' | 'note']
}>()

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
</script>

<template>
  <div class="nav-panel__title px-3 py-2 border-bottom flex-shrink-0 d-flex align-items-center justify-content-between">
    <span class="fw-semibold text-truncate">{{ level.title }}</span>
    <div class="dropdown-wrapper">
      <button
        ref="btnEl"
        class="btn btn-sm p-1 d-flex align-items-center"
        type="button"
        @click.stop="toggle"
      >
        <IconPlus :size="16" />
      </button>
      <ul class="dropdown-menu" :class="{ show: open }" :style="menuStyle">
        <li><button class="dropdown-item d-flex align-items-center gap-2" type="button" @click="select('note')"><IconNote :size="22" stroke-width="1" />Nueva nota</button></li>
        <li><button class="dropdown-item d-flex align-items-center gap-2" type="button" @click="select('notebook')"><IconNotebook :size="22" stroke-width="1" />Nueva libreta</button></li>
      </ul>
    </div>
  </div>
</template>
