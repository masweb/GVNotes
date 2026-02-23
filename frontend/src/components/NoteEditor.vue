<script lang="ts" setup>
import Bold from '@tiptap/extension-bold'
import CharacterCount from '@tiptap/extension-character-count'
import Color from '@tiptap/extension-color'
import Document from '@tiptap/extension-document'
import Heading from '@tiptap/extension-heading'
import History from '@tiptap/extension-history'
import { LocalImage } from '../composables/useLocalImage'
import { useImageServer } from '../composables/useImageServer'
import Italic from '@tiptap/extension-italic'
import Paragraph from '@tiptap/extension-paragraph'
import Strike from '@tiptap/extension-strike'
import Text from '@tiptap/extension-text'
import TextAlign from '@tiptap/extension-text-align'
import { TextStyle } from '@tiptap/extension-text-style'
import { Editor, EditorContent } from '@tiptap/vue-3'
import {
 IconAlignCenter,
 IconAlignJustified,
 IconAlignLeft,
 IconAlignRight,
 IconArrowBackUp,
 IconArrowForwardUp,
 IconBold,
 IconH1,
 IconH2,
 IconH3,
 IconH4,
 IconH5,
 IconH6,
 IconItalic,
 IconPalette,
 IconPhoto,
 IconStrikethrough
} from '@tabler/icons-vue'
import { GetNote, SaveImage, UpdateNoteContent, UpdateNoteTitle } from '../../wailsjs/go/main/App'
import type { dto } from '../../wailsjs/go/models'

const props = defineProps<{ noteId: string }>()
const emit = defineEmits<{ titleChanged: [id: string, title: string] }>()

const note = ref<dto.NoteDetail | null>(null)
const loading = ref(false)
const saving = ref(false)
const saveTimer = ref<ReturnType<typeof setTimeout> | null>(null)

// Renombrar título
const editingTitle = ref(false)
const titleInput = ref<HTMLInputElement | null>(null)
const { handleSubmit: handleTitleSubmit, resetForm: resetTitleForm } = useForm({ validateOnMount: false })
const { value: titleValue, errorMessage: titleError } = useField<string>('title', 'required', {
 validateOnValueUpdate: false
})

const startEditTitle = () => {
 if (!note.value) return
 titleValue.value = note.value.title
 editingTitle.value = true
 nextTick(() => {
  titleInput.value?.select()
 })
}

const cancelEditTitle = () => {
 editingTitle.value = false
 resetTitleForm()
}

// Color picker
const colorOpen = ref(false)
const colorBtn = ref<HTMLElement | null>(null)
const colorMenuStyle = ref<Record<string, string>>({})

const COLORS = [
 '#000000',
 '#374151',
 '#6b7280',
 '#9ca3af',
 '#ffffff',
 '#ef4444',
 '#f97316',
 '#eab308',
 '#22c55e',
 '#3b82f6',
 '#8b5cf6',
 '#ec4899',
 '#06b6d4',
 '#14b8a6',
 '#84cc16',
 '#dc2626',
 '#ea580c',
 '#ca8a04',
 '#16a34a',
 '#2563eb'
]

const toggleColorPicker = () => {
 if (!colorOpen.value && colorBtn.value) {
  const rect = colorBtn.value.getBoundingClientRect()
  colorMenuStyle.value = {
   position: 'fixed',
   top: `${rect.bottom + 4}px`,
   left: `${rect.left}px`,
   zIndex: '9999'
  }
 }
 colorOpen.value = !colorOpen.value
}

const pickColor = (color: string) => {
 editor.chain().focus().setColor(color).run()
 colorOpen.value = false
}

const clearColor = () => {
 editor.chain().focus().unsetColor().run()
 colorOpen.value = false
}

const onDocClickColor = (e: MouseEvent) => {
 if (colorBtn.value && !colorBtn.value.contains(e.target as Node)) {
  colorOpen.value = false
 }
}

const { init: initImageServer } = useImageServer()

