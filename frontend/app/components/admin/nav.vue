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
nav(class="surface-section rounded-3xl p-2 flex flex-wrap gap-2")
  nuxt-link(
    v-for="link in links"
    :key="link.key"
    :to="link.to"
    class="px-4 py-2.5 rounded-full inline-flex items-center gap-2 text-sm font-semibold transition-all duration-200"
    :class="current === link.key ? 'btn-primary' : 'btn-secondary opacity-85 hover:opacity-100'"
  )
    u-icon(:name="link.icon" class="w-4 h-4")
    span {{ link.label }}
</template>
