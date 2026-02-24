import { CreateNotebook, CreateNote, DeleteNotebook, DeleteNote } from '../../wailsjs/go/main/App'
import { dto } from '../../wailsjs/go/models'
import type { NavItem, NavLevel } from '@/stores/navigation'

const { NotebookListItem, NoteListItem } = dto

export const useNavActions = () => {
  const nav = useNavigationStore()

  const createNotebook = async (level: NavLevel, title: string) => {
    const nb = await CreateNotebook({ parentId: level.parentId ?? undefined, title })
    const nbItem = new NotebookListItem({ id: nb.id, title: nb.title, position: nb.position, createdAt: nb.createdAt, updatedAt: nb.updatedAt })
    const target = nav.stack.find(l => l.parentId === level.parentId)
    if (target) {
      target.items = [
        ...target.items,
        { kind: 'notebook', data: nbItem },
      ]
    }
    await nav.openNotebook(nbItem)
  }

  const createNote = async (level: NavLevel, title: string) => {
    const note = await CreateNote({ notebookId: level.parentId ?? undefined, title })
    const target = nav.stack.find(l => l.parentId === level.parentId)
    if (target) {
      target.items = [
        ...target.items,
        { kind: 'note', data: new NoteListItem({ id: note.id, title: note.title, position: note.position, createdAt: note.createdAt, updatedAt: note.updatedAt }) },
      ]
    }
    nav.selectNote(note.id)
  }

  const deleteItem = async (level: NavLevel, item: NavItem) => {
    if (item.kind === 'notebook') {
      await DeleteNotebook(item.data.id)
    } else {
      await DeleteNote(item.data.id)
    }
    nav.removeItem(level.parentId, item.data.id, item.kind)
  }

  return { createNotebook, createNote, deleteItem }
}
