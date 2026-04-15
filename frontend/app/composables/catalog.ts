import { computed, ref } from "vue"

export interface Nutrition {
  calories: number
  proteins: number
  fats: number
  carbs: number
}

export interface CatalogItem {
  id: string
  name: string
  description?: string
  category: string
  photo_url: string
  price: number // В копейках.
  stock_amount: number
  nutrition?: Nutrition
}

export interface Location {
  id: string
  name: string
  address: string
  slug: string
}

// Заглушки пока бек не готов.
const MOCK_ITEMS: CatalogItem[] = [
  {
    id: "1",
    name: "Green Sanctuary Bowl",
    description: "A vibrant mix of quinoa, avocado, roast chicken, and fresh greens with a citrus zest dressing.",
    category: "lunch",
    photo_url: "/images/food-placeholder.png",
    price: 45000,
    stock_amount: 5,
    nutrition: { calories: 450, proteins: 28, fats: 18, carbs: 42 },
  },
  {
    id: "2",
    name: "Morning Vigor Bowl",
    description: "Overnight oats with chia seeds, fresh berries, and a touch of organic honey.",
    category: "breakfast",
    photo_url: "/images/food-placeholder.png",
    price: 32000,
    stock_amount: 8,
    nutrition: { calories: 310, proteins: 12, fats: 8, carbs: 54 },
  },
  {
    id: "3",
    name: "Emerald Smoothie",
    description: "Cold-pressed spinach, apple, kale, and ginger for a refreshing energy boost.",
    category: "drinks",
    photo_url: "/images/food-placeholder.png",
    price: 28000,
    stock_amount: 12,
    nutrition: { calories: 120, proteins: 2, fats: 0.5, carbs: 28 },
  },
  {
    id: "4",
    name: "Almond Vitality Bar",
    description: "Raw almonds with dark chocolate and sea salt. The perfect high-protein snack.",
    category: "snacks",
    photo_url: "/images/food-placeholder.png",
    price: 15000,
    stock_amount: 20,
    nutrition: { calories: 190, proteins: 6, fats: 14, carbs: 12 },
  },
]

const MOCK_LOCATIONS: Record<string, Location> = {
  "550e8400-e29b-41d4-a716-446655440000": {
    id: "550e8400-e29b-41d4-a716-446655440000",
    name: "Central Avenue Vending",
    address: "123 Business St, Suite 100",
    slug: "central-ave",
  },
}

export const useCatalog = (locationId: string) => {
  const items = ref<CatalogItem[]>(MOCK_ITEMS)
  const location = ref<Location | null>(MOCK_LOCATIONS[locationId] || null)
  const categories = computed(() => ["all", ...new Set(items.value.map(item => item.category))])
  const selectedCategory = ref("all")

  const filteredItems = computed(() => {
    if (selectedCategory.value === "all") {
      return items.value
    }
    return items.value.filter(item => item.category === selectedCategory.value)
  })

  return {
    items,
    location,
    categories,
    selectedCategory,
    filteredItems,
  }
}
