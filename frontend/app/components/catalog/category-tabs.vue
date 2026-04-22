<script setup lang="ts">
import { CATEGORY_LABELS } from "~/composables/catalog"

defineProps<{
  categories: string[]
  modelValue: string
}>()

const emit = defineEmits<{
  "update:modelValue": [string]
}>()
</script>

<template lang="pug">
div(class="flex gap-3 overflow-x-auto no-scrollbar py-2 -mx-4 px-4 md:mx-0 md:px-0")
  button(
    v-for="category in categories"
    :key="category"
    class="px-5 py-2.5 rounded-full whitespace-nowrap transition-all duration-300 transform active:scale-95"
    :class="modelValue === category ? 'btn-primary' : 'btn-secondary'"
    @click="emit('update:modelValue', category)"
  )
    | {{ CATEGORY_LABELS[category] || category }}
</template>

<style scoped>
.no-scrollbar::-webkit-scrollbar {
  display: none;
}

.no-scrollbar {
  -ms-overflow-style: none;
  scrollbar-width: none;
}
</style>
