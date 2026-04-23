<script setup lang="ts">
import type { FormSubmitEvent } from "@nuxt/ui"
import * as v from "valibot"

import type { AdminLocationForm } from "~/types/admin"

const route = useRoute()
const { createLocationForm, getLocationById, saveLocation, isLoading, error } = useAdmin()

const locationId = computed(() => String(route.params.id || ""))
const isNew = computed(() => locationId.value === "new")
const currentLocationName = ref("")
const notFound = ref(false)
const form = reactive<AdminLocationForm>(createLocationForm())

const schema = v.object({
  slug: v.pipe(
    v.string(),
    v.nonEmpty("Slug обязателен"),
    v.minLength(4, "Минимум 4 символа"),
    v.maxLength(10, "Максимум 10 символов"),
    v.regex(/^[\w-]+$/, "Только латиница, цифры, дефис и подчеркивание"),
  ),
  name: v.pipe(
    v.string(),
    v.nonEmpty("Название обязательно"),
    v.minLength(4, "Минимум 4 символа"),
  ),
  address: v.pipe(
    v.string(),
    v.nonEmpty("Адрес обязателен"),
    v.minLength(20, "Минимум 20 символов"),
  ),
  is_active: v.boolean(),
})

type Schema = v.InferOutput<typeof schema>

async function loadLocation() {
  if (isNew.value) {
    notFound.value = false
    Object.assign(form, createLocationForm())
    return
  }

  const location = await getLocationById(locationId.value)
  if (!location) {
    notFound.value = true
    return
  }

  notFound.value = false
  currentLocationName.value = location.name
  Object.assign(form, createLocationForm(location))
}

async function onSubmit(event: FormSubmitEvent<Schema>) {
  const payload = {
    ...event.data,
    slug: event.data.slug.trim(),
    name: event.data.name.trim(),
    address: event.data.address.trim(),
  }

  try {
    await saveLocation(payload, isNew.value ? undefined : locationId.value)
    await navigateTo("/admin/locations")
  } catch {
    // Ошибка уже нормализована в useAdmin.error
  }
}

await loadLocation()

useHead({
  title: computed(() => {
    if (isNew.value) {
      return "Foodstock | Новая локация"
    }
    return currentLocationName.value
      ? `Foodstock | ${currentLocationName.value}`
      : "Foodstock | Редактирование локации"
  }),
})
</script>

<template lang="pug">
div(class="flex flex-col gap-10")
  section(class="flex flex-col gap-3")
    h1(class="display-lg font-extrabold") {{ isNew ? "Новая локация" : "Редактирование" }}
    p(class="body-md opacity-60 flex items-center gap-2")
      u-icon(name="i-heroicons-pencil-square" class="w-5 h-5 text-primary")
      | {{ isNew ? "Создание новой точки выдачи заказов." : `Настройка параметров для ${currentLocationName}` }}

  admin-nav(current="locations")

  section(v-if="notFound" class="surface-card container-pad flex flex-col items-center gap-6 text-center py-20")
    u-icon(name="i-heroicons-face-frown" class="w-16 h-16 opacity-10")
    div
      h2(class="headline-md font-bold") Локация не найдена
      p(class="body-md opacity-50 mt-2") Проверьте ссылку или вернитесь к списку локаций.
    u-button(to="/admin/locations" variant="ghost" class="btn-secondary px-8 py-3 rounded-full font-bold transform transition-all active:scale-95")
      | Вернуться к списку

  u-form(v-else :schema="schema" :state="form" class="surface-card container-pad flex flex-col gap-10" @submit="onSubmit")
    div(v-if="error" class="bg-red-50 dark:bg-red-900/10 p-4 rounded-2xl border-l-4 border-red-500")
      p(class="text-sm font-bold text-red-600 flex items-center gap-2")
        u-icon(name="i-heroicons-exclamation-triangle" class="w-4 h-4")
        | {{ error }}

    div(class="grid grid-cols-1 md:grid-cols-2 gap-8")
      u-form-field(name="slug" class="flex flex-col gap-3")
        template(#label)
          span(class="text-sm font-black uppercase tracking-widest opacity-40") Slug (URL-путь)
        u-input(v-model="form.slug" size="xl" placeholder="example-location")

      u-form-field(name="name" class="flex flex-col gap-3")
        template(#label)
          span(class="text-sm font-black uppercase tracking-widest opacity-40") Название
        u-input(v-model="form.name" size="xl" placeholder="Вендинг на Пушкина")

      u-form-field(name="address" class="md:col-span-2 flex flex-col gap-3 w-full")
        template(#label)
          span(class="text-sm font-black uppercase tracking-widest opacity-40") Полный адрес
        u-textarea(v-model="form.address" :rows="3" size="xl" class="resize-none w-full" placeholder="г. Новосибирск, ул. Пушкина, д. 10, этаж 1")

    div(class="flex flex-col md:flex-row items-center gap-8 border-t border-gray-100 dark:border-gray-800 pt-8")
      label(class="bg-gray-50 dark:bg-gray-900/50 rounded-full px-6 py-4 flex items-center gap-4 cursor-pointer transition-all hover:bg-gray-100 dark:hover:bg-gray-900")
        u-checkbox(v-model="form.is_active" size="lg" color="primary")
        span(class="font-bold") Локация активна и видна клиентам

    div(class="flex flex-wrap items-center gap-4")
      u-button(
        type="submit"
        size="xl"
        class="btn-primary px-10 py-4 transform transition-all active:scale-95 shadow-lg"
        :loading="isLoading"
      )
        | {{ isNew ? "Создать локацию" : "Сохранить изменения" }}

      u-button(
        to="/admin/locations"
        variant="ghost"
        class="btn-secondary px-8 py-4 rounded-full transition-all duration-300 transform active:scale-95"
      )
        | Отмена
</template>
