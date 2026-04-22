<script setup lang="ts">
const { login, isLoggedIn, lastError } = useAdminAuth()

const loginValue = ref("")
const passwordValue = ref("")
const loading = ref(false)
const error = ref("")

definePageMeta({
  layout: false,
})

onMounted(() => {
  if (isLoggedIn.value) {
    navigateTo("/admin")
  }
})

async function handleLogin() {
  loading.value = true
  error.value = ""

  const success = await login(loginValue.value, passwordValue.value)

  if (success) {
    await navigateTo("/admin")
  } else {
    error.value = lastError.value || "Не удалось выполнить вход"
  }

  loading.value = false
}
</script>

<template lang="pug">
.login-page.min-h-screen.surface-base.flex.items-center.justify-center.container-pad
  .surface-card.w-full.max-w-md.p-8.md.p-10
    h1.headline-lg.mb-2 Foodstock
    p.body-md.text-on-surface.opacity-60.mb-8 Вход в админ-панель

    form.space-y-6(@submit.prevent="handleLogin")
      .text-left
        label.block.font-semibold.mb-2(for="login") Логин
        u-input#login(
          v-model="loginValue"
          type="text"
          size="xl"
          placeholder="Введите логин"
          required
        )

      .text-left
        label.block.font-semibold.mb-2(for="password") Пароль
        u-input#password(
          v-model="passwordValue"
          type="password"
          size="xl"
          placeholder="Введите пароль"
          required
        )

      p.text-red-500.text-sm(v-if="error") {{ error }}

      u-button.w-full(type="submit" size="xl" :loading="loading")
        | Войти
</template>
