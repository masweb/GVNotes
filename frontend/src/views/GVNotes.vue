<script lang="ts" setup>
import { Splitpanes, Pane } from 'splitpanes'
import { IconArrowLeft } from '@tabler/icons-vue'
import type { NavItem } from '@/stores/navigation'

const nav = useNavigationStore()

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
</script>

<template>
  <splitpanes style="height: 100vh" class="split-dark">
    <!-- Paneles laterales -->
    <pane :size="25" :min-size="10" :max-size="60" class="d-flex flex-column">
      <!-- Barra de navegación -->
      <div class="d-flex align-items-center px-2 py-1 border-bottom bg-body-tertiary flex-shrink-0">
        <button
          class="btn btn-sm d-flex align-items-center gap-1 p-1"
          :disabled="!nav.canGoBack"
          @click="nav.goBack()"
        >
          <IconArrowLeft :size="18" />
        </button>
      </div>

      <!-- Panel único (raíz sin selección) -->
      <template v-if="!nav.rightPanel">
        <NavPanel
          v-if="nav.leftPanel"
          class="flex-grow-1"
          :level="nav.leftPanel"
          :active-note-id="nav.selectedNoteId"
          @notebook-click="onNotebookClickRight"
          @note-click="onNoteClick"
        />
      </template>

      <!-- Dos paneles (profundidad > 1) -->
      <splitpanes v-else class="split-dark flex-grow-1">
        <pane :min-size="20" :max-size="80">
          <NavPanel
            :level="nav.leftPanel!"
            :active-note-id="null"
            :active-notebook-id="nav.activeNotebookId"
            @notebook-click="onNotebookClickLeft"
            @note-click="onNoteClick"
          />
        </pane>
        <pane :min-size="20" :max-size="80">
          <NavPanel
            :level="nav.rightPanel"
            :active-note-id="nav.selectedNoteId"
            @notebook-click="onNotebookClickRight"
            @note-click="onNoteClick"
          />
        </pane>
      </splitpanes>
    </pane>

    <!-- Contenido principal -->
    <pane :min-size="40">
      <div class="h-100 p-4 overflow-auto">
        <p v-if="!nav.selectedNoteId" class="text-secondary">Selecciona una nota</p>
        <p v-else class="text-secondary font-monospace small">note: {{ nav.selectedNoteId }}</p>
      </div>
    </pane>
  </splitpanes>
</template>

