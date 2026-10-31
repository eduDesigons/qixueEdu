<script setup>
import { NImagePreview } from 'naive-ui'
import { useImagePreviewState } from '@/composables/useImagePreview'

const { show, src } = useImagePreviewState()
</script>

<template>
  <NImagePreview v-model:show="show" :src="src" />
</template>
