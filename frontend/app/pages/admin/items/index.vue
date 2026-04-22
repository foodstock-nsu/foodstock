<script setup lang="ts">
import type { AdminCategory, AdminItem } from "~/types/admin"
import { formatNumber } from "~/utils/intl-format"

const { store, removeItem, isLoading } = useAdmin()

const deletingId = ref<string | null>(null)
const itemSearch = ref(store.itemSearch.value)
const itemCategory = ref<AdminCategory | "all">(store.itemCategory.value)
const view = reactive<{ categories: AdminCategory[], items: AdminItem[] }>({
  categories: [],
  items: [],
})

watch(itemSearch, value => store.itemSearch.value = value)
watch(itemCategory, value => store.itemCategory.value = value)

watchEffect(() => {
  view.categories = [...store.categories.value]
  view.items = [...store.sortedItems.value]
})

async function onDeleteItem(id: string) {
  deletingId.value = id
  try {
    await removeItem(id)
  } finally {
    deletingId.value = null
  }
}

useHead({
  title: "Foodstock | Админка товаров",
})
</script>

<template lang="pug">
div(class="flex flex-col gap-8")
  section(class="flex flex-col gap-2")
    h1(class="headline-lg font-bold") Товары
    p(class="body-md opacity-70") Список карточек, цен и остатков.

  admin-nav(current="items")

  section(class="surface-section rounded-3xl container-pad flex flex-col gap-4")
    div(class="grid grid-cols-1 md:grid-cols-[1fr_220px_auto] gap-3 items-center")
      u-input(
        v-model="itemSearch"
        placeholder="Поиск по названию или описанию"
        size="xl"
        icon="i-heroicons-magnifying-glass"
      )

      u-select(
        v-model="itemCategory"
        placeholder="Все категории"
        size="xl"
      )
        option(value="all") Все категории
        option(v-for="category in view.categories" :key="category" :value="category") {{ category }}

      u-button(to="/admin/items/new" size="xl") Добавить товар

  section(v-if="view.items.length === 0" class="surface-card container-pad text-center flex flex-col gap-2")
    h2(class="headline-md") Ничего не найдено
    p(class="body-md opacity-70") Уточните фильтры или создайте новую карточку товара.

  section(v-else class="list-soft")
    article(
      v-for="item in view.items"
      :key="item.id"
      class="surface-card container-pad flex flex-col gap-4"
    )
      div(class="flex flex-col md:flex-row md:items-start md:justify-between gap-3")
        div(class="flex flex-col gap-1")
          h3(class="headline-md font-bold") {{ item.name }}
          p(class="body-md opacity-70 line-clamp-2") {{ item.description || "Без описания" }}

        div(class="flex items-center gap-2 text-xs uppercase tracking-wider")
          u-badge(variant="soft" color="neutral") {{ item.category }}
          u-badge(variant="soft" color="neutral") Остаток: {{ item.stock_amount }}

      div(class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3")
        div(class="text-xl font-bold text-primary") {{ formatNumber(item.price / 100) }} ₽

        div(class="flex items-center gap-2")
          u-button(:to="`/admin/items/${item.id}`" variant="soft" color="neutral")
            | Редактировать

          u-button(
            variant="outline"
            color="primary"
            :loading="deletingId === item.id"
            :disabled="isLoading"
            @click="onDeleteItem(item.id)"
          )
            | Удалить
</template>

<style scoped>
.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>
