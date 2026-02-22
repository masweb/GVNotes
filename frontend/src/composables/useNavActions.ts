import { CreateNotebook, CreateNote } from '../../wailsjs/go/main/App'
import { dto } from '../../wailsjs/go/models'
import type { NavLevel } from '@/stores/navigation'

const { NotebookListItem, NoteListItem } = dto

export const useNavActions = () => {
  const nav = useNavigationStore()

  const createNotebook = async (level: NavLevel, title: string) => {
    const nb = await CreateNotebook({ parentId: level.parentId ?? undefined, title })
    const idx = nav.stack.findIndex(l => l.parentId === level.parentId)
    if (idx !== -1) {
      const target = nav.stack[idx]
      nav.stack[idx] = {
        ...target,
        items: [
          ...target.items,
          { kind: 'notebook', data: new NotebookListItem({ id: nb.id, title: nb.title, position: nb.position, createdAt: nb.createdAt, updatedAt: nb.updatedAt }) },
        ],
      }
    }
  }

  const createNote = async (level: NavLevel, title: string) => {
    const note = await CreateNote({ notebookId: level.parentId ?? undefined, title })
    const idx = nav.stack.findIndex(l => l.parentId === level.parentId)
    if (idx !== -1) {
      const target = nav.stack[idx]
      nav.stack[idx] = {
        ...target,
        items: [
          ...target.items,
          { kind: 'note', data: new NoteListItem({ id: note.id, title: note.title, position: note.position, createdAt: note.createdAt, updatedAt: note.updatedAt }) },
        ],
      }
    }
  }

  return { createNotebook, createNote }
}
