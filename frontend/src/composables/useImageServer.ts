import { GetImageServerURL } from '../../wailsjs/go/main/App'

// Base URL of the Go image HTTP server. Resolved once at app startup.
const imageServerURL = ref('')

export const useImageServer = () => {
  const init = async () => {
    if (!imageServerURL.value) {
      imageServerURL.value = await GetImageServerURL()
    }
  }

  // Converts a stored relative image path (/images/uuid.ext) to an absolute
  // URL served by the Go image server.
  const resolveImageSrc = (src: string): string => {
    if (src.startsWith('/images/') && imageServerURL.value) {
      return `${imageServerURL.value}${src}`
    }
    return src
  }

  return { imageServerURL, init, resolveImageSrc }
}
