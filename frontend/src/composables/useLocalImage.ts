import BaseImage from '@tiptap/extension-image'
import { useImageServer } from './useImageServer'

// Extends the official TipTap Image extension so that stored relative src
// paths (/images/<filename>) are resolved to absolute URLs served by the
// Go image HTTP server at render time.
//
// The JSON content always stores the relative path — only the <img> element
// rendered by ProseMirror gets the absolute URL, keeping content portable
// across sessions regardless of the dynamic server port.
export const LocalImage = BaseImage.extend({
  addNodeView() {
    const { resolveImageSrc } = useImageServer()

    // Delegate to the base node view (handles ResizableNodeView when resize
    // is enabled). Patch node.attrs.src before passing to the base view so
    // the rendered img element always gets the absolute URL.
    const base = BaseImage.config.addNodeView?.bind(this)?.()
    if (base) {
      return (props) => {
        const resolvedSrc = resolveImageSrc(props.node.attrs.src ?? '')
        return base({
          ...props,
          node: {
            ...props.node,
            attrs: { ...props.node.attrs, src: resolvedSrc },
          },
        } as unknown as typeof props)
      }
    }

    // Fallback when resize is disabled: plain img element.
    return ({ node }) => {
      const img = document.createElement('img')
      img.src = resolveImageSrc(node.attrs.src ?? '')
      if (node.attrs.alt) img.alt = node.attrs.alt
      if (node.attrs.title) img.title = node.attrs.title
      if (node.attrs.width) img.style.width = `${node.attrs.width}px`
      if (node.attrs.height) img.style.height = `${node.attrs.height}px`
      return { dom: img }
    }
  },
})
