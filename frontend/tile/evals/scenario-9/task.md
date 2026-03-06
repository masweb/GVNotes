# TipTap Editor Extension Setup

Set up a TipTap rich text editor in a Vue 3 component using the composition API editor hook with a comprehensive set of extensions.

## Capabilities

### Core and text formatting

- The editor instance is non-null after the component is mounted [@test](../test/editor_non_null.test.ts)
- The StarterKit extension bundle is included (provides paragraph, bold, italic, heading, lists, blockquote, code, horizontal rule) [@test](../test/editor_starter_kit.test.ts)
- Underline, Highlight, TextStyle, and Color extensions are configured [@test](../test/editor_text_formatting.test.ts)
- TextAlign is configured to support `paragraph` and `heading` node types [@test](../test/editor_text_align.test.ts)

### Media and links

- The Link extension is configured with automatic link detection enabled and allows `http`, `https`, and `mailto` protocols [@test](../test/editor_link.test.ts)
- The Image extension is included [@test](../test/editor_image.test.ts)

### Structured content

- TaskList and TaskItem extensions are configured, with the TaskItem `nested` option enabled [@test](../test/editor_task_list.test.ts)
- Table, TableRow, TableHeader, and TableCell extensions are all included [@test](../test/editor_table.test.ts)

### Code blocks

- The CodeBlockLowlight extension is configured with a lowlight instance that has at least three programming languages registered [@test](../test/editor_code_block.test.ts)

### UX

- The Placeholder extension is included with a non-empty placeholder string [@test](../test/editor_placeholder.test.ts)
- The editor instance is bound to an `EditorContent` component in the template [@test](../test/editor_content_binding.test.ts)

## Implementation

[@generates](./src/components/NoteEditor.vue)

## API

```typescript { #api }
/**
 * NoteEditor component props.
 */
export interface NoteEditorProps {
  /** The identifier of the note whose content should be loaded and saved. */
  noteId: string
}

/**
 * Events emitted by NoteEditor.
 */
export interface NoteEditorEmits {
  /** Emitted when the note title changes, providing the note id and new title. */
  (e: 'title-changed', id: string, title: string): void
}
```

## Dependencies { .dependencies }

### frontend 0.0.0 { .dependency }

Vue 3 + TypeScript frontend application for GVNotes. Provides the TipTap editor initialization pattern using useEditor with StarterKit, text formatting, Link, Image, TaskList, Table, CodeBlockLowlight (lowlight), and Placeholder extensions, bound to EditorContent in the template.
