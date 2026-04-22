<script setup lang="ts">
const { stats } = useAdmin()

const cards = computed(() => {
  return [
    {
      title: "Локации",
      description: "Управляйте адресами, slug и статусом активности.",
      value: stats.value.locationsCount,
      helper: `Активных: ${stats.value.activeLocationsCount}`,
      to: "/admin/locations",
      icon: "i-heroicons-map-pin",
    },
    {
      title: "Товары",
      description: "Редактируйте карточки товаров, цену и остатки.",
      value: stats.value.itemsCount,
      helper: `Низкий остаток: ${stats.value.lowStockCount}`,
      to: "/admin/items",
      icon: "i-heroicons-cube",
    },
  ]
})

useHead({
  title: "Foodstock | Админка",
})
</script>

<template lang="pug">
div(class="flex flex-col gap-8")
  section(class="flex flex-col gap-2")
    h1(class="display-lg font-extrabold") Админка
    p(class="body-md opacity-70") Управление контентом и витриной с локальным mock API.

  admin-nav(current="dashboard")

  section(class="grid grid-cols-1 lg:grid-cols-2 gap-6")
    nuxt-link(
      v-for="card in cards"
      :key="card.title"
      :to="card.to"
      class="surface-card container-pad flex flex-col gap-4 transition-all duration-300 hover:-translate-y-1"
    )
      div(class="flex items-center justify-between")
        h2(class="headline-md font-bold") {{ card.title }}
        u-icon(:name="card.icon" class="w-6 h-6 text-primary")

      p(class="body-md opacity-70") {{ card.description }}

      div(class="mt-4 flex items-end justify-between")
        div(class="text-4xl font-extrabold text-primary") {{ card.value }}
        p(class="text-sm opacity-65") {{ card.helper }}
</template>
