<script setup lang="ts">
defineProps<{
  nutrition: Nutrition
}>()
</script>

<template lang="pug">
div(class="grid grid-cols-4 gap-2")
  div(class="nutrition-card")
    span(class="nutrition-value") {{ nutrition.calories }}
    span(class="nutrition-label") Ккал

  div(class="nutrition-card")
    span(class="nutrition-value") {{ nutrition.proteins }}г
    span(class="nutrition-label") Белки

  div(class="nutrition-card")
    span(class="nutrition-value") {{ nutrition.fats }}г
    span(class="nutrition-label") Жиры

  div(class="nutrition-card")
    span(class="nutrition-value") {{ nutrition.carbs }}г
    span(class="nutrition-label") Углеводы
</template>

<style scoped>
.nutrition-card {
  background: var(--color-surface-container);
  border-radius: 16px;
  padding: 12px 8px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
}

.nutrition-value {
  font-size: 1.125rem;
  font-weight: 800;
  color: var(--color-on-surface);
}

.nutrition-label {
  font-size: 0.75rem;
  text-transform: lowercase;
  color: var(--color-on-surface);
  opacity: 0.6;
}
</style>
