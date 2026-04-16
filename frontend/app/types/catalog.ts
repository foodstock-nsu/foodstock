export type Range = [number, number]

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
  price: number
  stock_amount: number
  nutrition?: Nutrition
}

export interface Location {
  id: string
  name: string
  address: string
  slug: string
}

export interface Filters {
  calories: Range
  proteins: Range
  fats: Range
  carbs: Range
}

export interface CartItem {
  item: CatalogItem
  quantity: number
}
