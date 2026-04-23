import type { components } from "#open-fetch-schemas/api"

export type Range = [number, number]

export type Nutrition = NonNullable<components["schemas"]["Item"]["nutrition"]>

type CatalogItemSchema = components["schemas"]["CatalogItem"]

export type CatalogItem = Omit<CatalogItemSchema, "id" | "price" | "stock_amount" | "photo_url" | "nutrition"> & Required<Pick<CatalogItemSchema, "id" | "price" | "stock_amount" | "photo_url">> & {
  nutrition?: Nutrition
}

type LocationSchema = components["schemas"]["LocationResponse"]

export type Location = Omit<LocationSchema, "id" | "is_active"> & Required<Pick<LocationSchema, "id" | "is_active">>

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
