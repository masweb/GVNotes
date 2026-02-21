<script lang="ts" setup>
import { IconFolder, IconFileText, IconGripVertical } from '@tabler/icons-vue'
import Sortable, { type SortableEvent } from 'sortablejs'
import type { NavItem, NavLevel } from '@/stores/navigation'

const props = defineProps<{
  level: NavLevel
  activeNoteId?: string | null
  activeNotebookId?: string | null
}>()

const emit = defineEmits<{
  notebookClick: [item: NavItem & { kind: 'notebook' }]
  noteClick: [item: NavItem & { kind: 'note' }]
  reorder: [level: NavLevel, kind: 'notebook' | 'note', newItems: NavItem[]]
}>()

const searchQuery = ref('')
const isSearching = computed(() => searchQuery.value !== '')

watch(() => props.level, () => { searchQuery.value = '' })

const notebooks = computed(() =>
  props.level.items.filter(i => i.kind === 'notebook') as (NavItem & { kind: 'notebook' })[]
)

const notes = computed(() =>
  props.level.items.filter(i => i.kind === 'note') as (NavItem & { kind: 'note' })[]
)

const filteredNotebooks = computed(() =>
  isSearching.value
    ? notebooks.value.filter(i => i.data.title.toLowerCase().includes(searchQuery.value.toLowerCase()))
    : notebooks.value
)

const filteredNotes = computed(() =>
  isSearching.value
    ? notes.value.filter(i => i.data.title.toLowerCase().includes(searchQuery.value.toLowerCase()))
    : notes.value
)

const noResults = computed(() => isSearching.value && filteredNotebooks.value.length === 0 && filteredNotes.value.length === 0)

const notebooksEl = useTemplateRef<HTMLElement>('notebooksEl')
const notesEl = useTemplateRef<HTMLElement>('notesEl')

const makeSortable = (el: HTMLElement, getList: () => NavItem[], kind: 'notebook' | 'note'): Sortable => {
  return Sortable.create(el, {
    handle: '.drag-handle',
    animation: 150,
    ghostClass: 'drag-ghost',
    // @ts-expect-error supportPointer no está en los tipos pero sí en la implementación
    supportPointer: false,
    onEnd: (evt: SortableEvent) => {
      const { oldIndex, newIndex } = evt
      if (oldIndex === undefined || newIndex === undefined || oldIndex === newIndex) return
      const list = [...getList()]
      const moved = list.splice(oldIndex, 1)[0]
      if (!moved) return
      list.splice(newIndex, 0, moved)
      setTimeout(() => emit('reorder', props.level, kind, list as NavItem[]), 0)
    },
  })
}

let sortableNotebooks: Sortable | null = null
let sortableNotes: Sortable | null = null

onMounted(() => {
  if (notebooksEl.value) sortableNotebooks = makeSortable(notebooksEl.value, () => notebooks.value, 'notebook')
  if (notesEl.value) sortableNotes = makeSortable(notesEl.value, () => notes.value, 'note')
})

onUnmounted(() => {
  sortableNotebooks?.destroy()
  sortableNotes?.destroy()
})

watch(isSearching, (v) => {
  sortableNotebooks?.option('disabled', v)
  sortableNotes?.option('disabled', v)
})
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
      <ul ref="notebooksEl" class="list-unstyled m-0">
        <li
          v-for="item in filteredNotebooks"
          :key="item.data.id"
          class="list-group-item d-flex align-items-center gap-2 py-2 px-3 nav-draggable-item"
          :class="{ active: item.data.id === activeNotebookId }"
          role="button"
          @click="emit('notebookClick', item)"
        >
          <IconGripVertical
            :size="14"
            class="drag-handle flex-shrink-0 text-secondary drag-handle-icon"
            :class="{ invisible: isSearching }"
          />
          <IconFolder :size="16" class="flex-shrink-0 text-secondary" />
          <span class="text-truncate">{{ item.data.title }}</span>
        </li>
      </ul>

      <ul ref="notesEl" class="list-unstyled m-0">
        <li
          v-for="item in filteredNotes"
          :key="item.data.id"
          class="list-group-item d-flex align-items-center gap-2 py-2 px-3 nav-draggable-item"
          :class="{ active: item.data.id === activeNoteId }"
          role="button"
          @click="emit('noteClick', item)"
        >
          <IconGripVertical
            :size="14"
            class="drag-handle flex-shrink-0 text-secondary drag-handle-icon"
            :class="{ invisible: isSearching }"
          />
          <IconFileText :size="16" class="flex-shrink-0" />
          <span class="text-truncate">{{ item.data.title }}</span>
        </li>
      </ul>

      <li v-if="level.items.length === 0" class="list-group-item text-secondary small px-3 py-2">
        Vacío
      </li>
      <li v-else-if="noResults" class="list-group-item text-secondary small px-3 py-2">
        Sin resultados
      </li>
    </ul>
  </div>
</template>