onMounted(async () => {
 await initImageServer()
 document.addEventListener('click', onDocClickColor)
})
onBeforeUnmount(() => document.removeEventListener('click', onDocClickColor))

// Insertar imagen
const insertImage = () => {
 const input = document.createElement('input')
 input.type = 'file'
 input.accept = 'image/png,image/jpeg,image/gif,image/webp'
 input.onchange = async () => {
  const file = input.files?.[0]
  if (!file || !note.value) return
  const buffer = await file.arrayBuffer()
  const bytes = Array.from(new Uint8Array(buffer))
  const img = await SaveImage(note.value.id, file.type, bytes)
  editor.chain().focus().setImage({ src: `/images/${img.filename}` }).run()
 }
 input.click()
}

const submitTitle = handleTitleSubmit(async values => {
 if (!note.value || values.title === note.value.title) {
  cancelEditTitle()
  return
 }
 await UpdateNoteTitle(note.value.id, { title: values.title })
 note.value.title = values.title
 emit('titleChanged', note.value.id, values.title)
 cancelEditTitle()
})

const editor = new Editor({
 content: '',
 extensions: [
  Bold,
  CharacterCount,
  Color,
  Document,
  Heading.configure({ levels: [1, 2, 3, 4, 5, 6] }),
  History,
  LocalImage.configure({
   inline: false,
   allowBase64: false,
   resize: {
    enabled: true,
    directions: ['bottom-right', 'bottom-left', 'bottom'],
    minWidth: 80,
    minHeight: 40,
    alwaysPreserveAspectRatio: true,
   },
  }),
  Italic,
  Paragraph,
  Strike,
  Text,
  TextAlign.configure({ types: ['heading', 'paragraph'] }),
  TextStyle
 ],
 onUpdate: () => {
  scheduleSave()
 }
})

const scheduleSave = () => {
 if (saveTimer.value) clearTimeout(saveTimer.value)
 saveTimer.value = setTimeout(async () => {
  if (!note.value) return
  saving.value = true
  try {
   await UpdateNoteContent(note.value.id, { content: JSON.stringify(editor.getJSON()) })
   await new Promise(r => setTimeout(r, 600))
  } finally {
   saving.value = false
  }
 }, 800)
}

const loadNote = async (id: string) => {
 loading.value = true
 try {
  note.value = await GetNote(id)
  editor.commands.setContent(note.value.content ? JSON.parse(note.value.content) : '')
 } finally {
  loading.value = false
 }
}

watch(
 () => props.noteId,
 id => {
  loadNote(id)
 },
 { immediate: true }
)

onBeforeUnmount(() => {
 if (saveTimer.value) clearTimeout(saveTimer.value)
 editor.destroy()
})
</script>

