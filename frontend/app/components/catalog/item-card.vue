<script setup lang="ts">
defineProps<{
  item: CatalogItem
}>()
</script>

<template lang="pug">
div(class="surface-card p-4 flex flex-col gap-4 relative transition-all duration-300 hover:-translate-y-1 hover:shadow-xl")
  div(class="ingredient-float aspect-square overflow-hidden rounded-md bg-surface-container-low")
    img(
      :src="item.photo_url"
      :alt="item.name"
      class="w-full h-full object-cover food-image"
    )

  div(class="flex flex-col gap-1")
    div(class="flex justify-between items-start")
      h3(class="headline-md text-on-surface line-clamp-1") {{ item.name }}
      div(class="text-primary font-bold text-lg") {{ formatNumber(item.price / 100) }} ₽

    p(class="body-md text-on-surface opacity-70 line-clamp-2") {{ item.description }}

  div(class="mt-auto flex items-center justify-between")
    div(v-if="item.stock_amount > 0" class="flex items-center gap-2")
      div(class="w-2 h-2 rounded-full bg-primary")
      span(class="text-xs font-medium uppercase tracking-wider text-on-surface opacity-60") {{ item.stock_amount }} in stock
    div(v-else class="text-xs font-medium uppercase text-red-500") Out of stock

    button(class="btn-primary px-6 py-2 text-sm") Add
</template>

<style scoped>
.line-clamp-1 {
  display: -webkit-box;
  -webkit-line-clamp: 1;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>
