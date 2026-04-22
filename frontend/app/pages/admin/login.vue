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
  .surface-card.w-full.max-w-md.p-8.md.p-12.flex.flex-col.items-center.text-center.transition-all.duration-500.hover.shadow-2xl
    .w-16.h-16.rounded-2xl.bg-primary.flex.items-center.justify-center.mb-8.shadow-lg
      u-icon(name="i-heroicons-lock-closed" class="w-8 h-8 text-white")

    h1.display-lg.font-black.mb-2 Foodstock
    p.body-md.opacity-60.mb-10 Вход в панель управления

    form.w-full.space-y-6(@submit.prevent="handleLogin")
      .text-left
        label.block.text-sm.font-bold.mb-2.opacity-60(for="login") Логин
        u-input#login(
          v-model="loginValue"
          type="text"
          size="xl"
          placeholder="admin"
          required
          class="rounded-full"
        )

      .text-left
        label.block.text-sm.font-bold.mb-2.opacity-60(for="password") Пароль
        u-input#password(
          v-model="passwordValue"
          type="password"
          size="xl"
          placeholder="••••••••"
          required
          class="rounded-full"
        )

      p.text-red-500.text-sm.font-bold(v-if="error") {{ error }}

      u-button.w-full.btn-primary.py-4.text-lg.transform.transition-all.active:scale-95(
        type="submit"
        size="xl"
        :loading="loading"
      )
        | Войти в систему
</template>
