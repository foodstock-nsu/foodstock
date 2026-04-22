import { DEFAULT_FILTERS, MOCK_ITEMS, MOCK_LOCATIONS } from "~/data/mock"
import type { CatalogItem, Filters, Location } from "~/types/catalog"
import { areFiltersActive, cloneFilters, itemMatchesCatalogFilters } from "~/utils/filter"

export const CATEGORY_LABELS: Record<string, string> = {
  Все: "Все",
  lunch: "Обеды",
  breakfast: "Завтраки",
  drinks: "Напитки",
  snacks: "Закуски",
  desserts: "Десерты",
}

export const useCatalog = (locationSlug: string) => {
  const config = useRuntimeConfig()

  const items = ref<CatalogItem[]>([])
  const location = ref<Location | null>(MOCK_LOCATIONS[locationSlug] || null)
  const categories = computed(() => ["Все", ...new Set(items.value.map(item => item.category))])
  const selectedCategory = ref("Все")
  const selectedItem = ref<CatalogItem | null>(null)
  const isLoading = ref(false)
  const error = ref<string | null>(null)

  const filters = reactive<Filters>(cloneFilters(DEFAULT_FILTERS))

  async function loadCatalog() {
    isLoading.value = true
    error.value = null

    try {
      // Пытаемся загрузить данные по slug или id (бекэнд должен поддерживать оба или мы резолвим это)
      const identifier = location.value?.id || locationSlug
      const response = await $fetch<{ items: CatalogItem[] }>(`/api/v1/client/locations/${identifier}/catalog`, {
        method: "GET",
        baseURL: config.public.baseUrl || undefined,
      })

      items.value = response.items || []
    } catch {
      // Endpoint каталога может быть временно недоступен в dev-среде.
      items.value = MOCK_ITEMS
      error.value = "Каталог временно недоступен, показаны тестовые данные"
    } finally {
      isLoading.value = false
    }
  }

  if (import.meta.client) {
    void loadCatalog()
  }

  const resetFilters = () => {
    Object.assign(filters, cloneFilters(DEFAULT_FILTERS))
  }

  const isFiltersActive = computed(() => areFiltersActive(filters, DEFAULT_FILTERS))

  const filteredItems = computed(() => {
    return items.value.filter(item => itemMatchesCatalogFilters(item, selectedCategory.value, filters, isFiltersActive.value))
  })

  return {
    items,
    location,
    categories,
    selectedCategory,
    selectedItem,
    isLoading,
    error,
    filters,
    isFiltersActive,
    resetFilters,
    filteredItems,
    loadCatalog,
  }
}
