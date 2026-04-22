<script setup lang="ts">
import type { AdminLocationForm } from "~/types/admin"

const route = useRoute()
const { createLocationForm, getLocationById, saveLocation, isLoading } = useAdmin()

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

  await saveLocation(payload, isNew.value ? undefined : locationId.value)
  await navigateTo("/admin/locations")
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
div(class="flex flex-col gap-8")
  section(class="flex flex-col gap-2")
    h1(class="headline-lg font-bold") {{ isNew ? "Новая локация" : "Редактирование локации" }}
    p(class="body-md opacity-70") Базовые данные точки и ее доступность в клиентском каталоге.

  admin-nav(current="locations")

  section(v-if="notFound" class="surface-card container-pad flex flex-col gap-2 text-center")
    h2(class="headline-md") Локация не найдена
    p(class="body-md opacity-70") Проверьте ссылку или вернитесь к списку локаций.
    nuxt-link(to="/admin/locations" class="btn-secondary px-5 py-2.5 self-center") К списку локаций

  form(v-else class="surface-card container-pad flex flex-col gap-6" @submit.prevent="onSubmit")
    div(class="grid grid-cols-1 md:grid-cols-2 gap-4")
      label(class="flex flex-col gap-2")
        span(class="text-sm font-semibold") Slug
        input(v-model="form.slug" class="input-minimal px-4 py-3" required)

      label(class="flex flex-col gap-2")
        span(class="text-sm font-semibold") Название
        input(v-model="form.name" class="input-minimal px-4 py-3" required)

      label(class="md:col-span-2 flex flex-col gap-2")
        span(class="text-sm font-semibold") Адрес
        textarea(v-model="form.address" rows="3" class="input-minimal px-4 py-3 resize-y" required)

    label(class="surface-section rounded-3xl p-4 inline-flex items-center gap-3 w-fit")
      input(v-model="form.is_active" type="checkbox" class="h-5 w-5")
      span(class="font-semibold") Локация активна

    div(class="flex flex-wrap items-center gap-3")
      button(type="submit" class="btn-primary px-6 py-3" :disabled="isLoading")
        span(v-if="isLoading") Сохраняем...
        span(v-else) Сохранить

      nuxt-link(to="/admin/locations" class="btn-secondary px-6 py-3") Отмена
</template>
