<script setup lang="ts">
import { useRoute } from "vue-router"

import type { CatalogItem } from "~/types/catalog"

const route = useRoute()
const locationSlug = route.params.slug as string

const {
  location,
  categories,
  selectedCategory,
  selectedItem,
  filters,
  isFiltersActive,
  resetFilters,
  filteredItems,
} = useCatalog(locationSlug)

const isDetailOpen = ref(false)
const isFilterOpen = ref(false)

function openDetails(item: CatalogItem) {
  selectedItem.value = item
  isDetailOpen.value = true
}

useHead({
  title: location.value ? `Заказ из ${location.value.name}` : "Foodstock | Вендинг",
  meta: [
    { name: "description", content: "Свежая еда из наших вендинговых точек." },
  ],
})
</script>

<template lang="pug">
div(class="flex flex-col gap-10")
  section(v-if="location" class="flex flex-col gap-2")
    h1(class="display-lg font-extrabold") {{ location.name }}
    p(class="body-md opacity-60 flex items-center gap-2")
      u-icon(name="i-heroicons-map-pin" class="w-4 h-4")
      | {{ location.address }}

  section(v-else class="flex flex-col gap-4 py-20 items-center text-center")
    u-icon(name="i-heroicons-exclamation-triangle" class="w-12 h-12 text-primary opacity-20")
    h1(class="headline-lg") Автомат не найден
    p(class="body-md") QR код возможно устарел или неверный. Пожалуйста, попробуйте отсканировать другой код.

  div(v-if="location" class="flex flex-col gap-8")
    catalog-category-tabs(
      v-model="selectedCategory"
      :categories="categories"
    )

    div(class="flex justify-start items-center")
      u-chip(color="primary" :show="isFiltersActive" size="2xs")
        u-button(
          label="Фильтры КБЖУ"
          icon="i-heroicons-adjustments-horizontal"
          variant="ghost"
          class="btn-secondary px-5 py-2.5 rounded-full transition-all duration-300 transform active:scale-95"
          @click="isFilterOpen = true"
        )

    div(v-if="filteredItems.length === 0" class="flex flex-col items-center gap-4 py-20 text-center")
      u-icon(name="i-heroicons-magnifying-glass" class="w-12 h-12 opacity-20")
      p(class="body-lg opacity-60") Ничего не найдено. Попробуйте сбросить фильтры.
      u-button(variant="ghost" color="primary" @click="resetFilters") Сбросить фильтры

    div(class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6")
      catalog-item-card(
        v-for="item in filteredItems"
        :key="item.id"
        :item="item"
        @select="openDetails"
      )

  catalog-details-drawer(
    v-model:open="isDetailOpen"
    :item="selectedItem"
  )

  catalog-filter-drawer(
    v-model:open="isFilterOpen"
    v-model:filters="filters"
    @reset="resetFilters"
  )
</template>
