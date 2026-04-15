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
    name: "Боул «Зелёный рай»",
    description: "Яркое сочетание киноа, авокадо, запечённой курицы и свежей зелени с цитрусовой заправкой.",
    category: "ужины",
    photo_url: "/images/food-placeholder.png",
    price: 45000,
    stock_amount: 5,
    nutrition: { calories: 450, proteins: 28, fats: 18, carbs: 42 },
  },
  {
    id: "2",
    name: "Боул «Утренний заряд»",
    description: "Ночная овсянка с семенами чиа, свежими ягодами и лёгкой ноткой органического мёда.",
    category: "завтраки",
    photo_url: "/images/food-placeholder.png",
    price: 32000,
    stock_amount: 8,
    nutrition: { calories: 310, proteins: 12, fats: 8, carbs: 54 },
  },
  {
    id: "3",
    name: "Смузи «Изумрудный»",
    description: "Холодный отжим из шпината, яблока, кейла и имбиря для освежающего заряда энергии.",
    category: "напитки",
    photo_url: "/images/food-placeholder.png",
    price: 28000,
    stock_amount: 12,
    nutrition: { calories: 120, proteins: 2, fats: 0.5, carbs: 28 },
  },
  {
    id: "4",
    name: "Батончик «Миндальная энергия»",
    description: "Сырые миндальные орехи с тёмным шоколадом и морской солью. Идеальный высокобелковый перекус.",
    category: "закуски",
    photo_url: "/images/food-placeholder.png",
    price: 15000,
    stock_amount: 20,
    nutrition: { calories: 190, proteins: 6, fats: 14, carbs: 12 },
  },
]

const MOCK_LOCATIONS: Record<string, Location> = {
  "550e8400-e29b-41d4-a716-446655440000": {
    id: "550e8400-e29b-41d4-a716-446655440000",
    name: "Торговый автомат на Центральной",
    address: "ул. Деловая, 123, офис 100",
    slug: "central-ave",
  },
}

export const useCatalog = (locationId: string) => {
  const items = ref<CatalogItem[]>(MOCK_ITEMS)
  const location = ref<Location | null>(MOCK_LOCATIONS[locationId] || null)
  const categories = computed(() => ["Все", ...new Set(items.value.map(item => item.category))])
  const selectedCategory = ref("Все")
  const selectedItem = ref<CatalogItem | null>(null)

  const filteredItems = computed(() => {
    if (selectedCategory.value === "Все") {
      return items.value
    }
    return items.value.filter(item => item.category === selectedCategory.value)
  })

  return {
    items,
    location,
    categories,
    selectedCategory,
    selectedItem,
    filteredItems,
  }
}
