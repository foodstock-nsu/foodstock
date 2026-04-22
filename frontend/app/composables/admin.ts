import type {
  AdminItem,
  AdminItemForm,
  AdminLocation,
  AdminLocationForm,
  AdminStats,
} from "~/types/admin"

const STUB_DELAY_MS = 110

function sleep(ms: number) {
  return new Promise(resolve => setTimeout(resolve, ms))
}

function createId() {
  if (typeof crypto !== "undefined" && typeof crypto.randomUUID === "function") {
    return crypto.randomUUID()
  }
  return `local-${Date.now()}-${Math.floor(Math.random() * 1000)}`
}

export const useAdmin = () => {
  const store = useAdminStore()
  const auth = useAdminAuth()
  const config = useRuntimeConfig()

  const isLoading = ref(false)
  const error = ref<string | null>(null)

  async function runStub<T>(task: () => T): Promise<T> {
    isLoading.value = true
    error.value = null
    try {
      await sleep(STUB_DELAY_MS)
      return task()
    } catch {
      error.value = "Не удалось выполнить операцию"
      throw new Error(error.value)
    } finally {
      isLoading.value = false
    }
  }

  function getApiBaseUrl() {
    return config.public.baseUrl || undefined
  }

  function withAuthHeaders() {
    const token = auth.getBearerToken()
    if (!token) {
      throw new Error("Требуется авторизация")
    }

    return {
      Authorization: `Bearer ${token}`,
    }
  }

  function mapApiLocation(location: {
    id: string
    slug: string
    name: string
    address: string
    is_active: boolean
  }): AdminLocation {
    return {
      id: location.id,
      slug: location.slug,
      name: location.name,
      address: location.address,
      is_active: location.is_active,
    }
  }

  function parseErrorMessage(err: unknown, fallback: string) {
    if (err && typeof err === "object" && "data" in err) {
      const data = (err as { data?: { error?: string } }).data
      if (data?.error) {
        return data.error
      }
    }

    if (err instanceof Error && err.message) {
      return err.message
    }

    return fallback
  }

  async function runApi<T>(task: () => Promise<T>, fallbackError: string): Promise<T> {
    isLoading.value = true
    error.value = null
    try {
      return await task()
    } catch (err) {
      error.value = parseErrorMessage(err, fallbackError)
      throw new Error(error.value)
    } finally {
      isLoading.value = false
    }
  }

  async function loadLocations() {
    return runApi(async () => {
      const response = await $fetch<{ locations: Array<{
        id: string
        slug: string
        name: string
        address: string
        is_active: boolean
      }> }>("/api/v1/admin/locations", {
        method: "GET",
        baseURL: getApiBaseUrl(),
        headers: withAuthHeaders(),
      })

      store.locations.value = (response.locations || []).map(mapApiLocation)
      store.isLocationsLoaded.value = true
      return store.locations.value
    }, "Не удалось загрузить локации")
  }

  if (import.meta.client && auth.isLoggedIn.value && !store.isLocationsLoaded.value) {
    void loadLocations()
  }

  const stats = computed<AdminStats>(() => {
    const locations = store.locations.value
    const items = store.items.value
    return {
      itemsCount: items.length,
      lowStockCount: items.filter((item) => {
        const stockAmount = item.stock_amount || 0
        return stockAmount <= 3
      }).length,
      locationsCount: locations.length,
      activeLocationsCount: locations.filter(location => location.is_active !== false).length,
    }
  })

  function createItemForm(item?: AdminItem | null): AdminItemForm {
    const priceRub = item?.price ? Math.floor(item.price / 100) : 0

    return {
      name: item?.name ?? "",
      description: item?.description ?? "",
      category: item?.category ?? "lunch",
      photo_url: item?.photo_url ?? "/images/food-placeholder.png",
      priceRub,
      stock_amount: item?.stock_amount ?? 0,
      nutrition: {
        calories: item?.nutrition?.calories ?? 0,
        proteins: item?.nutrition?.proteins ?? 0,
        fats: item?.nutrition?.fats ?? 0,
        carbs: item?.nutrition?.carbs ?? 0,
      },
    }
  }

  function createLocationForm(location?: AdminLocation | null): AdminLocationForm {
    return {
      slug: location?.slug ?? "",
      name: location?.name ?? "",
      address: location?.address ?? "",
      is_active: location?.is_active ?? true,
    }
  }

  async function getItemById(id: string) {
    return runStub(() => store.items.value.find(item => item.id === id) || null)
  }

  async function getLocationById(id: string) {
    if (!store.isLocationsLoaded.value && auth.isLoggedIn.value) {
      await loadLocations()
    }

    return store.locations.value.find(location => location.id === id) || null
  }

  async function saveItem(form: AdminItemForm, id?: string) {
    return runStub(() => {
      const payload: AdminItem = {
        id: id || createId(),
        name: form.name.trim(),
        description: form.description.trim(),
        category: form.category,
        photo_url: form.photo_url.trim() || "/images/food-placeholder.png",
        price: Math.max(0, Math.round(form.priceRub * 100)),
        stock_amount: Math.max(0, Math.round(form.stock_amount)),
        nutrition: {
          calories: Math.max(0, Math.round(form.nutrition.calories ?? 0)),
          proteins: Number(form.nutrition.proteins ?? 0),
          fats: Number(form.nutrition.fats ?? 0),
          carbs: Number(form.nutrition.carbs ?? 0),
        },
      }

      const index = store.items.value.findIndex(item => item.id === payload.id)
      if (index === -1) {
        store.items.value = [payload, ...store.items.value]
      } else {
        store.items.value[index] = payload
      }

      return payload
    })
  }

  async function saveLocation(form: AdminLocationForm, id?: string) {
    return runApi(async () => {
      const payload = {
        slug: form.slug.trim(),
        name: form.name.trim(),
        address: form.address.trim(),
      }

      if (!id) {
        const response = await $fetch<{ location: {
          id: string
          slug: string
          name: string
          address: string
          is_active: boolean
        } }>("/api/v1/admin/locations", {
          method: "POST",
          baseURL: getApiBaseUrl(),
          headers: withAuthHeaders(),
          body: payload,
        })

        const mapped = mapApiLocation(response.location)
        store.locations.value = [mapped, ...store.locations.value]
        store.isLocationsLoaded.value = true
        return mapped
      }

      const response = await $fetch<{ location: {
        id: string
        slug: string
        name: string
        address: string
        is_active: boolean
      } }>(`/api/v1/admin/locations/${id}`, {
        method: "PUT",
        baseURL: getApiBaseUrl(),
        headers: withAuthHeaders(),
        body: {
          ...payload,
          is_active: form.is_active,
        },
      })

      const mapped = mapApiLocation(response.location)
      const index = store.locations.value.findIndex(location => location.id === mapped.id)
      if (index === -1) {
        store.locations.value = [mapped, ...store.locations.value]
      } else {
        store.locations.value[index] = mapped
      }
      store.isLocationsLoaded.value = true
      return mapped
    }, "Не удалось сохранить локацию")
  }

  async function removeItem(id: string) {
    return runStub(() => {
      store.items.value = store.items.value.filter(item => item.id !== id)
    })
  }

  async function removeLocation(id: string) {
    return runApi(async () => {
      await $fetch(`/api/v1/admin/locations/${id}`, {
        method: "DELETE",
        baseURL: getApiBaseUrl(),
        headers: withAuthHeaders(),
      })

      store.locations.value = store.locations.value.filter(location => location.id !== id)
    }, "Не удалось удалить локацию")
  }

  return {
    store,
    isLoading,
    error,
    stats,
    createItemForm,
    createLocationForm,
    getItemById,
    getLocationById,
    saveItem,
    saveLocation,
    removeItem,
    removeLocation,
    loadLocations,
  }
}
