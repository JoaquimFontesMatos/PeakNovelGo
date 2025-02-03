<script setup lang="ts">
const isTop = ref(true);
const isMenuOpen = ref(false);

const handleScroll = () => {
  isTop.value = window.scrollY === 0;
};

const toggleMenu = () => {
  isMenuOpen.value = !isMenuOpen.value;
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
  isMenuOpen.value = false;
};

// Close menu when clicking outside
const clickOutside = (event: MouseEvent) => {
  const target = event.target as HTMLElement;
  if (!target.closest('.mobile-menu') && !target.closest('.hamburger')) {
    isMenuOpen.value = false;
  }
};

watch(isMenuOpen, (val) => {
  if (val) {
    window.addEventListener('click', clickOutside);
  } else {
    window.removeEventListener('click', clickOutside);
  }
});
</script>

<template>
  <header
    :class="[
      'transition-colors duration-300 fixed z-20 top-0 pl-4 border-b-2 border-accent-gold h-14 flex items-center justify-between w-full',
      isTop
      ? 'bg-transparent border-opacity-50'
      : 'bg-primary bg-opacity-50 backdrop-blur-md',
    ]"
  >
    <h1 @click="handleClickHome" class="text-2xl font-bold hover:cursor-pointer">PeakNovelGo</h1>

    <!-- Desktop Navigation -->
    <div class="hidden md:flex flex-1 justify-end items-center gap-5 pr-4">
      <div class="flex gap-5">
        <NuxtLink class="hover:text-accent-gold hover:underline" to="/novels">Novels</NuxtLink>
        <NuxtLink v-if="user" class="hover:text-accent-gold hover:underline" to="/settings">
          Settings
        </NuxtLink>
      </div>

      <div v-if="!user" class="flex gap-5">
        <NuxtLink class="hover:text-accent-gold hover:underline" to="/auth/login">Login</NuxtLink>
        <NuxtLink class="hover:text-accent-gold hover:underline" to="/auth/sign-up">
          Register
        </NuxtLink>
      </div>
      <div v-else class="flex gap-5 items-center">
        <p class="hidden sm:block">Hello, {{ user.username }}</p>
        <p
          @click="handleLogout()"
          class="hover:text-accent-gold hover:underline cursor-pointer"
        >
          Logout
        </p>
      </div>
    </div>

    <!-- Mobile Hamburger -->
    <button
      @click="toggleMenu"
      class="md:hidden p-2 mr-4 focus:outline-none hamburger"
      aria-label="Toggle menu"
    >
      <Icon name="fluent:line-horizontal-3-20-filled" class="w-6 h-6"
            :class="isMenuOpen ? 'hidden' : 'block'"/>
      <Icon name="fluent:dismiss-20-filled" class="w-6 h-6"
            :class="isMenuOpen ? 'block' : 'hidden'"/>
    </button>
  </header>

  <!-- Mobile Menu Drawer -->
  <div
    class="md:hidden fixed top-14 right-0 h-full w-64 bg-primary bg-opacity-50 backdrop-blur-md transform transition-transform duration-300 ease-in-out mobile-menu z-50"
    :class="isMenuOpen ? 'translate-x-0' : 'translate-x-full'"
  >
    <div class="flex flex-col p-4 space-y-4">
      <NuxtLink
        @click="isMenuOpen = false"
        class="hover:text-accent-gold hover:underline"
        to="/novels"
      >
        Novels
      </NuxtLink>
      <NuxtLink
        v-if="user"
        @click="isMenuOpen = false"
        class="hover:text-accent-gold hover:underline"
        to="/settings"
      >
        Settings
      </NuxtLink>

      <template v-if="!user">
        <NuxtLink
          @click="isMenuOpen = false"
          class="hover:text-accent-gold hover:underline"
          to="/auth/login"
        >
          Login
        </NuxtLink>
        <NuxtLink
          @click="isMenuOpen = false"
          class="hover:text-accent-gold hover:underline"
          to="/auth/sign-up"
        >
          Register
        </NuxtLink>
      </template>

      <template v-else>
        <p class="pt-4">Hello, {{ user.username }}</p>
        <p
          @click="handleLogout()"
          class="hover:text-accent-gold hover:underline cursor-pointer"
        >
          Logout
        </p>
      </template>
    </div>
  </div>
</template>