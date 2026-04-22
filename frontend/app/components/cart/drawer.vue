<script setup lang="ts">
const { items, totalQuantity } = useCartStore()

const open = defineModel<boolean>("open", { required: true })
</script>

<template lang="pug">
u-drawer(v-model:open="open" direction="bottom" :ui="{ content: 'max-h-[85dvh] flex flex-col --color-surface-container-low' }")
  template(#header)
    div.cart-header
      span.headline-md.font-bold Корзина
      span.cart-badge(v-if="totalQuantity > 0") {{ totalQuantity }}
      u-button(
        id="cart-close-btn"
        icon="i-heroicons-x-mark"
        variant="ghost"
        color="primary"
        @click="open = false"
      )

  template(#body)
    cart-empty(v-if="items.length === 0")
    div.cart-item-list(v-else)
      transition-group(name="cart-item")
        cart-item(v-for="ci in items" :key="ci.item.id" :item="ci")

  template(v-if="items.length > 0" #footer)
    cart-footer
</template>

<style scoped>
/* Header */
.cart-header {
  display: flex;
  align-items: center;
  gap: 12px;
}

.cart-header .headline-md {
  flex: 1;
}

.cart-badge {
  background: linear-gradient(120deg, var(--color-primary) 0%, var(--color-primary-container) 100%);
  color: var(--color-on-primary);
  font-weight: 700;
  font-size: 0.75rem;
  min-width: 22px;
  height: 22px;
  border-radius: var(--radius-full);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0 6px;
}

/* Item list */
.cart-item-list {
  display: flex;
  flex-direction: column;
  overflow-y: auto;
}

/* Transitions */
.cart-item-enter-active,
.cart-item-leave-active {
  transition: all 0.25s ease;
}

.cart-item-enter-from,
.cart-item-leave-to {
  opacity: 0;
  transform: translateX(-16px);
}
</style>
