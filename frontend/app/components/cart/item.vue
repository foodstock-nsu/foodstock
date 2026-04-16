<script setup lang="ts">
import type { CartItem } from "~/types/catalog"

defineProps<{
  item: CartItem
}>()

const { decrement, increment } = useCartStore()
</script>

<template lang="pug">
div.cart-item
  img.cart-item-img(:src="item.item.photo_url" :alt="item.item.name")
  div.cart-item-info
    p.cart-item-name {{ item.item.name }}
    p.cart-item-price {{ formatNumber(item.item.price / 100) }} ₽
  div.cart-item-controls
    button.cart-qty-btn(:id="`cart-dec-${item.item.id}`" @click="decrement(item.item.id)")
      u-icon(name="i-heroicons-minus")
    span.cart-qty {{ item.quantity }}
    button.cart-qty-btn(:id="`cart-inc-${item.item.id}`" @click="increment(item.item.id)")
      u-icon(name="i-heroicons-plus")
</template>

<style scoped>
.cart-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px 0;
  border-bottom: 1px solid var(--ghost-border);
}

.cart-item:last-child {
  border-bottom: none;
}

.cart-item-img {
  width: 56px;
  height: 56px;
  border-radius: var(--radius-md);
  object-fit: cover;
  flex-shrink: 0;
}

.cart-item-info {
  flex: 1;
  min-width: 0;
}

.cart-item-name {
  font-weight: 600;
  font-size: 0.9375rem;
  color: var(--color-on-surface);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.cart-item-price {
  font-size: 0.875rem;
  color: var(--color-primary);
  font-weight: 700;
  margin-top: 2px;
}

.cart-item-controls {
  display: flex;
  align-items: center;
  gap: 8px;
  background: var(--color-surface-container-low);
  border-radius: var(--radius-full);
  padding: 4px 8px;
}

.cart-qty-btn {
  width: 26px;
  height: 26px;
  border-radius: var(--radius-full);
  border: none;
  background: transparent;
  color: var(--color-primary);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1rem;
  transition: background 0.15s;
}

.cart-qty-btn:hover {
  background: var(--color-secondary-container);
}

.cart-qty {
  font-weight: 700;
  font-size: 0.9375rem;
  color: var(--color-on-surface);
  min-width: 20px;
  text-align: center;
}
</style>
