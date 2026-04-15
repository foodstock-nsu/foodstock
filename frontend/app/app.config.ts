export default defineAppConfig({
  ui: {
    drawer: {
      slots: {
        // Перебиваем только фон — остальные классы унаследуются из defaults
        content: "fixed bg-[var(--color-surface-container-lowest)] ring ring-[var(--ghost-border)] flex focus:outline-none",
      },
    },
  },
})
