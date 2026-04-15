<template lang="pug">
div(class="min-h-screen surface-base")
  header(class="sticky top-0 z-50 glass-nav container-pad py-4 flex items-center justify-between shadow-soft")
    div(class="headline-md font-bold text-primary") Foodstock

  main(class="container-pad pb-36")
    slot

  cart-bar
</template>

<style scoped>
header {
  border-bottom: 1px solid var(--ghost-border);
}
</style>
