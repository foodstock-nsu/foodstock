<script setup lang="ts">
import type { CatalogItem } from "~/types/catalog"

defineProps<{
  item: CatalogItem | null
}>()

const open = defineModel<boolean>("open", { required: true })
</script>

<template lang="pug">
u-drawer(
  v-model:open="open"
  direction="bottom"
  :ui="{ content: 'max-h-[90dvh] flex flex-col --color-surface-container-low rounded-t-3xl overflow-hidden' }"
)
  template(#header)
    div(class="flex justify-between items-center w-full px-2")
      span(class="headline-sm font-bold opacity-0") Details
      u-button(
        icon="i-heroicons-x-mark"
        variant="ghost"
        color="primary"
        class="rounded-full"
        @click="open = false"
      )

  template(#body)
    div(v-if="item" class="flex flex-col gap-6 pb-4")
      div(class="aspect-[4/3] w-full overflow-hidden rounded-2xl bg-surface-container-low shadow-inner")
        img(
          :src="item.photo_url"
          :alt="item.name"
          class="w-full h-full object-cover"
        )

      div(class="flex flex-col gap-2")
        div(class="flex justify-between items-baseline")
          h2(class="display-sm font-extrabold text-on-surface") {{ item.name }}
          span(class="headline-md text-primary font-bold") {{ formatNumber(item.price / 100) }} ₽

        p(class="body-lg text-on-surface opacity-80 leading-relaxed") {{ item.description }}

      div(v-if="item.nutrition" class="mt-4")
        h3(class="label-lg uppercase tracking-wider text-on-surface opacity-50 mb-4") Пищевая ценность
        catalog-details-nutrition(:nutrition="item.nutrition")

  template(v-if="item" #footer)
    catalog-details-footer(:item="item")
</template>
