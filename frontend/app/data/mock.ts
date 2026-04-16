import type { CatalogItem, Filters, Location, Range } from "~/types/catalog"

const createRange = (min: number, max: number): Range => [min, max]

export const DEFAULT_FILTERS: Filters = {
  calories: createRange(0, 1000),
  proteins: createRange(0, 100),
  fats: createRange(0, 100),
  carbs: createRange(0, 200),
}

export const MOCK_ITEMS: CatalogItem[] = [
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

export const MOCK_LOCATIONS: Record<string, Location> = {
  "550e8400-e29b-41d4-a716-446655440000": {
    id: "550e8400-e29b-41d4-a716-446655440000",
    name: "Торговый автомат на Центральной",
    address: "ул. Деловая, 123, офис 100",
    slug: "central-ave",
  },
}
