<script setup lang="ts">
import type { AdminLocation } from "~/types/admin"

const { store, removeLocation, getLocationQRCode, isLoading, error } = useAdmin()

const deletingId = ref<string | null>(null)
const qrLoadingId = ref<string | null>(null)
const locationSearch = ref(store.locationSearch.value)
const view = reactive<{ locations: AdminLocation[] }>({
  locations: [],
})
const qrPreview = reactive({
  isOpen: false,
  imageUrl: "",
  title: "",
})

watch(locationSearch, value => store.locationSearch.value = value)

watch(() => qrPreview.isOpen, (isOpen) => {
  if (isOpen || !qrPreview.imageUrl) {
    return
  }

  URL.revokeObjectURL(qrPreview.imageUrl)
  qrPreview.imageUrl = ""
  qrPreview.title = ""
})

watchEffect(() => {
  view.locations = [...store.sortedLocations.value]
})

async function onDeleteLocation(id: string) {
  deletingId.value = id
  try {
    await removeLocation(id)
  } finally {
    deletingId.value = null
  }
}

async function onViewQrCode(location: AdminLocation) {
  qrLoadingId.value = location.id
  try {
    if (qrPreview.imageUrl) {
      URL.revokeObjectURL(qrPreview.imageUrl)
      qrPreview.imageUrl = ""
    }

    const imageUrl = await getLocationQRCode(location.id)
    qrPreview.imageUrl = imageUrl
    qrPreview.title = location.name
    qrPreview.isOpen = true
  } finally {
    qrLoadingId.value = null
  }
}

onBeforeUnmount(() => {
  if (qrPreview.imageUrl) {
    URL.revokeObjectURL(qrPreview.imageUrl)
  }
})

useHead({
  title: "Foodstock | Админка локаций",
})
</script>

<template lang="pug">
div(class="flex flex-col gap-10")
  section(class="flex flex-col gap-3")
    h1(class="display-lg font-extrabold") Локации
    p(class="body-md opacity-60 flex items-center gap-2")
      u-icon(name="i-heroicons-map-pin" class="w-5 h-5 text-primary")
      | Управление адресами и доступностью точек.

  admin-nav(current="locations")

  section(class="flex flex-col gap-6")
    div(class="flex flex-col md:flex-row gap-4 items-center justify-between bg-white/50 dark:bg-gray-900/50 p-2 rounded-[2.5rem] shadow-sm border border-gray-100 dark:border-gray-800")
      div(class="relative w-full md:max-w-md")
        u-input(
          v-model="locationSearch"
          icon="i-heroicons-magnifying-glass"
          size="xl"
          class="w-full"
          placeholder="Поиск по названию или адресу..."
          variant="none"
        )

      u-button(
        to="/admin/locations/new"
        size="xl"
        class="btn-primary w-full md:w-auto px-8 py-3 transform transition-all active:scale-95 shadow-md"
      )
        template(#leading)
          u-icon(name="i-heroicons-plus-circle" class="w-5 h-5")
        | Добавить локацию

  section(v-if="view.locations.length === 0" class="surface-card container-pad text-center py-20 flex flex-col items-center gap-4")
    u-icon(name="i-heroicons-map" class="w-16 h-16 opacity-10")
    h2(class="headline-md font-bold") Локации не найдены
    p(class="body-md opacity-50") Попробуйте изменить параметры поиска.

  section(v-if="error" class="surface-card container-pad border-l-4 border-red-500 bg-red-50 dark:bg-red-900/10")
    div(class="flex items-center gap-3 text-red-600 font-bold")
      u-icon(name="i-heroicons-exclamation-circle" class="w-5 h-5")
      p {{ error }}

  section(v-else class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-8")
    admin-location-card(
      v-for="location in view.locations"
      :key="location.id"
      :location="location"
      :is-busy="isLoading"
      :is-deleting="deletingId === location.id"
      :is-qr-loading="qrLoadingId === location.id"
      @delete="onDeleteLocation"
      @view-qr="onViewQrCode"
    )

  u-modal(v-model:open="qrPreview.isOpen")
    template(#content)
      div(class="surface-card w-full max-w-md p-8 md:p-10 flex flex-col items-center gap-6 text-center")
        div(class="w-full flex items-center justify-between mb-2")
          div(class="text-left")
            h2(class="headline-md font-extrabold") QR-код
            p(class="text-sm opacity-60") {{ qrPreview.title }}
          u-button(
            variant="ghost"
            color="neutral"
            class="rounded-full h-10 w-10 !p-0"
            @click="qrPreview.isOpen = false"
          )
            u-icon(name="i-heroicons-x-mark" class="w-6 h-6")

        div(class="bg-white p-6 rounded-3xl shadow-inner border border-gray-50")
          img(
            v-if="qrPreview.imageUrl"
            :src="qrPreview.imageUrl"
            :alt="`QR-код локации ${qrPreview.title}`"
            class="w-64 h-64"
          )

        p(class="body-md opacity-60 px-4") Отсканируйте этот код для быстрого перехода на витрину этой точки.

        u-button(
          variant="ghost"
          color="primary"
          class="btn-secondary w-full py-4 rounded-full font-bold"
          @click="qrPreview.isOpen = false"
        )
          | Понятно
</template>
