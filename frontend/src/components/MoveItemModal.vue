<script lang="ts" setup>
const { t } = useI18n()

import { IconHome } from '@tabler/icons-vue'
import { MoveNote, MoveNotebook, ListNotebooks } from '../../wailsjs/go/main/App'
import type { NavItem, NavLevel } from '@/stores/navigation'
import type { dto } from '../../wailsjs/go/models'

const props = defineProps<{
  visible: boolean
  item: NavItem | null
  level: NavLevel | null
}>()

const emit = defineEmits<{
  'update:visible': [value: boolean]
}>()

const nav = useNavigationStore()

const rootNotebooks = ref<dto.NotebookListItem[]>([])
const loadingTree = ref(false)
const loading = ref(false)
const error = ref<string | null>(null)

// null = raíz, string = id del notebook destino
const selectedDestination = ref<string | null>(null)

const currentParentId = computed(() => props.level?.parentId ?? null)

// El notebook que se está moviendo no puede ser su propio destino
const disabledId = computed(() =>
  props.item?.kind === 'notebook' ? props.item.data.id : null,
)

const isCurrentLocation = computed(() => selectedDestination.value === currentParentId.value)

const close = () => emit('update:visible', false)

const loadTree = async () => {
  loadingTree.value = true
  try {
    rootNotebooks.value = await ListNotebooks('')
  } finally {
    loadingTree.value = false
  }
}

watch(() => props.visible, async (v) => {
  if (v) {
    error.value = null
    selectedDestination.value = currentParentId.value
    await loadTree()
  }
})

const onMove = async () => {
  if (!props.item || !props.level) return
  if (isCurrentLocation.value) {
    error.value = t('move.same_location')
    return
  }

  loading.value = true
  error.value = null

  try {
    const destId = selectedDestination.value

    if (props.item.kind === 'note') {
      await MoveNote(props.item.data.id, { notebookId: destId ?? undefined })
    } else {
      await MoveNotebook(props.item.data.id, { parentId: destId ?? undefined })
    }

    nav.moveItem(props.level.parentId, props.item, destId)
    close()
  } catch (e: any) {
    error.value = e?.message ?? t('move.error')
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <CModal :visible="visible" alignment="top" @close="close">
    <CModalHeader>
      <h5 class="modal-title">
        {{ item?.kind === 'notebook' ? t('move.title_notebook') : t('move.title_note') }}
      </h5>
    </CModalHeader>

    <div class="modal-body p-2">
      <p v-if="item" class="small text-secondary mb-2 px-1">
        {{ t('move.select_destination') }}: <strong>{{ item.data.title }}</strong>
      </p>

      <div style="max-height: 55vh; overflow-y: auto">
        <ul class="list-unstyled mb-0">
          <!-- Opción raíz -->
          <li>
            <div
              class="d-flex align-items-center gap-2 px-2 py-1 rounded"
              :class="{
                'bg-primary text-white': selectedDestination === null,
                'tree-node-selectable': selectedDestination !== null,
              }"
              style="cursor: pointer"
              @click="selectedDestination = null"
            >
              <span style="width: 16px" class="flex-shrink-0" />
              <IconHome :size="16" stroke-width="1.5" class="flex-shrink-0" />
              <span class="small">{{ t('move.root') }}</span>
            </div>
          </li>

          <!-- Árbol de notebooks -->
          <li v-if="loadingTree" class="text-secondary small px-2 py-2">…</li>
          <NotebookTreeNode
            v-for="nb in rootNotebooks"
            :key="nb.id"
            :notebook="nb"
            :selected-id="selectedDestination"
            :disabled-id="disabledId"
            @select="selectedDestination = $event"
          />
        </ul>
      </div>

      <p v-if="error" class="text-danger small mt-2 mb-0 px-1">{{ error }}</p>
    </div>

    <div class="modal-footer">
      <button type="button" class="btn btn-secondary btn-sm" @click="close">
        {{ t('common.cancel') }}
      </button>
      <button
        type="button"
        class="btn btn-primary btn-sm"
        :disabled="loading || isCurrentLocation"
        @click="onMove"
      >
        {{ loading ? t('move.moving') : t('move.btn') }}
      </button>
    </div>
  </CModal>
</template>
