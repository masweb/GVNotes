<script lang="ts" setup>
import { IconFolder, IconFileText } from '@tabler/icons-vue'
import type { NavItem, NavLevel } from '@/stores/navigation'

const props = defineProps<{
  level: NavLevel
  activeNoteId?: string | null
  activeNotebookId?: string | null
}>()

const emit = defineEmits<{
  notebookClick: [item: NavItem & { kind: 'notebook' }]
  noteClick: [item: NavItem & { kind: 'note' }]
}>()

const searchQuery = ref('')

watch(() => props.level, () => { searchQuery.value = '' })

const notebooks = computed(() =>
  props.level.items
    .filter(i => i.kind === 'notebook')
    .filter(i => i.data.title.toLowerCase().includes(searchQuery.value.toLowerCase())) as (NavItem & { kind: 'notebook' })[]
)

const notes = computed(() =>
  props.level.items
    .filter(i => i.kind === 'note')
    .filter(i => i.data.title.toLowerCase().includes(searchQuery.value.toLowerCase())) as (NavItem & { kind: 'note' })[]
)

const noResults = computed(() => searchQuery.value !== '' && notebooks.value.length === 0 && notes.value.length === 0)
</script>

<template>
  <div class="nav-panel h-100 d-flex flex-column overflow-hidden">
    <div class="nav-panel__title px-3 py-2 fw-semibold text-truncate border-bottom flex-shrink-0">
      {{ level.title }}
    </div>
    <div class="px-2 py-2 border-bottom flex-shrink-0">
      <input
        v-model="searchQuery"
        type="search"
        class="form-control form-control-sm"
        placeholder="Buscar..."
      />
    </div>
    <ul class="list-group list-group-flush overflow-y-auto flex-grow-1">
      <li
        v-for="item in notebooks"
        :key="item.data.id"
        class="list-group-item list-group-item-action d-flex align-items-center gap-2 py-2 px-3"
        :class="{ active: item.data.id === activeNotebookId }"
        role="button"
        @click="emit('notebookClick', item)"
      >
        <IconFolder :size="16" class="flex-shrink-0 text-secondary" />
        <span class="text-truncate">{{ item.data.title }}</span>
      </li>
      <li
        v-for="item in notes"
        :key="item.data.id"
        class="list-group-item list-group-item-action d-flex align-items-center gap-2 py-2 px-3"
        :class="{ active: item.data.id === activeNoteId }"
        role="button"
        @click="emit('noteClick', item)"
      >
        <IconFileText :size="16" class="flex-shrink-0" />
        <span class="text-truncate">{{ item.data.title }}</span>
      </li>
      <li v-if="level.items.length === 0" class="list-group-item text-secondary small px-3 py-2">
        Vacío
      </li>
      <li v-else-if="noResults" class="list-group-item text-secondary small px-3 py-2">
        Sin resultados
      </li>
    </ul>
  </div>
</template>
