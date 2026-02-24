<script lang="ts" setup>
import { IconChevronRight, IconChevronDown, IconNotebook } from '@tabler/icons-vue'
import { ListNotebooks } from '../../wailsjs/go/main/App'
import type { dto } from '../../wailsjs/go/models'

const props = defineProps<{
  notebook: dto.NotebookListItem
  selectedId: string | null
  disabledId: string | null
}>()

const emit = defineEmits<{
  select: [id: string]
}>()

const expanded = ref(false)
const children = ref<dto.NotebookListItem[]>([])
const loadingChildren = ref(false)

const isSelected = computed(() => props.selectedId === props.notebook.id)
const isDisabled = computed(() => props.disabledId === props.notebook.id)

const toggle = async (e: Event) => {
  e.stopPropagation()
  if (isDisabled.value) return
  expanded.value = !expanded.value
  if (expanded.value && children.value.length === 0) {
    loadingChildren.value = true
    try {
      children.value = await ListNotebooks(props.notebook.id)
    } finally {
      loadingChildren.value = false
    }
  }
}

const select = () => {
  if (isDisabled.value) return
  emit('select', props.notebook.id)
}
</script>

<template>
  <li class="list-unstyled mb-0">
    <div
      class="d-flex align-items-center gap-1 px-2 py-1 rounded"
      :class="{
        'bg-primary text-white': isSelected,
        'text-secondary': isDisabled,
        'tree-node-selectable': !isDisabled && !isSelected,
      }"
      style="cursor: pointer"
      @click="select"
    >
      <button
        class="btn btn-sm p-0 border-0 d-flex align-items-center flex-shrink-0"
        style="width: 16px"
        type="button"
        :disabled="isDisabled"
        tabindex="-1"
        @click="toggle"
      >
        <IconChevronDown v-if="expanded" :size="14" />
        <IconChevronRight v-else :size="14" />
      </button>
      <IconNotebook :size="16" stroke-width="1.5" class="flex-shrink-0" />
      <span class="text-truncate small flex-grow-1">{{ notebook.title }}</span>
    </div>
    <ul v-if="expanded" class="list-unstyled ps-3 mb-0">
      <li v-if="loadingChildren" class="text-secondary small px-2 py-1">…</li>
      <NotebookTreeNode
        v-for="child in children"
        :key="child.id"
        :notebook="child"
        :selected-id="selectedId"
        :disabled-id="disabledId"
        @select="emit('select', $event)"
      />
    </ul>
  </li>
</template>
