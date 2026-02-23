<script lang="ts" setup>
const { t } = useI18n()

import { Splitpanes, Pane } from 'splitpanes'
import { IconArrowLeft, IconPower, IconSun, IconMoon } from '@tabler/icons-vue'
import type { NavItem, NavLevel } from '@/stores/navigation'

const nav = useNavigationStore()
const auth = useAuthStore()
const { currentTheme, toggleTheme } = useTheme()

const logout = () => {
 nav.reset()
 auth.logout()
}

onMounted(() => nav.init())

const onNotebookClickLeft = (item: NavItem & { kind: 'notebook' }) => {
 nav.openNotebookFromLeft(item.data)
}

const onNotebookClickRight = (item: NavItem & { kind: 'notebook' }) => {
 nav.openNotebook(item.data)
}

const onNoteClick = (item: NavItem & { kind: 'note' }) => {
 nav.selectNote(item.data.id)
}

const onNoteClickFromLeft = (item: NavItem & { kind: 'note' }) => {
 nav.selectNoteFromRoot(item.data.id)
}

const onReorder = (level: NavLevel, kind: 'notebook' | 'note', newItems: NavItem[]) => {
 nav.reorderItems(level, kind, newItems)
}

// Modal de alta
const modalVisible = ref(false)
const modalKind = ref<'notebook' | 'note'>('notebook')
const modalLevel = ref<NavLevel | null>(null)

const onCreate = (level: NavLevel, kind: 'notebook' | 'note') => {
 modalLevel.value = level
 modalKind.value = kind
 modalVisible.value = true
}

// Modal de borrado
const { deleteItem } = useNavActions()
const deleteModalVisible = ref(false)
const deleteLevel = ref<NavLevel | null>(null)
const deleteTarget = ref<NavItem | null>(null)

const onDelete = (level: NavLevel, item: NavItem) => {
 deleteLevel.value = level
 deleteTarget.value = item
 nextTick(() => {
  deleteModalVisible.value = true
 })
}

const onDeleteClose = () => {
 deleteModalVisible.value = false
 deleteTarget.value = null
 deleteLevel.value = null
}

const onDeleteConfirm = async () => {
 if (!deleteLevel.value || !deleteTarget.value) return
 await deleteItem(deleteLevel.value, deleteTarget.value)
 onDeleteClose()
}

const onRename = (level: NavLevel, newTitle: string) => {
 if (!level.parentId) return
 nav.renameNotebook(level.parentId, newTitle)
}
</script>

<template>
 <splitpanes style="height: 100vh" :class="currentTheme === 'dark' ? 'split-dark' : 'default-theme'">
  <!-- Paneles laterales -->
  <pane :size="25" :min-size="10" :max-size="60" class="d-flex flex-column" style="min-width: 250px">
   <!-- Barra de navegación -->
   <div class="border-bottom flex-shrink-0 mainbar">
    <!-- Piso 1: título de la app (arrastrable) -->
    <div
     class="d-flex align-items-center justify-content-center toolbar-bg"
     style="height: 48px; --wails-draggable: drag"
    >
     <span class="fw-semibold" style="padding-top: 0.2rem">GVNotes</span>
    </div>
    <!-- Piso 2: controles de navegación -->
    <div class="d-flex align-items-center px-2 gap-1" style="height: 42px">
     <button class="btn border-0 d-flex align-items-center p-1" tabindex="-1" :disabled="!nav.canGoBack" @click="nav.goBack()">
      <IconArrowLeft :size="22" stroke-width="1" />
     </button>
     <div class="ms-auto d-flex align-items-center gap-1">
      <button class="btn btn-sm d-flex align-items-center p-1" tabindex="-1" @click="toggleTheme">
       <IconSun v-if="currentTheme === 'dark'" :size="22" stroke-width="1" />
       <IconMoon v-else :size="22" stroke-width="1" />
      </button>
      <button class="btn btn-sm d-flex align-items-center p-1" tabindex="-1" @click="logout">
       <IconPower :size="22" stroke-width="1" />
      </button>
     </div>
    </div>
   </div>

   <!-- Panel único (raíz sin selección) -->
   <template v-if="!nav.rightPanel">
    <NavPanel
     v-if="nav.leftPanel"
     :key="nav.leftPanel.parentId ?? 'root'"
     class="flex-grow-1"
     :level="nav.leftPanel"
     :active-note-id="nav.selectedNoteId"
     @notebook-click="onNotebookClickRight"
     @note-click="onNoteClick"
     @reorder="onReorder"
     @create="onCreate(nav.leftPanel, $event)"
     @delete="onDelete(nav.leftPanel, $event)"
     @rename="onRename"
    />
   </template>

   <!-- Dos paneles (profundidad > 1) -->
   <splitpanes v-else class="split-dark flex-grow-1">
    <pane :min-size="20" :max-size="80">
     <NavPanel
      :level="nav.leftPanel!"
      :active-note-id="nav.selectedNoteId"
      :active-notebook-id="nav.activeNotebookId"
      @notebook-click="onNotebookClickLeft"
      @note-click="onNoteClickFromLeft"
      @reorder="onReorder"
      @create="onCreate(nav.leftPanel!, $event)"
      @delete="onDelete(nav.leftPanel!, $event)"
      @rename="onRename"
     />
    </pane>
    <pane :min-size="20" :max-size="80">
     <NavPanel
      :level="nav.rightPanel"
      :active-note-id="nav.selectedNoteId"
      @notebook-click="onNotebookClickRight"
      @note-click="onNoteClick"
      @reorder="onReorder"
      @create="onCreate(nav.rightPanel!, $event)"
      @delete="onDelete(nav.rightPanel!, $event)"
      @rename="onRename"
     />
    </pane>
   </splitpanes>
  </pane>

  <!-- Contenido principal -->
  <pane :min-size="40">
   <div class="h-100 d-flex align-items-center justify-content-center" v-if="!nav.selectedNoteId">
    <p class="text-secondary mb-0">{{ t('note.select_hint') }}</p>
   </div>
   <NoteEditor
    v-else
    :key="nav.selectedNoteId"
    :note-id="nav.selectedNoteId"
    @title-changed="(id: string, title: string) => nav.renameItem(id, 'note', title)"
   />
  </pane>
 </splitpanes>

 <CreateItemModal v-model:visible="modalVisible" :kind="modalKind" :level="modalLevel" />

 <DeleteConfirmModal
  v-if="deleteTarget"
  :visible="deleteModalVisible"
  :kind="deleteTarget.kind"
  :title="deleteTarget.data.title"
  @confirm="onDeleteConfirm"
  @update:visible="onDeleteClose"
 />
</template>
