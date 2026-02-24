<script lang="ts" setup>
const { t } = useI18n()

import { Splitpanes, Pane } from 'splitpanes'
import { IconArrowLeft, IconSettings } from '@tabler/icons-vue'
import type { NavItem, NavLevel } from '@/stores/navigation'

const nav = useNavigationStore()
const { currentTheme } = useTheme()

onMounted(() => {
 nav.init()
 document.addEventListener('keydown', onKeydown)
})

onBeforeUnmount(() => {
 document.removeEventListener('keydown', onKeydown)
})

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

// Modal de mover
const moveModalVisible = ref(false)
const moveLevel = ref<NavLevel | null>(null)
const moveTarget = ref<NavItem | null>(null)

const onMove = (level: NavLevel, item: NavItem) => {
 moveLevel.value = level
 moveTarget.value = item
 nextTick(() => {
  moveModalVisible.value = true
 })
}

const onMoveClose = () => {
 moveModalVisible.value = false
 moveTarget.value = null
 moveLevel.value = null
}

// Modal de configuración
const settingsVisible = ref(false)

// Atajos de teclado globales
let nTimer: ReturnType<typeof setTimeout> | null = null
const onKeydown = (e: KeyboardEvent) => {
 if (!e.metaKey) return
 if (e.target instanceof HTMLInputElement || e.target instanceof HTMLTextAreaElement) return
 if (modalVisible.value || deleteModalVisible.value || settingsVisible.value || moveModalVisible.value) return
 const activeLevel = nav.rightPanel ?? nav.leftPanel
 if (!activeLevel) return
 if (e.key === 'n') {
  e.preventDefault()
  if (nTimer) {
   clearTimeout(nTimer)
   nTimer = null
   onCreate(activeLevel, 'notebook')
  } else {
   nTimer = setTimeout(() => {
    nTimer = null
    onCreate(activeLevel, 'note')
   }, 300)
  }
 }
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
     <button
      class="btn border-0 d-flex align-items-center p-1"
      tabindex="-1"
      :disabled="!nav.canGoBack"
      @click="nav.goBack()"
     >
      <IconArrowLeft :size="22" stroke-width="1" />
     </button>
     <div class="ms-auto d-flex align-items-center gap-1">
      <button class="btn btn-sm d-flex align-items-center p-1" tabindex="-1" @click="settingsVisible = true">
       <IconSettings :size="22" stroke-width="1" />
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
     @move="onMove(nav.leftPanel, $event)"
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
      @move="onMove(nav.leftPanel!, $event)"
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
      @move="onMove(nav.rightPanel!, $event)"
      @rename="onRename"
     />
    </pane>
   </splitpanes>
  </pane>

  <!-- Contenido principal -->
  <pane :min-size="40">
   <div class="h-100 d-flex flex-column align-items-center justify-content-center gap-2" v-if="!nav.selectedNoteId">
    <p class="text-secondary mb-2">{{ t('note.select_hint') }}</p>
    <p class="text-secondary mb-0" style="font-size: 1em; opacity: 0.7">
     <kbd>⌘N</kbd> {{ t('nav.new_note') }} &nbsp;·&nbsp; <kbd>⌘NN</kbd> {{ t('nav.new_notebook') }}
    </p>
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
 <AppSettingsModal v-model:visible="settingsVisible" />
 <MoveItemModal
  v-if="moveTarget"
  :visible="moveModalVisible"
  :item="moveTarget"
  :level="moveLevel"
  @update:visible="onMoveClose"
 />

 <DeleteConfirmModal
  v-if="deleteTarget"
  :visible="deleteModalVisible"
  :kind="deleteTarget.kind"
  :title="deleteTarget.data.title"
  @confirm="onDeleteConfirm"
  @update:visible="onDeleteClose"
 />
</template>
