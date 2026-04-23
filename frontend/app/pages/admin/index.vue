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
    p(class="body-md opacity-70") Управление контентом и витриной через backend API.

  admin-nav(current="dashboard")

  section(class="grid grid-cols-1 lg:grid-cols-2 gap-8")
    nuxt-link(
      v-for="card in cards"
      :key="card.title"
      :to="card.to"
      class="surface-card container-pad flex flex-col gap-6 transition-all duration-300 hover:-translate-y-1 hover:shadow-xl group active:scale-[0.98]"
    )
      div(class="flex items-center justify-between")
        h2(class="headline-md font-extrabold") {{ card.title }}
        div(class="w-12 h-12 rounded-2xl bg-primary-50 dark:bg-primary-900/10 flex items-center justify-center transition-colors group-hover:bg-primary group-hover:text-white")
          u-icon(:name="card.icon" class="w-6 h-6 text-primary group-hover:text-white transition-colors")

      p(class="body-md opacity-60") {{ card.description }}

      div(class="mt-4 flex items-end justify-between border-t border-gray-100 dark:border-gray-800 pt-6")
        div(class="flex flex-col")
          span(class="text-xs uppercase tracking-wider font-bold opacity-40") Всего
          div(class="text-5xl font-black text-primary") {{ card.value }}
        div(class="text-right")
          p(class="text-sm font-bold opacity-60") {{ card.helper }}
          p(class="text-xs opacity-40 mt-1") Нажмите для управления &rarr;
</template>
