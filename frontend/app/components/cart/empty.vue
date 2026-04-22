<template lang="pug">
div.cart-empty
  u-icon(name="i-heroicons-shopping-cart" class="cart-empty-icon")
  p.body-md.cart-empty-text Ваша корзина пуста
  p.body-md.cart-empty-sub Добавьте что-нибудь вкусное
</template>

<style scoped>
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
</style>
