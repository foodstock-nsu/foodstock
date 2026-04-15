import { createNumberFormatFactory } from "intl-formats"

export const locale = navigator.languages[0] || "en"

// в древних Safari лимит 20
export const createNumberFormat = createNumberFormatFactory(locale, { maximumFractionDigits: 20 })

// Дефолтный форматтер для чисел.
export const formatNumber = createNumberFormat()
