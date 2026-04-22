export default defineNuxtRouteMiddleware((to) => {
  if (!to.path.startsWith("/admin")) {
    return
  }

  const { isLoggedIn } = useAdminAuth()
  const isLoginPage = to.path === "/admin/login"

  if (!isLoggedIn.value && !isLoginPage) {
    return navigateTo("/admin/login")
  }

  if (isLoggedIn.value && isLoginPage) {
    return navigateTo("/admin")
  }
})
