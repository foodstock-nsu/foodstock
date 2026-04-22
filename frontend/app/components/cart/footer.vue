<script setup lang="ts">
import { useRouter } from "vue-router"

const { totalPrice, clear } = useCartStore()
const router = useRouter()

function handleCheckout() {
  clear()
  router.push("/payment")
}
</script>

<template lang="pug">
div.cart-footer
  div.cart-total-row
    span.body-md Итого
    span.cart-total-price {{ formatNumber(totalPrice / 100) }} ₽
  button.btn-primary.cart-checkout-btn(id="cart-checkout-btn" @click="handleCheckout") Оформить заказ
</template>

<style scoped>
.cart-footer {
  display: flex;
  flex-direction: column;
  gap: 14px;
  width: 100%;
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
</style>
