export default defineAppConfig({
  ui: {
    // Твой Drawer (не трогаем)
    drawer: {
      slots: {
        content: "fixed bg-[var(--color-surface-container-lowest)] ring ring-[var(--ghost-border)] flex focus:outline-none",
      },
    },
    // Твой Slider (не трогаем)
    slider: {
      slots: {
        root: "relative flex items-center select-none touch-none w-full",
        track: "relative grow h-1.5 rounded-full bg-[var(--color-surface-container-high)]",
        range: "absolute h-full rounded-full bg-linear-to-r from-[var(--color-primary)] to-[var(--color-primary-container)]",
        thumb: "block w-5 h-5 bg-white border-2 border-[var(--color-primary)] rounded-full shadow-[var(--shadow-soft)] focus:outline-none focus:ring-2 focus:ring-[var(--color-primary)] focus:ring-offset-2 transition-colors cursor-pointer",
      },
    },

    // --- КНОПКИ ---
    button: {
      // Только базовые классы и скругление. Цвета берем из Nuxt UI
      slots: {
        base: "rounded-full font-bold transition-all duration-150 justify-center items-center inline-flex disabled:opacity-50",
      },
      // Удаляем цвета. Nuxt UI сам подставит primary, error, neutral
      // Оставляем только специфику для variant (если нужно)
      variants: {
        variant: {
          solid: "", // По умолчанию primary solid — это чистый primary цвет
          outline: {
            // Outline тоже будет в primary, если не задано иное
            base: "hover:bg-primary-50 dark:hover:bg-primary-900/10",
          },
          ghost: {
            base: "hover:bg-primary-50 dark:hover:bg-primary-900/10 px-4",
          },
        },
        // Убираем ручное задание размеров, используем дефолтные Nuxt UI
      },
    },

    // --- ПОЛЯ ВВОДА ---
    // Мы задаем ТОЛЬКО скругление и отступы. Фокус и цвета делает Nuxt UI.
    input: {
      slots: {
        root: "relative w-full",
        // Убрали bg-[var(...)], Nuxt UI поставит нейтральный фон
        base: "w-full border rounded-full px-4 py-3 placeholder:text-gray-400 dark:placeholder:text-gray-500 transition-all",
      },
      // Добавляем варианты для фокуса
      defaultVariants: {
        variant: "outline", // Обычный outline input
        color: "neutral", // Нейтральная обводка по умолчанию
      },
    },

    select: {
      slots: {
        base: "w-full border rounded-full px-4 py-3 appearance-none transition-all",
      },
      defaultVariants: {
        variant: "outline",
        color: "neutral",
      },
    },

    textarea: {
      slots: {
        base: "w-full border rounded-2xl px-4 py-3 transition-all",
      },
      defaultVariants: {
        variant: "outline",
        color: "neutral",
      },
    },

    // --- БЕЙДЖИ ---
    badge: {
      slots: {
        // Делаем их аккуратными и маленькими
        base: "rounded-full font-bold px-2.5 py-1 text-[10px] uppercase tracking-widest inline-flex items-center justify-center",
      },
      variants: {
        color: {
          primary: {
            base: "bg-primary-100 text-primary-700 dark:bg-primary-900/30 dark:text-primary-400",
          },
          neutral: {
            base: "bg-gray-100 text-gray-700 dark:bg-gray-800 dark:text-gray-300",
          },
          error: {
            base: "bg-red-100 text-red-700 dark:bg-red-900/30 dark:text-red-400",
          },
        },
      },
      defaultVariants: {
        color: "primary", // Бейдж АКТИВНА по умолчанию primary
      },
    },
  },
})
