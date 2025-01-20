<script setup lang="ts">
const isTop = ref(true);

const handleScroll = () => {
  isTop.value = window.scrollY === 0;
};

onMounted(() => {
  window.addEventListener("scroll", handleScroll);
  handleScroll();
});

onUnmounted(() => {
  window.removeEventListener("scroll", handleScroll);
});
</script>

<template>
  <header
    :class="[
      'transition-colors duration-300 fixed z-20 top-0 pl-4 border-b-2 border-accent-gold h-14 flex flex-row items-center w-full',
      isTop
        ? 'bg-transparent border-opacity-50'
        : 'bg-primary  bg-opacity-50 backdrop-blur-md',
    ]"
    class=""
  >
    <h1 class="text-2xl font-bold">PeakNovelGo</h1>

    <HorizontalSpacer />

    <div class="flex flex-row justify-between gap-5 w-full">
      <div class="flex gap-5">
        <NuxtLink class="hover:text-accent-gold hover:underline" to="/"
          >Home</NuxtLink
        >
        <NuxtLink class="hover:text-accent-gold hover:underline" to="/novels"
          >Novels</NuxtLink
        >
        <NuxtLink class="hover:text-accent-gold hover:underline" to="/settings"
          >Settings</NuxtLink
        >
      </div>

      <div class="flex gap-5 float-right pr-4">
        <NuxtLink
          class="hover:text-accent-gold hover:underline"
          to="/auth/login"
          >Login</NuxtLink
        >
        <NuxtLink
          class="hover:text-accent-gold hover:underline"
          to="/auth/sign-in"
          >Register</NuxtLink
        >
      </div>
    </div>
  </header>
</template>
