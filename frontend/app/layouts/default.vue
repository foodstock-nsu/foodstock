<script setup lang="ts">
import { useRoute } from "vue-router"

const { totalQuantity } = useCartStore()
const { logout } = useAdminAuth()

const route = useRoute()
const isLocationPage = computed(() => route.path.startsWith("/l/"))
</script>

<template lang="pug">
div(class="min-h-screen surface-base")
  header(class="sticky top-0 z-50 glass-nav container-pad py-4 flex items-center justify-between shadow-soft")
    div(class="headline-md font-bold text-primary" @click="logout") Foodstock

  main(class="container-pad")
    slot

    div(
      v-if="isLocationPage"
      class="transition-all duration-300 ease-in-out"
      :class="totalQuantity > 0 ? 'pb-20' : ''"
    )

  cart-bar(v-if="isLocationPage")
</template>

<style scoped>
header {
  border-bottom: 1px solid var(--ghost-border);
}
</style>
