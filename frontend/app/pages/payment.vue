<script setup lang="ts">
import { useRouter } from "vue-router"

const router = useRouter()
const isProcessing = ref(true)

onMounted(() => {
  // Simulate payment processing time
  setTimeout(() => {
    isProcessing.value = false
  }, 2000)
})

function goBack() {
  // Ideally, this should redirect to the previous location or a success page, but we'll go to home for now.
  router.push("/")
}
</script>

<template lang="pug">
div(class="min-h-[100dvh] flex flex-col items-center justify-center p-6 bg-surface text-on-surface")
  div(class="max-w-sm w-full flex flex-col items-center gap-8 text-center")

    // Mock SBP Logo or Icon
    div(class="w-24 h-24 rounded-full bg-primary/10 flex items-center justify-center")
      u-icon(name="i-heroicons-qr-code" class="w-12 h-12 text-primary")

    div(class="flex flex-col gap-2")
      h1(class="headline-lg font-bold") Оплата СБП
      p(class="body-md opacity-70") Тестовая страница оплаты для демонстрации работы

    // Loading State
    div(v-if="isProcessing" class="flex flex-col items-center gap-4 py-8")
      u-icon(name="i-heroicons-arrow-path" class="w-8 h-8 text-primary animate-spin")
      p(class="body-md font-medium animate-pulse") Ожидание оплаты...

    // Success State
    div(v-else class="flex flex-col items-center gap-6 py-8 w-full")
      div(class="w-16 h-16 rounded-full bg-green-100 flex items-center justify-center text-green-600 mb-2")
        u-icon(name="i-heroicons-check" class="w-8 h-8")
      h2(class="headline-md text-green-600 font-bold") Оплата прошла успешно!

      u-button(
        size="lg"
        block
        color="primary"
        class="mt-4 py-4 text-base font-bold rounded-full"
        @click="goBack"
      ) Вернуться на главную
</template>

<style scoped>
/* Scoped styles if needed */
</style>
