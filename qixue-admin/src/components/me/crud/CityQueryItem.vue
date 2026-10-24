<template>
  <MeQueryItem v-if="show && !bare" label="大区" :label-width="labelWidth">
    <NSelect
      :value="modelValue || null"
      clearable
      filterable
      placeholder="全部大区"
      :options="options"
      @update:value="emit('update:modelValue', $event ?? null)"
    />
  </MeQueryItem>
  <NSelect
    v-else-if="show && bare"
    :value="modelValue || null"
    clearable
    filterable
    placeholder="全部大区"
    :options="options"
    @update:value="emit('update:modelValue', $event ?? null)"
  />
</template>

<script setup>
import { NSelect } from 'naive-ui'
import { useUserStore } from '@/store'
import api from '@/api'
import MeQueryItem from './QueryItem.vue'

defineProps({
  modelValue: { default: null },
  labelWidth: { type: Number, default: 80 },
  bare: { type: Boolean, default: false },
})
const emit = defineEmits(['update:modelValue'])

const userStore = useUserStore()
const show = computed(() => !Number(userStore.userInfo?.city_id || 0))
const options = ref([])

if (show.value) {
  api.common.getAreaTree().then(({ data = [] }) => {
    options.value = (data || []).map(item => ({
      label: item.name,
      value: Number(item.id),
    }))
  }).catch(() => {
    options.value = []
  })
}
</script>
