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

const categoryItems = computed(() => [
  { label: "Все категории", value: "all" },
  ...view.categories.map(category => ({ label: category, value: category })),
])

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
div(class="flex flex-col gap-10")
  section(class="flex flex-col gap-3")
    h1(class="display-lg font-extrabold") Товары
    p(class="body-md opacity-60 flex items-center gap-2")
      u-icon(name="i-heroicons-cube" class="w-5 h-5 text-primary")
      | Список карточек, цен и остатков по всем точкам.

  admin-nav(current="items")

  section(class="flex flex-col gap-6")
    div(class="flex flex-col lg:flex-row gap-4 items-center justify-between bg-white/50 dark:bg-gray-900/50 p-2 rounded-[2.5rem] shadow-sm border border-gray-100 dark:border-gray-800")
      div(class="flex flex-col md:flex-row gap-3 w-full lg:max-w-2xl")
        u-input(
          v-model="itemSearch"
          placeholder="Поиск товара..."
          size="xl"
          icon="i-heroicons-magnifying-glass"
          variant="none"
          class="flex-1"
        )

        div(class="h-8 w-px bg-gray-100 dark:bg-gray-800 hidden md:block self-center")

        u-select(
          v-model="itemCategory"
          size="xl"
          variant="none"
          class="w-full md:w-64"
          :items="categoryItems"
        )

      u-button(
        to="/admin/items/new"
        size="xl"
        class="btn-primary w-full lg:w-auto px-8 py-3 transform transition-all active:scale-95 shadow-md"
      )
        template(#leading)
          u-icon(name="i-heroicons-plus-circle" class="w-5 h-5")
        | Добавить товар

  section(v-if="view.items.length === 0" class="surface-card container-pad text-center py-20 flex flex-col items-center gap-4")
    u-icon(name="i-heroicons-cube-transparent" class="w-16 h-16 opacity-10")
    h2(class="headline-md font-bold") Товары не найдены
    p(class="body-md opacity-50") Попробуйте изменить фильтры или создать новый товар.

  section(v-else class="grid grid-cols-1 xl:grid-cols-2 gap-8")
    article(
      v-for="item in view.items"
      :key="item.id"
      class="surface-card container-pad flex flex-col gap-6 transition-all duration-300 hover:shadow-lg relative overflow-hidden"
    )
      div(class="flex flex-col md:flex-row md:items-start gap-6")
        div(class="w-24 h-24 rounded-2xl bg-gray-50 dark:bg-gray-900 flex-shrink-0 flex items-center justify-center border border-gray-100 dark:border-gray-800")
          u-icon(name="i-heroicons-photo" class="w-8 h-8 opacity-20")

        div(class="flex-1 flex flex-col gap-2")
          div(class="flex items-start justify-between gap-4")
            h3(class="headline-md font-extrabold") {{ item.name }}
            div(class="text-2xl font-black text-primary whitespace-nowrap") {{ formatNumber(item.price / 100) }} ₽

          p(class="body-md opacity-60 line-clamp-2 min-h-[3rem]") {{ item.description || "Описание отсутствует." }}

          div(class="flex flex-wrap items-center gap-3 mt-2")
            u-badge(color="primary" variant="soft" class="px-3") {{ item.category }}
            u-badge(
              :color="item.stock_amount < 5 ? 'error' : 'neutral'"
              variant="soft"
              class="px-3"
            )
              | Остаток: {{ item.stock_amount }} шт.

      div(class="flex items-center gap-3 border-t border-gray-100 dark:border-gray-800 pt-6 mt-2")
        u-button(
          :to="`/admin/items/${item.id}`"
          variant="ghost"
          class="btn-secondary flex-1 py-3 rounded-full font-bold transform transition-all active:scale-95"
        )
          | Редактировать

        u-button(
          variant="ghost"
          color="error"
          class="hover:bg-red-50 dark:hover:bg-red-900/10 px-6 py-3 rounded-full transition-all duration-300 transform active:scale-95"
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