<template>
 <div class="note-editor h-100 d-flex flex-column">
  <!-- Toolbar -->
  <div
   class="editor-toolbar border-0 d-flex align-items-center flex-wrap gap-1 px-2 py-1 flex-shrink-0"
  >
   <!-- History -->
   <div class="d-flex">
    <button
     type="button"
     class="btn btn-sm"
     :disabled="!editor.can().undo()"
     @click="editor.chain().focus().undo().run()"
    >
     <IconArrowBackUp :size="22" stroke-width="1" />
    </button>
    <button
     type="button"
     class="btn btn-sm"
     :disabled="!editor.can().redo()"
     @click="editor.chain().focus().redo().run()"
    >
     <IconArrowForwardUp :size="22" stroke-width="1" />
    </button>
   </div>

   <!-- Formato básico -->
   <div class="d-flex">
    <button
     type="button"
     class="btn btn-sm"
     :class="{ active: editor.isActive('bold') }"
     @click="editor.chain().focus().toggleBold().run()"
    >
     <IconBold :size="22" stroke-width="1" />
    </button>
    <button
     type="button"
     class="btn btn-sm"
     :class="{ active: editor.isActive('italic') }"
     @click="editor.chain().focus().toggleItalic().run()"
    >
     <IconItalic :size="22" stroke-width="1" />
    </button>
    <button
     type="button"
     class="btn btn-sm"
     :class="{ active: editor.isActive('strike') }"
     @click="editor.chain().focus().toggleStrike().run()"
    >
     <IconStrikethrough :size="22" stroke-width="1" />
    </button>
   </div>

   <!-- Encabezados -->
   <div class="d-flex">
    <button
     type="button"
     class="btn btn-sm"
     :class="{ active: editor.isActive('heading', { level: 1 }) }"
     @click="editor.chain().focus().toggleHeading({ level: 1 }).run()"
    >
     <IconH1 :size="22" stroke-width="1" />
    </button>
    <button
     type="button"
     class="btn btn-sm"
     :class="{ active: editor.isActive('heading', { level: 2 }) }"
     @click="editor.chain().focus().toggleHeading({ level: 2 }).run()"
    >
     <IconH2 :size="22" stroke-width="1" />
    </button>
    <button
     type="button"
     class="btn btn-sm"
     :class="{ active: editor.isActive('heading', { level: 3 }) }"
     @click="editor.chain().focus().toggleHeading({ level: 3 }).run()"
    >
     <IconH3 :size="22" stroke-width="1" />
    </button>
    <button
     type="button"
     class="btn btn-sm"
     :class="{ active: editor.isActive('heading', { level: 4 }) }"
     @click="editor.chain().focus().toggleHeading({ level: 4 }).run()"
    >
     <IconH4 :size="22" stroke-width="1" />
    </button>
    <button
     type="button"
     class="btn btn-sm"
     :class="{ active: editor.isActive('heading', { level: 5 }) }"
     @click="editor.chain().focus().toggleHeading({ level: 5 }).run()"
    >
     <IconH5 :size="22" stroke-width="1" />
    </button>
    <button
     type="button"
     class="btn btn-sm"
     :class="{ active: editor.isActive('heading', { level: 6 }) }"
     @click="editor.chain().focus().toggleHeading({ level: 6 }).run()"
    >
     <IconH6 :size="22" stroke-width="1" />
    </button>
   </div>

   <!-- Alineación -->
   <div class="d-flex">
    <button
     type="button"
     class="btn btn-sm"
     :class="{ active: editor.isActive({ textAlign: 'left' }) }"
     @click="editor.chain().focus().setTextAlign('left').run()"
    >
     <IconAlignLeft :size="22" stroke-width="1" />
    </button>
    <button
     type="button"
     class="btn btn-sm"
     :class="{ active: editor.isActive({ textAlign: 'center' }) }"
     @click="editor.chain().focus().setTextAlign('center').run()"
    >
     <IconAlignCenter :size="22" stroke-width="1" />
    </button>
    <button
     type="button"
     class="btn btn-sm"
     :class="{ active: editor.isActive({ textAlign: 'right' }) }"
     @click="editor.chain().focus().setTextAlign('right').run()"
    >
     <IconAlignRight :size="22" stroke-width="1" />
    </button>
    <button
     type="button"
     class="btn btn-sm"
     :class="{ active: editor.isActive({ textAlign: 'justify' }) }"
     @click="editor.chain().focus().setTextAlign('justify').run()"
    >
     <IconAlignJustified :size="22" stroke-width="1" />
    </button>
   </div>

   <!-- Color -->
   <div class="d-flex">
    <button
     ref="colorBtn"
     type="button"
     class="btn btn-sm"
     :class="{ active: !!editor.getAttributes('textStyle').color }"
     @click.stop="toggleColorPicker"
    >
     <IconPalette :size="22" stroke-width="1" :style="{ color: editor.getAttributes('textStyle').color }" />
    </button>
   </div>

   <!-- Imagen -->
   <div class="d-flex">
    <button type="button" class="btn btn-sm" :disabled="!note" @click="insertImage">
     <IconPhoto :size="22" stroke-width="1" />
    </button>
   </div>
   <Teleport to="body">
    <div v-if="colorOpen" :style="colorMenuStyle" class="color-picker-menu border rounded shadow-sm bg-body p-2">
     <div class="d-flex flex-wrap gap-1" style="width: 134px">
      <button
       v-for="c in COLORS"
       :key="c"
       type="button"
       class="color-swatch"
       :style="{
        background: c,
        outline:
         editor.getAttributes('textStyle').color === c
          ? '2px solid var(--bs-primary)'
          : '1px solid var(--bs-border-color)'
       }"
       @click="pickColor(c)"
      />
     </div>
     <button type="button" class="btn btn-sm btn-link px-0 mt-1 text-secondary small" @click="clearColor">
      Sin color
     </button>
    </div>
   </Teleport>
  </div>

  <Teleport to="body">
   <Transition name="fade">
    <span v-if="saving" class="saving-badge position-fixed text-secondary small">Guardando…</span>
   </Transition>
  </Teleport>

  <!-- Título -->
  <div v-if="note" class="px-4 pt-3 pb-1 flex-shrink-0">
   <form v-if="editingTitle" @submit.prevent="submitTitle" @keydown.esc="cancelEditTitle">
    <input
     ref="titleInput"
     v-model="titleValue"
     type="text"
     class="form-control form-control-sm fw-semibold fs-5 border-0 border-bottom rounded-0 px-0"
     :class="{ 'is-invalid': titleError }"
     @blur="submitTitle"
    />
    <div v-if="titleError" class="invalid-feedback">{{ titleError }}</div>
   </form>
   <h5 v-else class="mb-0 text-truncate cursor-pointer" title="Click para renombrar" @click="startEditTitle">
    {{ note.title }}
   </h5>
  </div>

  <!-- Contenido -->
  <div class="editor-content flex-grow-1 overflow-auto px-4 py-2" @click.self="editor.commands.focus()">
   <div v-if="loading" class="text-secondary small mt-3">Cargando…</div>
   <EditorContent v-else :editor="editor" />
  </div>

  <!-- Footer -->
  <div class="border-top px-3 py-1 text-end text-secondary small flex-shrink-0">
   <span class="me-3">{{ editor.storage.characterCount.characters() }} caracteres</span>
   <span>{{ editor.storage.characterCount.words() }} palabras</span>
  </div>
 </div>
