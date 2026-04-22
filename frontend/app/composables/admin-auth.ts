import { createGlobalState, useLocalStorage } from "@vueuse/core"

const ADMIN_TOKEN_KEY = "foodstock-admin-token"

function getApiBaseUrl() {
  const config = useRuntimeConfig()
  return config.public.baseUrl || undefined
}

export const useAdminAuth = createGlobalState(() => {
  const token = useLocalStorage<string>(ADMIN_TOKEN_KEY, "")

  const isLoggedIn = computed(() => token.value.length > 0)

  async function login(login: string, password: string): Promise<boolean> {
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
    } catch {
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
    login,
    logout,
    getBearerToken,
  }
})
