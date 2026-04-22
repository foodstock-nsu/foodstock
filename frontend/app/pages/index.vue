<script setup lang="ts">
// Перенаправляем на демонстрационную точку продаж
definePageMeta({
  middleware: [
    function () {
      return navigateTo("/l/central-ave")
    },
  ],
})
</script>

<template>
  <div class="flex items-center justify-center min-h-screen">
    <p class="text-sm opacity-50">Перенаправление...</p>
  </div>
</template>
