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
div(class="flex flex-col gap-8")
  section(class="flex flex-col gap-2")
    h1(class="headline-lg font-bold") {{ isNew ? "Новый товар" : "Редактирование товара" }}
    p(class="body-md opacity-70") Форма использует заглушки до готовности item-endpoints на backend.

  admin-nav(current="items")

  section(v-if="notFound" class="surface-card container-pad flex flex-col gap-2 text-center")
    h2(class="headline-md") Товар не найден
    p(class="body-md opacity-70") Проверьте ссылку или вернитесь к списку.
    nuxt-link(to="/admin/items" class="btn-secondary px-5 py-2.5 self-center") К списку товаров

  form(v-else class="surface-card container-pad flex flex-col gap-6" @submit.prevent="onSubmit")
    div(class="grid grid-cols-1 md:grid-cols-2 gap-4")
      label(class="flex flex-col gap-2")
        span(class="text-sm font-semibold") Название
        input(v-model="form.name" class="input-minimal px-4 py-3" required)

      label(class="flex flex-col gap-2")
        span(class="text-sm font-semibold") Категория
        select(v-model="form.category" class="input-minimal px-4 py-3")
          option(value="lunch") lunch
          option(value="breakfast") breakfast
          option(value="drinks") drinks
          option(value="snacks") snacks
          option(value="desserts") desserts

      label(class="md:col-span-2 flex flex-col gap-2")
        span(class="text-sm font-semibold") Описание
        textarea(v-model="form.description" rows="3" class="input-minimal px-4 py-3 resize-y")

      label(class="md:col-span-2 flex flex-col gap-2")
        span(class="text-sm font-semibold") Ссылка на фото
        input(v-model="form.photo_url" class="input-minimal px-4 py-3")

      label(class="flex flex-col gap-2")
        span(class="text-sm font-semibold") Цена, ₽
        input(v-model.number="form.priceRub" type="number" min="0" class="input-minimal px-4 py-3")

      label(class="flex flex-col gap-2")
        span(class="text-sm font-semibold") Остаток
        input(v-model.number="form.stock_amount" type="number" min="0" class="input-minimal px-4 py-3")

    div(class="surface-section rounded-3xl p-4 flex flex-col gap-4")
      h3(class="headline-md font-bold") КБЖУ
      div(class="grid grid-cols-2 md:grid-cols-4 gap-3")
        label(class="flex flex-col gap-2")
          span(class="text-xs font-semibold uppercase tracking-wider") Calories
          input(v-model.number="form.nutrition.calories" type="number" min="0" class="input-minimal px-4 py-3")

        label(class="flex flex-col gap-2")
          span(class="text-xs font-semibold uppercase tracking-wider") Proteins
          input(v-model.number="form.nutrition.proteins" type="number" min="0" step="0.1" class="input-minimal px-4 py-3")

        label(class="flex flex-col gap-2")
          span(class="text-xs font-semibold uppercase tracking-wider") Fats
          input(v-model.number="form.nutrition.fats" type="number" min="0" step="0.1" class="input-minimal px-4 py-3")

        label(class="flex flex-col gap-2")
          span(class="text-xs font-semibold uppercase tracking-wider") Carbs
          input(v-model.number="form.nutrition.carbs" type="number" min="0" step="0.1" class="input-minimal px-4 py-3")

    div(class="flex flex-wrap items-center gap-3")
      button(type="submit" class="btn-primary px-6 py-3" :disabled="isLoading")
        span(v-if="isLoading") Сохраняем...
        span(v-else) Сохранить

      nuxt-link(to="/admin/items" class="btn-secondary px-6 py-3") Отмена
</template>
