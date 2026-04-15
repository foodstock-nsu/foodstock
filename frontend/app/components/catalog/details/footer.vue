<script setup lang="ts">
const props = defineProps<{
  item: CatalogItem
}>()

const { items, addItem } = useCartStore()

const adding = ref(false)

function addToCart() {
  addItem(props.item)
  adding.value = true
  setTimeout(() => {
    adding.value = false
  }, 400)
}

const addLabel = computed(() => {
  const cartItem = items.value.find((ci: CartItem) => ci.item.id === props.item.id)
  return cartItem ? `В корзине (${cartItem.quantity})` : "Добавить в корзину"
})
</script>

<template lang="pug">
div(class="footer-container")
  u-button(
    block
    size="xl"
    class="btn-primary py-4 text-lg font-bold shadow-lg"
    :class="{ 'add-btn--bounce': adding }"
    :disabled="item.stock_amount === 0"
    @click="addToCart"
  ) {{ item.stock_amount > 0 ? addLabel : 'Нет в наличии' }}
</template>

<style scoped>
.footer-container {
  padding: 16px 0;
  background: linear-gradient(to top, var(--color-surface-container-low) 80%, transparent);
}

@keyframes add-bounce {
  0% { transform: scale(1); }
  40% { transform: scale(0.95); }
  75% { transform: scale(1.02); }
  100% { transform: scale(1); }
}

.add-btn--bounce {
  animation: add-bounce 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
}
</style>
