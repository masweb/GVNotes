<script lang="ts" setup>
const { t } = useI18n()

import Bold from '@tiptap/extension-bold'
import CharacterCount from '@tiptap/extension-character-count'
import Color from '@tiptap/extension-color'
import Document from '@tiptap/extension-document'
import Heading from '@tiptap/extension-heading'
import History from '@tiptap/extension-history'
import Image from '@tiptap/extension-image'
import Italic from '@tiptap/extension-italic'
import Link from '@tiptap/extension-link'
import Underline from '@tiptap/extension-underline'
import Paragraph from '@tiptap/extension-paragraph'
import Strike from '@tiptap/extension-strike'
import Text from '@tiptap/extension-text'
import TextAlign from '@tiptap/extension-text-align'
import { TextStyle } from '@tiptap/extension-text-style'
import { BubbleMenu } from '@tiptap/vue-3/menus'
import { Editor, EditorContent } from '@tiptap/vue-3'
import {
 IconAlignCenter,
 IconAlignJustified,
 IconAlignLeft,
 IconAlignRight,
 IconArrowBackUp,
 IconArrowForwardUp,
 IconBold,
 IconExternalLink,
 IconHeading,
 IconH1,
 IconH2,
 IconH3,
 IconH4,
 IconH5,
 IconH6,
 IconItalic,
 IconLink,
 IconLinkOff,
 IconPalette,
 IconPhoto,
 IconStrikethrough,
 IconUnderline
} from '@tabler/icons-vue'
import { BrowserOpenURL } from '../../wailsjs/runtime/runtime'
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
// Headings dropdown
const headingOpen = ref(false)
const headingBtn = ref<HTMLElement | null>(null)
const headingMenuStyle = ref<Record<string, string>>({})

const toggleHeadingDropdown = () => {
 if (!headingOpen.value && headingBtn.value) {
  const rect = headingBtn.value.getBoundingClientRect()
  headingMenuStyle.value = {
   position: 'fixed',
   top: `${rect.bottom + 4}px`,
   left: `${rect.left}px`,
   zIndex: '9999'
  }
 }
 headingOpen.value = !headingOpen.value
}

const pickHeading = (level: 1 | 2 | 3 | 4 | 5 | 6) => {
 editor.chain().focus().toggleHeading({ level }).run()
 headingOpen.value = false
}

const onDocClickHeading = (e: MouseEvent) => {
 if (headingBtn.value && !headingBtn.value.contains(e.target as Node)) {
  headingOpen.value = false
 }
}

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

// Link
const linkOpen = ref(false)
const linkInput = ref<HTMLInputElement | null>(null)
const linkUrl = ref('')
const linkBtn = ref<HTMLElement | null>(null)
const linkPopover = ref<HTMLElement | null>(null)
const linkMenuStyle = ref<Record<string, string>>({})
let savedLinkRange: { from: number; to: number } | null = null

const openLinkModal = () => {
 const { from, to } = editor.state.selection
 savedLinkRange = { from, to }
 linkUrl.value = editor.getAttributes('link').href ?? ''
 if (linkBtn.value) {
  const rect = linkBtn.value.getBoundingClientRect()
  linkMenuStyle.value = {
   position: 'fixed',
   top: `${rect.bottom + 4}px`,
   left: `${rect.left}px`,
   zIndex: '9999'
  }
 }
 linkOpen.value = true
 nextTick(() => {
  linkInput.value?.focus()
  linkInput.value?.select()
 })
}

const applyLink = () => {
 const url = linkUrl.value.trim()
 if (!url || !savedLinkRange) return
 const { from, to } = savedLinkRange
 savedLinkRange = null
 linkOpen.value = false
 // Transacción atómica de ProseMirror: restaurar selección + aplicar mark
 const { state, dispatch } = editor.view
 const markType = state.schema.marks['link']
 if (!markType) return
 const { tr } = state
 tr.addMark(from, to, markType.create({ href: url }))
 dispatch(tr)
 editor.commands.focus()
}

