<script setup lang="ts">
import type { AdminLocation } from "~/types/admin"

const { store, removeLocation, isLoading } = useAdmin()

const deletingId = ref<string | null>(null)
const locationSearch = ref(store.locationSearch.value)
const view = reactive<{ locations: AdminLocation[] }>({
  locations: [],
})

watch(locationSearch, value => store.locationSearch.value = value)

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

  section(class="surface-section rounded-3xl container-pad flex flex-col gap-4")
    div(class="grid grid-cols-1 md:grid-cols-[1fr_auto] gap-3")
      input(
        v-model="locationSearch"
        class="input-minimal px-4 py-3 w-full"
        placeholder="Поиск по имени, slug или адресу"
      )

      nuxt-link(to="/admin/locations/new" class="btn-primary px-6 py-3 text-center") Добавить локацию

  section(v-if="view.locations.length === 0" class="surface-card container-pad text-center flex flex-col gap-2")
    h2(class="headline-md") Локации не найдены
    p(class="body-md opacity-70") Измените фильтр поиска или создайте новую локацию.

  section(v-else class="list-soft")
    article(
      v-for="location in view.locations"
      :key="location.id"
      class="surface-card container-pad flex flex-col gap-4"
    )
      div(class="flex flex-col md:flex-row md:items-start md:justify-between gap-3")
        div(class="flex flex-col gap-1")
          h3(class="headline-md font-bold") {{ location.name }}
          p(class="body-md opacity-70") {{ location.address }}
          p(class="text-sm opacity-55") /l/{{ location.slug }}

        span(
          class="px-3 py-1 rounded-full text-xs font-semibold uppercase tracking-wider"
          :class="location.is_active !== false ? 'status-ready' : 'btn-secondary'"
        )
          | {{ location.is_active !== false ? "Активна" : "Отключена" }}

      div(class="flex items-center gap-2")
        nuxt-link(:to="`/admin/locations/${location.id}`" class="btn-secondary px-4 py-2") Редактировать
        button(
          class="btn-tertiary px-4 py-2"
          :disabled="deletingId === location.id || isLoading"
          @click="onDeleteLocation(location.id)"
        )
          span(v-if="deletingId === location.id") Удаляем...
          span(v-else) Удалить
</template>
