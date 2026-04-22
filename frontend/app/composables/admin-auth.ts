import { createGlobalState, useLocalStorage } from "@vueuse/core"

const ADMIN_TOKEN_KEY = "foodstock-admin-token"

function getApiBaseUrl() {
  const config = useRuntimeConfig()
  return config.public.baseUrl || undefined
}

export const useAdminAuth = createGlobalState(() => {
  const token = useLocalStorage<string>(ADMIN_TOKEN_KEY, "")
  const lastError = ref("")

  const isLoggedIn = computed(() => token.value.length > 0)

  function mapAuthApiError(err: unknown) {
    if (!err || typeof err !== "object") {
      return "Не удалось выполнить вход"
    }

    const maybeError = err as {
      statusCode?: number
      response?: { status?: number }
      data?: { error?: string }
    }

    const statusCode = maybeError.statusCode ?? maybeError.response?.status
    const errorCode = String(maybeError.data?.error || "").trim().toLowerCase()

    if (errorCode === "invalid login or password" || statusCode === 401) {
      return "Неверный логин или пароль"
    }

    if (errorCode === "invalid json" || statusCode === 400) {
      return "Некорректный формат запроса"
    }

    if (errorCode === "internal error" || statusCode === 500) {
      return "Сервер временно недоступен. Попробуйте позже"
    }

    return "Не удалось выполнить вход"
  }

  async function login(login: string, password: string): Promise<boolean> {
    lastError.value = ""
    try {
      const payload = await $fetch<{ token: string }>("/api/v1/admin/auth", {
        method: "POST",
        baseURL: getApiBaseUrl(),
        body: {
          login: login.trim(),
          password,
        },
      })

      if (!payload?.token) {
        return false
      }

      token.value = payload.token
      return true
    } catch (err) {
      lastError.value = mapAuthApiError(err)
      return false
    }
  }

  function getBearerToken() {
    return token.value || ""
  }

  function logout() {
    token.value = ""
    navigateTo("/admin/login")
  }

  return {
    token,
    isLoggedIn,
    lastError,
    login,
    logout,
    getBearerToken,
  }
})
