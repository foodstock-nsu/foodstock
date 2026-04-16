import type { CatalogItem, Filters, Nutrition, Range } from "~/types/catalog"

const cloneRange = ([min, max]: Range): Range => [min, max]

const isEqualRange = (left: Range, right: Range) => {
  return left[0] === right[0] && left[1] === right[1]
}

const isValueInRange = (value: number, [min, max]: Range) => {
  return value >= min && value <= max
}

const matchesNutritionFilters = (nutrition: Nutrition, filters: Filters) => {
  return isValueInRange(nutrition.calories, filters.calories)
    && isValueInRange(nutrition.proteins, filters.proteins)
    && isValueInRange(nutrition.fats, filters.fats)
    && isValueInRange(nutrition.carbs, filters.carbs)
}

export const cloneFilters = (filters: Filters): Filters => {
  return {
    calories: cloneRange(filters.calories),
    proteins: cloneRange(filters.proteins),
    fats: cloneRange(filters.fats),
    carbs: cloneRange(filters.carbs),
  }
}

export const areFiltersActive = (filters: Filters, defaults: Filters) => {
  return !isEqualRange(filters.calories, defaults.calories)
    || !isEqualRange(filters.proteins, defaults.proteins)
    || !isEqualRange(filters.fats, defaults.fats)
    || !isEqualRange(filters.carbs, defaults.carbs)
}

export const itemMatchesCatalogFilters = (
  item: CatalogItem,
  selectedCategory: string,
  filters: Filters,
  hasActiveFilters: boolean,
) => {
  if (selectedCategory !== "Все" && item.category !== selectedCategory) {
    return false
  }

  if (!item.nutrition) {
    return !hasActiveFilters
  }

  return matchesNutritionFilters(item.nutrition, filters)
}
