import type { CatalogItem, Location, Nutrition } from "~/types/catalog"

export type AdminItem = CatalogItem

export type AdminLocation = Location

export type AdminCategory = Exclude<CatalogItem["category"], undefined>

export interface AdminItemForm {
  name: string
  description: string
  category: AdminCategory
  photo_url: string
  priceRub: number
  stock_amount: number
  nutrition: Nutrition
}

export interface AdminLocationForm {
  slug: string
  name: string
  address: string
  is_active: boolean
}

export interface AdminStats {
  itemsCount: number
  lowStockCount: number
  locationsCount: number
  activeLocationsCount: number
}
