<script setup lang="ts">
import type { AdminItemForm } from "~/types/admin"

const route = useRoute()
const { createItemForm, getItemById, saveItem, isLoading } = useAdmin()

const itemId = computed(() => String(route.params.id || ""))
const isNew = computed(() => itemId.value === "new")
const currentItemName = ref("")
const notFound = ref(false)
const form = reactive<AdminItemForm>(createItemForm())

async function loadItem() {
  if (isNew.value) {
    notFound.value = false
    Object.assign(form, createItemForm())
    return
  }

  const item = await getItemById(itemId.value)
  if (!item) {
    notFound.value = true
    return
  }

  notFound.value = false
  currentItemName.value = item.name
  Object.assign(form, createItemForm(item))
}

async function onSubmit() {
  const payload = {
    ...form,
    name: form.name.trim(),
    description: form.description.trim(),
    photo_url: form.photo_url.trim(),
  }

  if (!payload.name) {
    return
  }

  await saveItem(payload, isNew.value ? undefined : itemId.value)
  await navigateTo("/admin/items")
}

await loadItem()

useHead({
  title: computed(() => {
    if (isNew.value) {
      return "Foodstock | Новый товар"
    }
    return currentItemName.value
      ? `Foodstock | ${currentItemName.value}`
      : "Foodstock | Редактирование товара"
  }),
})
</script>

<template lang="pug">
div(class="flex flex-col gap-10")
  section(class="flex flex-col gap-3")
    h1(class="display-lg font-extrabold") {{ isNew ? "Новый товар" : "Редактирование" }}
    p(class="body-md opacity-60 flex items-center gap-2")
      u-icon(name="i-heroicons-cube" class="w-5 h-5 text-primary")
      | {{ isNew ? "Создание новой карточки товара для каталога." : `Редактирование характеристик ${currentItemName}` }}

  admin-nav(current="items")

  section(v-if="notFound" class="surface-card container-pad flex flex-col items-center gap-6 text-center py-20")
    u-icon(name="i-heroicons-face-frown" class="w-16 h-16 opacity-10")
    div
      h2(class="headline-md font-bold") Товар не найден
      p(class="body-md opacity-50 mt-2") Проверьте правильность ссылки или вернитесь к списку.
    u-button(to="/admin/items" variant="ghost" class="btn-secondary px-8 py-3 rounded-full font-bold transform transition-all active:scale-95")
      | Вернуться к списку

  form(v-else class="surface-card container-pad flex flex-col gap-10" @submit.prevent="onSubmit")
    div(class="grid grid-cols-1 md:grid-cols-2 gap-8")
      label(class="flex flex-col gap-3")
        span(class="text-sm font-black uppercase tracking-widest opacity-40") Название товара
        u-input(v-model="form.name" size="xl" required placeholder="Например: Салат Цезарь")

      label(class="flex flex-col gap-3")
        span(class="text-sm font-black uppercase tracking-widest opacity-40") Категория
        u-select(v-model="form.category" size="xl")
          option(value="lunch") Lunch
          option(value="breakfast") Breakfast
          option(value="drinks") Drinks
          option(value="snacks") Snacks
          option(value="desserts") Desserts

      label(class="md:col-span-2 flex flex-col gap-3")
        span(class="text-sm font-black uppercase tracking-widest opacity-40") Описание
        u-textarea(v-model="form.description" rows="3" size="xl" class="resize-none" placeholder="Краткое описание для карточки товара...")

      label(class="md:col-span-2 flex flex-col gap-3")
        span(class="text-sm font-black uppercase tracking-widest opacity-40") Ссылка на изображение (URL)
        u-input(v-model="form.photo_url" size="xl" placeholder="https://example.com/photo.jpg")

      label(class="flex flex-col gap-3")
        span(class="text-sm font-black uppercase tracking-widest opacity-40") Цена (в рублях)
        u-input(v-model.number="form.priceRub" type="number" min="0" size="xl" icon="i-heroicons-banknotes")

      label(class="flex flex-col gap-3")
        span(class="text-sm font-black uppercase tracking-widest opacity-40") Доступный остаток
        u-input(v-model.number="form.stock_amount" type="number" min="0" size="xl" icon="i-heroicons-circle-stack")

    div(class="bg-gray-50 dark:bg-gray-900/50 rounded-[2.5rem] p-8 md:p-10 flex flex-col gap-8 border border-gray-100 dark:border-gray-800")
      div(class="flex items-center gap-3")
        div(class="w-10 h-10 rounded-xl bg-primary/10 flex items-center justify-center")
          u-icon(name="i-heroicons-sparkles" class="w-5 h-5 text-primary")
        h3(class="headline-md font-extrabold") Пищевая ценность (КБЖУ)

      div(class="grid grid-cols-2 lg:grid-cols-4 gap-6")
        label(class="flex flex-col gap-3")
          span(class="text-[10px] font-black uppercase tracking-[0.2em] opacity-40") Калории
          u-input(v-model.number="form.nutrition.calories" type="number" min="0" size="xl")

        label(class="flex flex-col gap-3")
          span(class="text-[10px] font-black uppercase tracking-[0.2em] opacity-40") Белки (г)
          u-input(v-model.number="form.nutrition.proteins" type="number" min="0" step="0.1" size="xl")

        label(class="flex flex-col gap-3")
          span(class="text-[10px] font-black uppercase tracking-[0.2em] opacity-40") Жиры (г)
          u-input(v-model.number="form.nutrition.fats" type="number" min="0" step="0.1" size="xl")

        label(class="flex flex-col gap-3")
          span(class="text-[10px] font-black uppercase tracking-[0.2em] opacity-40") Углеводы (г)
          u-input(v-model.number="form.nutrition.carbs" type="number" min="0" step="0.1" size="xl")

    div(class="flex flex-wrap items-center gap-4 pt-4 border-t border-gray-100 dark:border-gray-800")
      u-button(
        type="submit"
        size="xl"
        class="btn-primary px-10 py-4 transform transition-all active:scale-95 shadow-lg"
        :loading="isLoading"
      )
        | {{ isNew ? "Создать товар" : "Сохранить изменения" }}

      u-button(
        to="/admin/items"
        variant="ghost"
        class="btn-secondary px-8 py-4 rounded-full transition-all duration-300 transform active:scale-95"
      )
        | Отмена
</template>
