import type {
  AdminItem,
  AdminItemForm,
  AdminLocation,
  AdminLocationForm,
  AdminStats,
} from "~/types/admin"

const MOCK_DELAY_MS = 110

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

  const isLoading = ref(false)
  const error = ref<string | null>(null)

  async function runMock<T>(task: () => T): Promise<T> {
    isLoading.value = true
    error.value = null
    try {
      await sleep(MOCK_DELAY_MS)
      return task()
    } catch {
      error.value = "Не удалось выполнить операцию"
      throw new Error(error.value)
    } finally {
      isLoading.value = false
    }
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
    return runMock(() => store.items.value.find(item => item.id === id) || null)
  }

  async function getLocationById(id: string) {
    return runMock(() => store.locations.value.find(location => location.id === id) || null)
  }

  async function saveItem(form: AdminItemForm, id?: string) {
    return runMock(() => {
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
    return runMock(() => {
      const payload: AdminLocation = {
        id: id || createId(),
        slug: form.slug.trim(),
        name: form.name.trim(),
        address: form.address.trim(),
        is_active: form.is_active,
      }

      const index = store.locations.value.findIndex(location => location.id === payload.id)
      if (index === -1) {
        store.locations.value = [payload, ...store.locations.value]
      } else {
        store.locations.value[index] = payload
      }

      return payload
    })
  }

  async function removeItem(id: string) {
    return runMock(() => {
      store.items.value = store.items.value.filter(item => item.id !== id)
    })
  }

  async function removeLocation(id: string) {
    return runMock(() => {
      store.locations.value = store.locations.value.filter(location => location.id !== id)
    })
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
  }
}
