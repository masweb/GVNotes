import { ListNotebooks, ListNotes } from '../../wailsjs/go/main/App'
import type { dto } from '../../wailsjs/go/models'

export type NavItem =
  | { kind: 'notebook'; data: dto.NotebookListItem }
  | { kind: 'note'; data: dto.NoteListItem }

export interface NavLevel {
  parentId: string | null // null = raíz
  title: string
  items: NavItem[]
}

export const useNavigationStore = defineStore('navigation', () => {
  const stack = ref<NavLevel[]>([])
  const selectedNoteId = ref<string | null>(null)
  const loading = ref(false)

  const leftPanel = computed(() => (stack.value.length > 1 ? stack.value[stack.value.length - 2] : stack.value[0] ?? null))
  const rightPanel = computed(() => (stack.value.length > 1 ? stack.value[stack.value.length - 1] : null))
  const canGoBack = computed(() => stack.value.length > 1)
  const activeNotebookId = computed(() => (stack.value.length > 1 ? stack.value[stack.value.length - 1].parentId : null))

  const loadLevel = async (parentId: string | null, title: string): Promise<NavLevel> => {
    const id = parentId ?? ''
    const [notebooks, notes] = await Promise.all([ListNotebooks(id), ListNotes(id)])
    const items: NavItem[] = [
      ...notebooks.map(n => ({ kind: 'notebook' as const, data: n })),
      ...notes.map(n => ({ kind: 'note' as const, data: n })),
    ]
    return { parentId, title, items }
  }

  const init = async () => {
    loading.value = true
    try {
      const root = await loadLevel(null, 'Inicio')
      stack.value = [root]
    } finally {
      loading.value = false
    }
  }

  // Panel izquierdo: reemplaza el nivel derecho (cambia qué notebook está abierto)
  const openNotebookFromLeft = async (notebook: dto.NotebookListItem) => {
    loading.value = true
    try {
      const level = await loadLevel(notebook.id, notebook.title)
      stack.value = [...stack.value.slice(0, -1), level]
      selectedNoteId.value = null
    } finally {
      loading.value = false
    }
  }

  // Panel derecho o panel único: navega más profundo
  const openNotebook = async (notebook: dto.NotebookListItem) => {
    loading.value = true
    try {
      const level = await loadLevel(notebook.id, notebook.title)
      stack.value = [...stack.value, level]
      selectedNoteId.value = null
    } finally {
      loading.value = false
    }
  }

  const selectNote = (noteId: string) => {
    selectedNoteId.value = noteId
  }

  const goBack = () => {
    if (stack.value.length > 1) {
      stack.value = stack.value.slice(0, -1)
      selectedNoteId.value = null
    }
  }

  return { stack, selectedNoteId, loading, leftPanel, rightPanel, canGoBack, activeNotebookId, init, openNotebook, openNotebookFromLeft, selectNote, goBack }
})
