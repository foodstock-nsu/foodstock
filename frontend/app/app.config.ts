export default defineAppConfig({
  ui: {
    drawer: {
      slots: {
        // Перебиваем только фон — остальные классы унаследуются из defaults
        content: "fixed bg-[var(--color-surface-container-lowest)] ring ring-[var(--ghost-border)] flex focus:outline-none",
      },
    },
    slider: {
      slots: {
        root: "relative flex items-center select-none touch-none w-full",
        track: "relative grow h-1.5 rounded-full bg-[var(--color-surface-container-high)]",
        range: "absolute h-full rounded-full bg-linear-to-r from-[var(--color-primary)] to-[var(--color-primary-container)]",
        thumb: "block w-5 h-5 bg-white border-2 border-[var(--color-primary)] rounded-full shadow-[var(--shadow-soft)] focus:outline-none focus:ring-2 focus:ring-[var(--color-primary)] focus:ring-offset-2 transition-colors cursor-pointer",
      },
    },
  },
})
