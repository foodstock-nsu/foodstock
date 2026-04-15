import { createGlobalState, useLocalStorage } from "@vueuse/core"

export interface CartItem {
  item: CatalogItem
  quantity: number
}

export const useCartStore = createGlobalState(() => {
  const items = useLocalStorage<CartItem[]>("foodstock-cart", [])

  const totalQuantity = computed(() => items.value.reduce((sum, ci) => sum + ci.quantity, 0))
  const totalPrice = computed(() => items.value.reduce((sum, ci) => sum + ci.item.price * ci.quantity, 0))

  function addItem(item: CatalogItem) {
    const existing = items.value.find(ci => ci.item.id === item.id)
    if (existing) {
      existing.quantity++
    } else {
      items.value.push({ item, quantity: 1 })
    }
  }

  function removeItem(itemId: string) {
    items.value = items.value.filter(ci => ci.item.id !== itemId)
  }

  function increment(itemId: string) {
    const entry = items.value.find(ci => ci.item.id === itemId)
    if (entry) {
      entry.quantity++
    }
  }

  function decrement(itemId: string) {
    const entry = items.value.find(ci => ci.item.id === itemId)
    if (!entry) {
      return
    }
    if (entry.quantity <= 1) {
      removeItem(itemId)
    } else {
      entry.quantity--
    }
  }

  function clear() {
    items.value = []
  }

  return {
    items,
    totalQuantity,
    totalPrice,
    addItem,
    removeItem,
    increment,
    decrement,
    clear,
  }
})
