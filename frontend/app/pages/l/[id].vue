<script setup lang="ts">
import { useRoute } from "vue-router"

import CategoryTabs from "~/components/catalog/category-tabs.vue"
import CatalogItemCard from "~/components/catalog/item-card.vue"

const route = useRoute()
const locationId = route.params.id as string

const {
  location,
  categories,
  selectedCategory,
  filteredItems,
} = useCatalog(locationId)

useHead({
  title: location.value ? `Order from ${location.value.name}` : "Foodstock | Vending",
  meta: [
    { name: "description", content: "Fresh organic food from our premium vending sanctuaries." },
  ],
})
</script>

<template lang="pug">
div(class="flex flex-col gap-10")
  section(v-if="location" class="flex flex-col gap-2")
    div(class="text-xs font-bold uppercase tracking-widest text-primary opacity-80") Foodstock
    h1(class="display-lg font-extrabold") {{ location.name }}
    p(class="body-md opacity-60 flex items-center gap-2")
      u-icon(name="i-heroicons-map-pin" class="w-4 h-4")
      | {{ location.address }}

  section(v-else class="flex flex-col gap-4 py-20 items-center text-center")
    u-icon(name="i-heroicons-exclamation-triangle" class="w-12 h-12 text-primary opacity-20")
    h1(class="headline-lg") Machine Not Found
    p(class="body-md") The QR code might be outdated or incorrect.

  div(v-if="location" class="flex flex-col gap-8")
    category-tabs(
      v-model="selectedCategory"
      :categories="categories"
    )

    div(class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-8")
      transition-group(name="list")
        catalog-item-card(
          v-for="item in filteredItems"
          :key="item.id"
          :item="item"
        )
</template>

<style scoped>
.list-enter-active,
.list-leave-active {
  transition: all 0.5s ease;
}

.list-enter-from,
.list-leave-to {
  opacity: 0;
  transform: translateY(20px);
}
</style>
