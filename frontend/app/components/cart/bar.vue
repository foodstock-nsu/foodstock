<script setup lang="ts">
const { totalPrice, totalQuantity } = useCartStore()

const drawerOpen = ref(false)
</script>

<template lang="pug">
transition(name="cart-bar")
  div.cart-bar(
    v-show="totalQuantity > 0"
    id="cart-bar"
    role="button"
    tabindex="0"
    aria-label="Open cart"
    @click="drawerOpen = true"
    @keydown.enter="drawerOpen = true"
  )
    div.cart-bar-left
      div.cart-bar-icon
        u-icon(name="i-heroicons-shopping-bag" class="cart-icon-svg")
        span.cart-bar-count {{ totalQuantity }}
      span.cart-bar-label Корзина

    span.cart-bar-total {{ formatNumber(totalPrice / 100) }} ₽

cart-drawer(v-model:open="drawerOpen")
</template>

<style scoped>
.cart-bar {
  position: fixed;
  bottom: 24px;
  left: 50%;
  transform: translateX(-50%);
  z-index: 100;

  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 20px;

  background: linear-gradient(120deg, var(--color-primary) 0%, var(--color-primary-container) 100%);
  color: var(--color-on-primary);
  border-radius: var(--radius-full);
  padding: 14px 20px 14px 16px;
  min-width: 260px;
  max-width: calc(100vw - 48px);

  cursor: pointer;
  user-select: none;

  box-shadow:
    0 8px 24px rgb(0 108 73 / 35%),
    0 2px 8px rgb(0 108 73 / 20%);

  transition:
    transform 0.2s ease,
    box-shadow 0.2s ease;
}

.cart-bar:hover {
  transform: translateX(-50%) translateY(-2px);
  box-shadow:
    0 12px 32px rgb(0 108 73 / 40%),
    0 4px 12px rgb(0 108 73 / 25%);
}

.cart-bar:active {
  transform: translateX(-50%) translateY(0);
}

.cart-bar-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.cart-bar-icon {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  background: rgb(255 255 255 / 20%);
  border-radius: var(--radius-full);
}

.cart-icon-svg {
  width: 18px;
  height: 18px;
}

.cart-bar-count {
  position: absolute;
  top: -4px;
  right: -4px;
  background: #fff;
  color: var(--color-primary);
  font-size: 0.625rem;
  font-weight: 800;
  min-width: 16px;
  height: 16px;
  border-radius: var(--radius-full);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0 3px;
  line-height: 1;
}

.cart-bar-label {
  font-weight: 700;
  font-size: 0.9375rem;
  font-family: var(--font-display);
}

.cart-bar-total {
  font-weight: 800;
  font-size: 1rem;
  font-family: var(--font-display);
  letter-spacing: -0.02em;
  white-space: nowrap;
}

/* Entrance / exit animation */
.cart-bar-enter-active,
.cart-bar-leave-active {
  transition: all 0.35s cubic-bezier(0.34, 1.56, 0.64, 1);
}

.cart-bar-enter-from,
.cart-bar-leave-to {
  opacity: 0;
  transform: translateX(-50%) translateY(32px) scale(0.9);
}
</style>
