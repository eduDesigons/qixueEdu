import { getAssetUrl } from '@/utils'

const show = ref(false)
const src = ref('')

/** 页内图片预览（Naive ImagePreview），勿用 window.open */
export function previewImage(url) {
  if (!url)
    return
  src.value = getAssetUrl(url)
  show.value = true
}

export function useImagePreviewState() {
  return { show, src }
}
