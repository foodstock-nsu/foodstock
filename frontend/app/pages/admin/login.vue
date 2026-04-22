<script setup lang="ts">
const { login, isLoggedIn } = useAdminAuth()

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
    error.value = "Неверный логин или пароль"
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
        input#login.input-minimal.w-full.p-4(
          v-model="loginValue"
          type="text"
          placeholder="Введите логин"
          required
        )

      .text-left
        label.block.font-semibold.mb-2(for="password") Пароль
        input#password.input-minimal.w-full.p-4(
          v-model="passwordValue"
          type="password"
          placeholder="Введите пароль"
          required
        )

      p.text-red-500.text-sm(v-if="error") {{ error }}

      button.btn-primary.w-full.p-4.mt-4.transition-all.duration-200(type="submit" :disabled="loading")
        span(v-if="!loading") Войти
        span(v-else) Входим...
</template>
