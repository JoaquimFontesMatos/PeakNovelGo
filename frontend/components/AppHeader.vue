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

const { user } = storeToRefs(useAuthStore());

const handleLogout = async() => {
  await useAuthStore().logout();
  navigateTo("/");
};

const handleClickHome = () => {
  navigateTo("/");
};
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
    <h1 @click="handleClickHome" class="text-2xl font-bold hover:cursor-pointer">PeakNovelGo</h1>

    <HorizontalSpacer/>

    <div class="flex flex-row justify-between gap-5 w-full">
      <div class="flex gap-5">
        <NuxtLink class="hover:text-accent-gold hover:underline" to="/novels"
        >Novels</NuxtLink
        >
        <NuxtLink v-if="user" class="hover:text-accent-gold hover:underline" to="/settings"
        >Settings</NuxtLink
        >
      </div>

      <div
        v-if="user === undefined || user === null"
        class="flex gap-5 float-right pr-4"
      >
        <NuxtLink
          class="hover:text-accent-gold hover:underline"
          to="/auth/login"
        >Login</NuxtLink
        >
        <NuxtLink
          class="hover:text-accent-gold hover:underline"
          to="/auth/sign-up"
        >Register</NuxtLink
        >
      </div>
      <div v-else class="flex gap-5 float-right pr-4">
        <p>Hello, {{ user.username }}</p>
        <p
          @click="handleLogout()"
          class="hover:text-accent-gold hover:underline cursor-pointer"
        >
          Logout
        </p>
      </div>
    </div>
  </header>
</template>
