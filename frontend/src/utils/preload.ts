/**
 * Preload image utility
 * @param url - Image URL to preload
 * @returns Promise that resolves when image is loaded
 */
export function preloadImage(url: string): Promise<void> {
  return new Promise((resolve, reject) => {
    const img = new Image()
    img.onload = () => resolve()
    img.onerror = () => reject(new Error(`Failed to load image: ${url}`))
    img.src = url
  })
}

/**
 * Preload multiple images with progress tracking
 * @param urls - Array of image URLs to preload
 * @param onProgress - Optional callback to track loading progress (0-100)
 * @returns Promise that resolves when all images are loaded
 */
export async function preloadImages(
  urls: string[],
  onProgress?: (progress: number) => void
): Promise<void> {
  if (urls.length === 0) {
    onProgress?.(100)
    return
  }

  let loaded = 0
  const total = urls.length

  await Promise.all(
    urls.map(async (url) => {
      try {
        await preloadImage(url)
      } catch (error) {
        console.error('Failed to preload image:', url, error)
        // Continue loading other images even if one fails
      } finally {
        loaded++
        const progress = Math.round((loaded / total) * 100)
        onProgress?.(progress)
      }
    })
  )
}

/**
 * Preload audio file
 * @param url - Audio URL to preload
 * @returns Promise that resolves when audio is loaded
 */
export function preloadAudio(url: string): Promise<void> {
  return new Promise((resolve, reject) => {
    const audio = new Audio()
    audio.oncanplaythrough = () => resolve()
    audio.onerror = () => reject(new Error(`Failed to load audio: ${url}`))
    audio.src = url
  })
}

/**
 * Preload multiple assets (images and audio) with progress tracking
 * @param imageUrls - Array of image URLs
 * @param audioUrls - Array of audio URLs
 * @param onProgress - Optional callback to track loading progress (0-100)
 * @returns Promise that resolves when all assets are loaded
 */
export async function preloadAssets(
  imageUrls: string[] = [],
  audioUrls: string[] = [],
  onProgress?: (progress: number) => void
): Promise<void> {
  const allUrls = [...imageUrls, ...audioUrls]

  if (allUrls.length === 0) {
    onProgress?.(100)
    return
  }

  let loaded = 0
  const total = allUrls.length

  const loadPromises = [
    ...imageUrls.map(async (url) => {
      try {
        await preloadImage(url)
      } catch (error) {
        console.error('Failed to preload image:', url, error)
      } finally {
        loaded++
        onProgress?.(Math.round((loaded / total) * 100))
      }
    }),
    ...audioUrls.map(async (url) => {
      try {
        await preloadAudio(url)
      } catch (error) {
        console.error('Failed to preload audio:', url, error)
      } finally {
        loaded++
        onProgress?.(Math.round((loaded / total) * 100))
      }
    })
  ]

  await Promise.all(loadPromises)
}
