<script setup lang="ts">
import type { Filters } from "~/types/catalog"

const emit = defineEmits<{
  reset: []
}>()

const filters = defineModel<Filters>("filters", { required: true })

const open = defineModel<boolean>("open", { required: true })
</script>

<template lang="pug">
u-drawer(
  v-model:open="open"
  direction="bottom"
  :ui="{ content: 'max-h-[90dvh] flex flex-col rounded-t-3xl overflow-hidden' }"
)
  template(#header)
    div(class="flex justify-between items-center w-full px-2 pt-2")
      u-button(
        variant="ghost"
        color="primary"
        class="rounded-full"
        @click="emit('reset'); open = false"
      ) Сбросить
      h2(class="headline-sm font-bold") Фильтры КБЖУ
      u-button(
        icon="i-heroicons-x-mark"
        variant="ghost"
        color="primary"
        class="rounded-full"
        @click="open = false"
      )

  template(#body)
    div(class="flex flex-col gap-10 px-2 pb-12 pt-4")
      // Калории
      div(class="flex flex-col gap-4")
        div(class="flex justify-between items-center")
          span(class="headline-sm font-bold") Калории
          span(class="body-md font-medium text-primary") {{ filters.calories[0] }} - {{ filters.calories[1] }} ккал
        u-slider(
          v-model="filters.calories"
          :min="0"
          :max="1000"
          :step="10"
        )

      // Белки
      div(class="flex flex-col gap-4")
        div(class="flex justify-between items-center")
          span(class="headline-sm font-bold") Белки
          span(class="body-md font-medium text-primary") {{ filters.proteins[0] }} - {{ filters.proteins[1] }} г
        u-slider(
          v-model="filters.proteins"
          :min="0"
          :max="100"
          :step="1"
        )

      // Жиры
      div(class="flex flex-col gap-4")
        div(class="flex justify-between items-center")
          span(class="headline-sm font-bold") Жиры
          span(class="body-md font-medium text-primary") {{ filters.fats[0] }} - {{ filters.fats[1] }} г
        u-slider(
          v-model="filters.fats"
          :min="0"
          :max="100"
          :step="1"
        )

      // Углеводы
      div(class="flex flex-col gap-4")
        div(class="flex justify-between items-center")
          span(class="headline-sm font-bold") Углеводы
          span(class="body-md font-medium text-primary") {{ filters.carbs[0] }} - {{ filters.carbs[1] }} г
        u-slider(
          v-model="filters.carbs"
          :min="0"
          :max="200"
          :step="1"
        )

  template(#footer)
    div(class="p-4 w-full")
      u-button(
        block
        size="xl"
        class="btn-primary py-4"
        @click="open = false"
      ) Применить
</template>
