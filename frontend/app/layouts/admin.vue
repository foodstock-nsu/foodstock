<template lang="pug">
.admin-layout.min-h-screen.surface-base(class="lg:flex")
  admin-sidebar

  main.admin-main.flex-1.min-w-0.container-pad
    .max-w-6xl.mx-auto.w-full
      slot
</template>

<style scoped>
.admin-main {
  padding-top: 1.5rem;
  padding-bottom: 2rem;
}

@media (width >= 1024px) {
  .admin-main {
    padding-top: 2rem;
  }
}
</style>
