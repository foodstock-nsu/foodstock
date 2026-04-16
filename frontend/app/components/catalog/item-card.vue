<script setup lang="ts">
import type { CatalogItem } from "~/types/catalog"

const props = defineProps<{
  item: CatalogItem
}>()

const emit = defineEmits<{
  (e: "select", item: CatalogItem): void
}>()

const { items, addItem } = useCartStore()

const adding = ref(false)

function addToCart(e: Event) {
  e.stopPropagation()
  addItem(props.item)
  adding.value = true
  setTimeout(() => {
    adding.value = false
  }, 400)
}

const addLabel = computed(() => {
  const cartItem = items.value.find((ci: CartItem) => ci.item.id === props.item.id)
  return cartItem ? `В корзине (${cartItem.quantity})` : "Добавить"
})
</script>

<template lang="pug">
div(
  class="surface-card p-4 flex flex-col gap-4 relative transition-all duration-300 hover:-translate-y-1 hover:shadow-xl cursor-pointer"
  @click="emit('select', item)"
)
  div(class="ingredient-float aspect-square overflow-hidden rounded-md bg-surface-container-low")
    img(
      :src="item.photo_url"
      :alt="item.name"
      class="w-full h-full object-cover food-image"
    )

  div(class="flex flex-col gap-1")
    div(class="flex justify-between items-start")
      h3(class="headline-md text-on-surface line-clamp-1") {{ item.name }}
      div(class="text-primary font-bold text-lg") {{ formatNumber(item.price / 100) }} ₽

    p(class="body-md text-on-surface opacity-70 line-clamp-2") {{ item.description }}

  div(class="mt-auto flex items-center justify-between")
    div(v-if="item.stock_amount > 0" class="flex items-center gap-2")
      div(class="w-2 h-2 rounded-full bg-primary")
      span(class="text-xs font-medium uppercase tracking-wider text-on-surface opacity-60") {{ item.stock_amount }} в наличии
    div(v-else class="text-xs font-medium uppercase text-red-500") Нет в наличии

    button(
      :id="`add-to-cart-${item.id}`"
      class="btn-primary px-6 py-2 text-sm add-btn"
      :class="{ 'add-btn--bounce': adding }"
      :disabled="item.stock_amount === 0"
      @click="addToCart"
    ) {{ addLabel }}
</template>

<style scoped>
.line-clamp-1 {
  display: -webkit-box;
  -webkit-line-clamp: 1;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.add-btn {
  transition:
    background 0.2s,
    transform 0.15s,
    opacity 0.15s;
}

.add-btn:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

@keyframes add-bounce {
  0% { transform: scale(1); }
  40% { transform: scale(0.92); }
  75% { transform: scale(1.06); }
  100% { transform: scale(1); }
}

.add-btn--bounce {
  animation: add-bounce 0.35s ease;
}
</style>
