<script setup lang="ts">
defineProps<{
  current: "dashboard" | "locations" | "items"
}>()

const links = [
  { key: "dashboard", label: "Дашборд", to: "/admin", icon: "i-heroicons-squares-2x2" },
  { key: "locations", label: "Локации", to: "/admin/locations", icon: "i-heroicons-map-pin" },
  { key: "items", label: "Товары", to: "/admin/items", icon: "i-heroicons-cube" },
] as const
</script>

<template lang="pug">
nav(class="surface-section rounded-3xl p-2 flex flex-wrap gap-3")
  nuxt-link(
    v-for="link in links"
    :key="link.key"
    :to="link.to"
    class="px-6 py-2.5 rounded-full inline-flex items-center gap-2.5 text-sm font-bold transition-all duration-300 transform active:scale-95 shadow-sm"
    :class="current === link.key ? 'btn-primary' : 'btn-secondary opacity-80 hover:opacity-100'"
  )
    u-icon(:name="link.icon" class="w-4.5 h-4.5")
    span {{ link.label }}
</template>