</template>

<style>
.tiptap.ProseMirror:focus,
.tiptap.ProseMirror:focus-visible {
 outline: none !important;
}

.tiptap img {
 max-width: 100%;
 height: auto;
 border-radius: 4px;
 display: block;
}

/* ResizableNodeView — container */
[data-resize-container] {
 display: inline-flex !important;
 max-width: 100%;
 margin: 4px 0;
}

/* Selección de la imagen */
[data-resize-container].ProseMirror-selectednode [data-resize-wrapper] {
 outline: 2px solid var(--bs-primary);
 border-radius: 4px;
}

/* Handles */
[data-resize-handle] {
 width: 10px;
 height: 10px;
 background: var(--bs-primary);
 border: 2px solid var(--bs-body-bg);
 border-radius: 50%;
 opacity: 0;
 transition: opacity 0.15s;
}

[data-resize-container].ProseMirror-selectednode [data-resize-handle] {
 opacity: 1;
}

/* Cursores por dirección */
[data-resize-handle="bottom-right"] { cursor: nwse-resize; }
[data-resize-handle="bottom-left"]  { cursor: nesw-resize; }
[data-resize-handle="bottom"]       { cursor: s-resize; width: 100%; height: 6px; border-radius: 0; }


.saving-badge {
 top: 12px;
 right: 16px;
}

.fade-enter-active,
.fade-leave-active {
 transition: opacity 0.3s;
}
.fade-enter-from,
.fade-leave-to {
 opacity: 0;
}

.color-swatch {
 width: 22px;
 height: 22px;
 border-radius: 3px;
 border: none;
 cursor: pointer;
 padding: 0;
 flex-shrink: 0;
}

.color-swatch:hover {
 transform: scale(1.15);
}
</style>
