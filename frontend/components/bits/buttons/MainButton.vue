<template>
  <button
    class="relative z-0 rounded-md border-[0.5px] border-accent-gold-dark bg-secondary px-6 py-3 text-lg font-bold text-secondary-content disabled:bg-linear-to-r disabled:from-accent-gold-dark disabled:to-accent-gold-light disabled:text-primary disabled:opacity-50 disabled:after:content-none! disabled:hover:cursor-not-allowed disabled:hover:no-underline"
  >
    <slot />
  </button>
</template>

<style scoped>
@reference "../assets/css/tailwind.css";
button:after {
  content: '';
  @apply absolute top-0 left-0 h-full w-full;
  background-image: url('/img/cloud_02.svg');
  background-size: auto 100%;
  background-repeat: space;
  transition: all 300ms ease-in;
  z-index: -1;
  width: 0;
}

button:hover:after {
  @apply text-accent-gold-dark;
  width: 100%;
}

button:hover {
  @apply bg-linear-to-r from-accent-gold-dark to-accent-gold-light text-gray-900;
  transition: all 300ms ease-in;
  text-decoration: underline;
}
</style>
