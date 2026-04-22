<script setup lang="ts">
import type { AdminLocationForm } from "~/types/admin"

const route = useRoute()
const { createLocationForm, getLocationById, saveLocation, isLoading, error } = useAdmin()

const locationId = computed(() => String(route.params.id || ""))
const isNew = computed(() => locationId.value === "new")
const currentLocationName = ref("")
const notFound = ref(false)
const form = reactive<AdminLocationForm>(createLocationForm())

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

async function onSubmit() {
  const payload = {
    ...form,
    slug: form.slug.trim(),
    name: form.name.trim(),
    address: form.address.trim(),
  }

  if (!payload.slug || !payload.name || !payload.address) {
    return
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

  form(v-else class="surface-card container-pad flex flex-col gap-10" @submit.prevent="onSubmit")
    div(v-if="error" class="bg-red-50 dark:bg-red-900/10 p-4 rounded-2xl border-l-4 border-red-500")
      p(class="text-sm font-bold text-red-600 flex items-center gap-2")
        u-icon(name="i-heroicons-exclamation-triangle" class="w-4 h-4")
        | {{ error }}

    div(class="grid grid-cols-1 md:grid-cols-2 gap-8")
      label(class="flex flex-col gap-3")
        span(class="text-sm font-black uppercase tracking-widest opacity-40") Slug (URL-путь)
        u-input(v-model="form.slug" size="xl" required placeholder="example-location")

      label(class="flex flex-col gap-3")
        span(class="text-sm font-black uppercase tracking-widest opacity-40") Название
        u-input(v-model="form.name" size="xl" required placeholder="Вендинг на Пушкина")

      label(class="md:col-span-2 flex flex-col gap-3")
        span(class="text-sm font-black uppercase tracking-widest opacity-40") Полный адрес
        u-textarea(v-model="form.address" rows="3" size="xl" class="resize-none" required placeholder="г. Новосибирск, ул. Пушкина, д. 10, этаж 1")

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
