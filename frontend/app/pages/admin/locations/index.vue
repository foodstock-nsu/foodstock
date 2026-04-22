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
div(class="flex flex-col gap-8")
  section(class="flex flex-col gap-2")
    h1(class="headline-lg font-bold") Локации
    p(class="body-md opacity-70") Управление адресами и доступностью точек.

  admin-nav(current="locations")

  section(class="rounded-3xl container-pad flex flex-col gap-4")
    div(class="grid grid-cols-1 md:grid-cols-[1fr_auto] gap-3")
      u-input(
        v-model="locationSearch"
        icon="i-heroicons-magnifying-glass"
        size="xl"
        class="w-full"
        placeholder="Поиск по имени, slug или адресу"
      )

      u-button(to="/admin/locations/new" size="xl" class="justify-center") Добавить локацию

  section(v-if="view.locations.length === 0" class="surface-card container-pad text-center flex flex-col gap-2")
    h2(class="headline-md") Локации не найдены
    p(class="body-md opacity-70") Измените фильтр поиска или создайте новую локацию.

  section(v-if="error" class="surface-card container-pad")
    p(class="text-sm font-semibold text-red-600") {{ error }}

  section(v-else class="list-soft")
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
      div(class="surface-card w-full max-w-md p-6 flex flex-col gap-4")
        div(class="flex items-start justify-between gap-3")
          h2(class="headline-md") QR-код: {{ qrPreview.title }}
          u-button(type="button" variant="soft" color="neutral" @click="qrPreview.isOpen = false") Закрыть

        img(
          v-if="qrPreview.imageUrl"
          :src="qrPreview.imageUrl"
          :alt="`QR-код локации ${qrPreview.title}`"
          class="w-full max-w-72 self-center"
        )

        p(class="text-xs opacity-70 text-center") Отсканируйте код, чтобы открыть страницу локации.
</template>
