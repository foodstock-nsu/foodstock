<script setup lang="ts">
const { items, totalQuantity, totalPrice, decrement, increment } = useCartStore()

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
        color="neutral"
        @click="open = false"
      )

  template(#body)
    //- Empty state
    div.cart-empty(v-if="items.length === 0")
      u-icon(name="i-heroicons-shopping-cart" class="cart-empty-icon")
      p.body-md.cart-empty-text Ваша корзина пуста
      p.body-md.cart-empty-sub Добавьте что-нибудь вкусное

    //- Item list
    div.cart-item-list(v-else)
      transition-group(name="cart-item")
        div.cart-item(v-for="ci in items" :key="ci.item.id")
          img.cart-item-img(:src="ci.item.photo_url" :alt="ci.item.name")
          div.cart-item-info
            p.cart-item-name {{ ci.item.name }}
            p.cart-item-price {{ formatNumber(ci.item.price / 100) }} ₽
          div.cart-item-controls
            button.cart-qty-btn(:id="`cart-dec-${ci.item.id}`" @click="decrement(ci.item.id)")
              u-icon(name="i-heroicons-minus")
            span.cart-qty {{ ci.quantity }}
            button.cart-qty-btn(:id="`cart-inc-${ci.item.id}`" @click="increment(ci.item.id)")
              u-icon(name="i-heroicons-plus")

  template(v-if="items.length > 0" #footer)
    div.cart-footer
      div.cart-total-row
        span.body-md Итого
        span.cart-total-price {{ formatNumber(totalPrice / 100) }} ₽
      button.btn-primary.cart-checkout-btn(id="cart-checkout-btn") Оформить заказ
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

/* Empty state */
.cart-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  padding: 40px 0 48px;
  text-align: center;
}

.cart-empty-icon {
  width: 64px;
  height: 64px;
  color: var(--color-primary);
  opacity: 0.2;
  margin-bottom: 8px;
}

.cart-empty-text {
  font-weight: 600;
  color: var(--color-on-surface);
}

.cart-empty-sub {
  color: var(--color-on-surface);
  opacity: 0.5;
}

/* Item list */
.cart-item-list {
  display: flex;
  flex-direction: column;
  overflow-y: auto;
}

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

/* Footer */
.cart-footer {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.cart-total-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.cart-total-price {
  font-family: var(--font-display);
  font-size: 1.25rem;
  font-weight: 800;
  color: var(--color-on-surface);
  letter-spacing: -0.02em;
}

.cart-checkout-btn {
  width: 100%;
  padding: 14px;
  font-size: 1rem;
  letter-spacing: 0.01em;
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
