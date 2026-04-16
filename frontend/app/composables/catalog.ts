import { DEFAULT_FILTERS, MOCK_ITEMS, MOCK_LOCATIONS } from "~/data/mock"
import type { CatalogItem, Filters, Location } from "~/types/catalog"
import { areFiltersActive, cloneFilters, itemMatchesCatalogFilters } from "~/utils/filter"

export const useCatalog = (locationId: string) => {
  const items = ref<CatalogItem[]>(MOCK_ITEMS)
  const location = ref<Location | null>(MOCK_LOCATIONS[locationId] || null)
  const categories = computed(() => ["Все", ...new Set(items.value.map(item => item.category))])
  const selectedCategory = ref("Все")
  const selectedItem = ref<CatalogItem | null>(null)

  const filters = reactive<Filters>(cloneFilters(DEFAULT_FILTERS))

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
    filters,
    isFiltersActive,
    resetFilters,
    filteredItems,
  }
}
