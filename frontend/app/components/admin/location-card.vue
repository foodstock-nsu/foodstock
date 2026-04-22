<script setup lang="ts">
import type { AdminLocation } from "~/types/admin"

const props = defineProps<{
  location: AdminLocation
  isBusy: boolean
  isDeleting: boolean
  isQrLoading: boolean
}>()

const emit = defineEmits<{
  delete: [string]
  viewQr: [AdminLocation]
}>()

const isActive = computed(() => props.location.is_active !== false)
const statusLabel = computed(() => isActive.value ? "Активна" : "Отключена")
const statusColor = computed(() => isActive.value ? "primary" : "neutral")

function onDelete() {
  emit("delete", props.location.id)
}

function onViewQr() {
  emit("viewQr", props.location)
}
</script>

<template lang="pug">
article(class="surface-card container-pad relative flex flex-col gap-4")
  div(class="absolute right-4 top-4")
    u-button(
      color="neutral"
      size="sm"
      class="h-9 w-9 !p-0 inline-flex items-center justify-center"
      :loading="isQrLoading"
      :disabled="isBusy"
      :title="`Показать QR-код: ${location.name}`"
      @click="onViewQr"
    )
      svg(
        v-if="!isQrLoading"
        xmlns="http://www.w3.org/2000/svg"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        stroke-width="1.8"
        class="h-5 w-5"
        aria-hidden="true"
      )
        rect(x="3" y="3" width="6" height="6" rx="1")
        rect(x="15" y="3" width="6" height="6" rx="1")
        rect(x="3" y="15" width="6" height="6" rx="1")
        path(d="M15 15h2v2h-2zM19 15h2v2h-2zM15 19h2v2h-2zM19 19h2v2h-2z")

  div(class="flex flex-col gap-1 pr-12")
    h3(class="headline-md font-bold") {{ location.name }}
    p(class="body-md opacity-70") {{ location.address }}
    p(class="text-sm opacity-55") /l/{{ location.slug }}

  div(class="flex flex-wrap items-center justify-between gap-2")
    u-badge(:color="statusColor" size="md") {{ statusLabel }}

    div(class="flex items-center gap-2")
      u-button(
        :to="`/admin/locations/${location.id}`"
        color="neutral"
        size="md"
      )
        | Редактировать

      u-button(
        variant="outline"
        color="error"
        size="md"
        :loading="isDeleting"
        :disabled="isBusy"
        @click="onDelete"
      )
        | Удалить
</template>
