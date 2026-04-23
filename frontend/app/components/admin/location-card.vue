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
article(class="surface-card container-pad relative flex flex-col gap-4 transition-all duration-300 hover:shadow-lg")
  div(v-if="location.is_active" class="absolute right-4 top-4")
    u-button(
      color="neutral"
      variant="ghost"
      class="h-10 w-10 !p-0 inline-flex items-center justify-center rounded-full transition-all duration-300 transform active:scale-90"
      :loading="isQrLoading"
      :disabled="isBusy"
      :title="`Показать QR-код: ${location.name}`"
      @click="onViewQr"
    )
      u-icon(
        v-if="!isQrLoading"
        name="i-heroicons-qr-code"
        class="h-6 w-6"
      )

  div(class="flex flex-col gap-1 pr-12")
    h3(class="headline-md font-extrabold") {{ location.name }}
    p(class="body-md opacity-60 flex items-center gap-1.5")
      u-icon(name="i-heroicons-map-pin" class="w-4 h-4")
      | {{ location.address }}
    p(class="text-xs font-mono opacity-40") /l/{{ location.slug }}

  div(class="flex flex-wrap items-center justify-between gap-4 mt-2")
    u-badge(:color="statusColor" size="md" class="px-3") {{ statusLabel }}

    div(class="flex items-center gap-3")
      u-button(
        :to="`/admin/locations/${location.id}`"
        variant="ghost"
        class="btn-secondary px-5 py-2 rounded-full transition-all duration-300 transform active:scale-95"
      )
        | Редактировать

      u-button(
        variant="ghost"
        color="error"
        class="hover:bg-red-50 dark:hover:bg-red-900/10 px-5 py-2 rounded-full transition-all duration-300 transform active:scale-95"
        :loading="isDeleting"
        :disabled="isBusy"
        @click="onDelete"
      )
        | Удалить
</template>
