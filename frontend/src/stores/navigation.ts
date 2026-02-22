import { ListNotebooks, ListNotes, UpdateNotebookPosition, UpdateNotePosition, UpdateNotebookTitle } from '../../wailsjs/go/main/App'
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
  const activeNotebookId = computed(() => (stack.value.length > 1 ? stack.value[stack.value.length - 1]?.parentId ?? null : null))

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

  const reorderItems = async (level: NavLevel, kind: 'notebook' | 'note', newItems: NavItem[]) => {
    // Buscar el nivel en el stack por parentId para mutar el objeto reactivo directamente
    const target = stack.value.find(l => l.parentId === level.parentId)
    if (target) {
      const otherKind = target.items.filter(i => i.kind !== kind)
      target.items = kind === 'notebook'
        ? [...newItems, ...otherKind]
        : [...otherKind, ...newItems]
    }
    // Persistir en backend
    await Promise.all(
      newItems.map((item, idx) =>
        kind === 'notebook'
          ? UpdateNotebookPosition(item.data.id, { position: idx * 1000 })
          : UpdateNotePosition(item.data.id, { position: idx * 1000 }),
      ),
    )
    newItems.forEach((item, idx) => { item.data.position = idx * 1000 })
  }

  const reset = () => {
    stack.value = []
    selectedNoteId.value = null
  }

  const renameNotebook = async (notebookId: string, newTitle: string) => {
    await UpdateNotebookTitle(notebookId, { title: newTitle })
    // Actualizar el título en el nivel cuyo parentId es este notebook
    const level = stack.value.find(l => l.parentId === notebookId)
    if (level) level.title = newTitle
    // Actualizar también el item en el nivel padre
    renameItem(notebookId, 'notebook', newTitle)
  }

  const renameItem = (itemId: string, kind: 'notebook' | 'note', newTitle: string) => {
    for (const level of stack.value) {
      const item = level.items.find(i => i.kind === kind && i.data.id === itemId)
      if (item) {
        item.data.title = newTitle
        break
      }
    }
  }

  const removeItem = (levelParentId: string | null, itemId: string, itemKind: 'notebook' | 'note') => {
    const target = stack.value.find(l => l.parentId === levelParentId)
    if (target) {
      target.items = target.items.filter(i => i.data.id !== itemId)
    }
    if (itemKind === 'notebook') {
      const stackIdx = stack.value.findIndex(l => l.parentId === itemId)
      if (stackIdx !== -1) {
        stack.value = stack.value.slice(0, stackIdx)
      }
    }
  }

  return { stack, selectedNoteId, loading, leftPanel, rightPanel, canGoBack, activeNotebookId, init, openNotebook, openNotebookFromLeft, selectNote, goBack, reorderItems, removeItem, renameItem, renameNotebook, reset }
})