const removeLink = () => {
 editor.chain().focus().unsetLink().run()
}

const openLinkHref = () => {
 const href = editor.getAttributes('link').href
 if (href) BrowserOpenURL(href)
}

const onDocClickLink = (e: MouseEvent) => {
 const target = e.target as Node
 if (linkBtn.value?.contains(target)) return
 if (linkPopover.value?.contains(target)) return
 linkOpen.value = false
}

onMounted(() => {
 document.addEventListener('click', onDocClickColor)
 document.addEventListener('click', onDocClickLink)
 document.addEventListener('click', onDocClickHeading)
})

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
  editor
   .chain()
   .focus()
   .setImage({ src: `/images/${img.filename}` })
   .run()
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
  Image.configure({
   inline: false,
   allowBase64: false,
   resize: {
    enabled: true,
    directions: ['bottom-right', 'bottom-left', 'bottom'],
    minWidth: 80,
    minHeight: 40,
    alwaysPreserveAspectRatio: true
   }
  }),
  Italic,
  Link.configure({ openOnClick: false }),
  Paragraph,
  Strike,
  Text,
  Underline,
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
  nextTick(() => { editor.commands.focus() })
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
 document.removeEventListener('click', onDocClickColor)
 document.removeEventListener('click', onDocClickLink)
 document.removeEventListener('click', onDocClickHeading)
 if (saveTimer.value) clearTimeout(saveTimer.value)
 editor.destroy()
})
</script>

