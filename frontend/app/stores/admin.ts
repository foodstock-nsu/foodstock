import { createGlobalState, useLocalStorage } from "@vueuse/core"

import { MOCK_ITEMS } from "~/data/mock"
import type { AdminCategory, AdminItem, AdminLocation } from "~/types/admin"

const ALL_CATEGORIES = "all"

function cloneItems(): AdminItem[] {
  return MOCK_ITEMS.map(item => ({
    ...item,
    nutrition: item.nutrition ? { ...item.nutrition } : undefined,
  }))
}

export const useAdminStore = createGlobalState(() => {
  const items = useLocalStorage<AdminItem[]>("foodstock-admin-items", cloneItems())
  const locations = useLocalStorage<AdminLocation[]>("foodstock-admin-locations", [])
  const isLocationsLoaded = ref(false)

  const itemSearch = useLocalStorage("foodstock-admin-item-search", "")
  const itemCategory = useLocalStorage<AdminCategory | typeof ALL_CATEGORIES>("foodstock-admin-item-category", ALL_CATEGORIES)
  const locationSearch = useLocalStorage("foodstock-admin-location-search", "")

  const normalizedItemSearch = computed(() => itemSearch.value.trim().toLowerCase())
  const normalizedLocationSearch = computed(() => locationSearch.value.trim().toLowerCase())

  const categories = computed(() => {
    const set = new Set<AdminCategory>()
    items.value.forEach((item) => {
      if (item.category) {
        set.add(item.category)
      }
    })
    return Array.from(set)
  })

  const filteredItems = computed(() => {
    return items.value.filter((item) => {
      const matchesSearch = !normalizedItemSearch.value
        || item.name.toLowerCase().includes(normalizedItemSearch.value)
        || item.description?.toLowerCase().includes(normalizedItemSearch.value)

      const matchesCategory = itemCategory.value === ALL_CATEGORIES || item.category === itemCategory.value

      return matchesSearch && matchesCategory
    })
  })

  const filteredLocations = computed(() => {
    return locations.value.filter((location) => {
      if (!normalizedLocationSearch.value) {
        return true
      }
      return (
        location.name.toLowerCase().includes(normalizedLocationSearch.value)
        || location.slug.toLowerCase().includes(normalizedLocationSearch.value)
        || location.address.toLowerCase().includes(normalizedLocationSearch.value)
      )
    })
  })

  const sortedItems = computed(() => {
    return [...filteredItems.value].sort((a, b) => a.name.localeCompare(b.name, "ru"))
  })

  const sortedLocations = computed(() => {
    return [...filteredLocations.value].sort((a, b) => a.name.localeCompare(b.name, "ru"))
  })

  return {
    items,
    locations,
    isLocationsLoaded,
    itemSearch,
    itemCategory,
    locationSearch,
    categories,
    sortedItems,
    sortedLocations,
  }
})
