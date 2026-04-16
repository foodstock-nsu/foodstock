import type { components } from "#open-fetch-schemas/api"

export type Range = [number, number]

export type Nutrition = NonNullable<components["schemas"]["Item"]["nutrition"]>

export type CatalogItem = components["schemas"]["CatalogItem"]

export type Location = components["schemas"]["Location"]

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