<template>
 <div class="note-editor h-100 d-flex flex-column">
  <!-- Toolbar -->
  <div class="editor-toolbar border-0 d-flex align-items-center flex-wrap gap-1 px-2 py-1 flex-shrink-0">
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
    <button
     type="button"
     class="btn btn-sm"
     :class="{ active: editor.isActive('underline') }"
     @click="editor.chain().focus().toggleUnderline().run()"
    >
     <IconUnderline :size="22" stroke-width="1" />
    </button>
   </div>

   <!-- Encabezados -->
   <div class="d-flex">
    <button
     ref="headingBtn"
     type="button"
     class="btn btn-sm"
     :class="{ active: editor.isActive('heading') }"
     @click.stop="toggleHeadingDropdown"
    >
     <IconHeading :size="22" stroke-width="1" /><span style="font-size: 9px; line-height: 1; margin-left: 1px;">▾</span>
    </button>
   </div>
   <Teleport to="body">
    <div v-if="headingOpen" :style="headingMenuStyle" class="heading-dropdown border rounded shadow-sm bg-body p-1">
     <button
      type="button"
      class="btn btn-sm w-100 text-start d-flex align-items-center gap-2"
      :class="{ active: editor.isActive('heading', { level: 1 }) }"
      @click="pickHeading(1)"
     >
      <IconH1 :size="20" stroke-width="1" />
     </button>
     <button
      type="button"
      class="btn btn-sm w-100 text-start d-flex align-items-center gap-2"
      :class="{ active: editor.isActive('heading', { level: 2 }) }"
      @click="pickHeading(2)"
     >
      <IconH2 :size="20" stroke-width="1" />
     </button>
     <button
      type="button"
      class="btn btn-sm w-100 text-start d-flex align-items-center gap-2"
      :class="{ active: editor.isActive('heading', { level: 3 }) }"
      @click="pickHeading(3)"
     >
      <IconH3 :size="20" stroke-width="1" />
     </button>
     <button
      type="button"
      class="btn btn-sm w-100 text-start d-flex align-items-center gap-2"
      :class="{ active: editor.isActive('heading', { level: 4 }) }"
      @click="pickHeading(4)"
     >
      <IconH4 :size="20" stroke-width="1" />
     </button>
     <button
      type="button"
      class="btn btn-sm w-100 text-start d-flex align-items-center gap-2"
      :class="{ active: editor.isActive('heading', { level: 5 }) }"
      @click="pickHeading(5)"
     >
      <IconH5 :size="20" stroke-width="1" />
     </button>
     <button
      type="button"
      class="btn btn-sm w-100 text-start d-flex align-items-center gap-2"
      :class="{ active: editor.isActive('heading', { level: 6 }) }"
      @click="pickHeading(6)"
     >
      <IconH6 :size="20" stroke-width="1" />
     </button>
    </div>
   </Teleport>

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

   <!-- Link -->
   <div class="d-flex">
    <button
     ref="linkBtn"
     type="button"
     class="btn btn-sm"
     :class="{ active: editor.isActive('link') }"
     @click.stop="openLinkModal"
    >
     <IconLink :size="22" stroke-width="1" />
    </button>
    <button type="button" class="btn btn-sm" :disabled="!editor.isActive('link')" @click="removeLink">
     <IconLinkOff :size="22" stroke-width="1" />
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
    <div
     v-if="linkOpen"
     ref="linkPopover"
     :style="linkMenuStyle"
     class="link-modal border rounded shadow-sm bg-body p-2"
    >
     <form class="d-flex gap-1" @submit.prevent="applyLink">
      <input
       ref="linkInput"
       v-model="linkUrl"
       type="text"
       class="form-control form-control-sm"
       placeholder="https://..."
       style="min-width: 220px"
      />
      <button type="submit" class="btn btn-sm btn-primary">OK</button>
     </form>
    </div>
   </Teleport>
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
      {{ t('color.none') }}
     </button>
    </div>
   </Teleport>
  </div>

  <Teleport to="body">
   <Transition name="fade">
    <span v-if="saving" class="saving-badge position-fixed text-secondary small">{{ t('note.saving') }}</span>
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
   <h5 v-else class="mb-0 text-truncate cursor-pointer" :title="t('nav.rename_hint')" @click="startEditTitle">
    {{ note.title }}
   </h5>
  </div>

  <!-- Contenido -->
  <div class="editor-content flex-grow-1 overflow-auto px-4 py-2" @click.self="editor.commands.focus()">
   <div v-if="loading" class="text-secondary small mt-3">{{ t('note.loading') }}</div>
   <EditorContent v-else :editor="editor" />
   <BubbleMenu
    :editor="editor"
    :should-show="() => editor.isActive('link')"
    class="link-bubble border rounded shadow-sm bg-body px-2 py-1 d-flex align-items-center gap-2"
   >
    <span class="text-truncate small" style="max-width: 200px">{{ editor.getAttributes('link').href }}</span>
    <a
     :href="editor.getAttributes('link').href"
     target="_blank"
     rel="noopener noreferrer"
     class="btn btn-sm p-0 text-primary"
     :title="t('link.open')"
     @click.prevent="openLinkHref"
    >
     <IconExternalLink :size="16" stroke-width="1.5" />
    </a>
    <button type="button" class="btn btn-sm p-0 text-secondary" :title="t('link.edit')" @click="openLinkModal">
     <IconLink :size="16" stroke-width="1.5" />
    </button>
    <button type="button" class="btn btn-sm p-0 text-danger" :title="t('link.remove')" @click="removeLink">
     <IconLinkOff :size="16" stroke-width="1.5" />
    </button>
   </BubbleMenu>
  </div>

  <!-- Footer -->
  <div class="border-top px-3 py-1 text-end text-secondary small flex-shrink-0">
   <span class="me-3">{{ editor.storage.characterCount.characters() }} {{ t('note.characters') }}</span>
   <span>{{ editor.storage.characterCount.words() }} {{ t('note.words') }}</span>
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
[data-resize-handle='bottom-right'] {
 cursor: nwse-resize;
}
[data-resize-handle='bottom-left'] {
 cursor: nesw-resize;
}
[data-resize-handle='bottom'] {
 cursor: s-resize;
 width: 100%;
 height: 6px;
 border-radius: 0;
}

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

.link-bubble {
 font-size: 0.8rem;
}
</style>
