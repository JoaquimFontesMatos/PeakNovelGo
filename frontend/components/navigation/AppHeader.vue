<script setup lang="ts">
    import { hasPermission } from '~/config/permissionsConfig';
    const router = useRouter();

    const isTop = ref(true);
    const isMenuOpen = ref(false);
    const lastScrollY = ref(0);
    const isScrollingUp = ref(false);
    const isHeaderVisible = ref(true);

    const handleScroll = () => {
        const currentScrollY = window.scrollY;

        // Check if we're at the top
        isTop.value = currentScrollY === 0;

        // Determine scroll direction
        isScrollingUp.value = currentScrollY < lastScrollY.value;

        // Always show header at top or when scrolling up
        isHeaderVisible.value = isTop.value || isScrollingUp.value;

        // Hide header only when scrolling down past a certain threshold (e.g., 100px)
        if (!isTop.value && currentScrollY > 100 && !isScrollingUp.value) {
            isHeaderVisible.value = false;
        }

        lastScrollY.value = currentScrollY;
    };

    const toggleMenu = () => {
        isMenuOpen.value = !isMenuOpen.value;
    };

    onMounted(() => {
        window.addEventListener('scroll', handleScroll);
        handleScroll();
    });

    onUnmounted(() => {
        window.removeEventListener('scroll', handleScroll);
    });

    const { user } = storeToRefs(useAuthStore());

    const handleLogout = async () => {
        await useAuthStore().logout();

        const previousRoute = sessionStorage.getItem('previousRoute') || '/';

        // If the previous route exists and is not the login page, navigate there
        await router.push(previousRoute); // Redirect to the previous route

        // Clear the previous route after use
        sessionStorage.removeItem('previousRoute');
    };

    const handleClickHome = () => {
        navigateTo('/');
        isMenuOpen.value = false;
    };

    // Close menu when clicking outside
    const clickOutside = (event: MouseEvent) => {
        const target = event.target as HTMLElement;
        if (!target.closest('.mobile-menu') && !target.closest('.hamburger')) {
            isMenuOpen.value = false;
        }
    };

    watch(isMenuOpen, val => {
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
            'fixed top-0 z-20 flex h-14 w-full items-center justify-between border-b-2 border-accent-gold pl-4 transition-all duration-300',
            isTop ? 'border-opacity-50 bg-transparent' : 'bg-opacity-50 bg-primary backdrop-blur-md',
            isHeaderVisible || isMenuOpen ? 'pointer-events-auto translate-y-0 opacity-100' : 'pointer-events-none -translate-y-2 opacity-0',
        ]"
        style="will-change: transform, opacity"
    >
        <div @click="handleClickHome()" class="flex cursor-pointer items-center gap-2">
            <img src="/android-chrome-512x512.png" alt="PeakNovelGo Logo" class="h-10 w-10 cursor-pointer brightness-105 hover:scale-105" />
            <h1 class="hidden text-2xl font-bold hover:cursor-pointer sm:block">PeakNovelGo</h1>
        </div>

        <!-- Desktop Navigation -->
        <div class="hidden flex-1 items-center justify-end gap-5 pr-4 md:flex">
            <div class="flex gap-5">
                <NuxtLink class="hover:text-accent-gold hover:underline" to="/novels">Novels</NuxtLink>
                <NuxtLink v-if="user" class="hover:text-accent-gold hover:underline" to="/settings">Settings</NuxtLink>
                <NuxtLink
                    v-if="user && hasPermission(user, 'novels', 'create')"
                    @click="isMenuOpen = false"
                    class="hover:text-accent-gold hover:underline"
                    to="/novels/import"
                >
                    Import Novel
                </NuxtLink>
                <NuxtLink
                    v-if="user && hasPermission(user, 'chapters', 'create')"
                    @click="isMenuOpen = false"
                    class="hover:text-accent-gold hover:underline"
                    to="/novels/import/chapters"
                >
                    Import Chapters
                </NuxtLink>
            </div>

            <div v-if="!user" class="flex gap-5">
                <NuxtLink class="hover:text-accent-gold hover:underline" to="/auth/login">Login</NuxtLink>
                <NuxtLink class="hover:text-accent-gold hover:underline" to="/auth/sign-up">Register</NuxtLink>
            </div>
            <div v-else class="flex items-center gap-5">
                <p class="hidden sm:block">Hello, {{ user.username }}</p>
                <p @click="handleLogout()" class="cursor-pointer hover:text-accent-gold hover:underline">Logout</p>
            </div>
        </div>

        <!-- Mobile Hamburger -->
        <button @click="toggleMenu()" class="hamburger mr-4 p-2 focus:outline-hidden md:hidden" aria-label="Toggle menu">
            <Icon name="fluent:line-horizontal-3-20-filled" class="h-6 w-6" v-show="!isMenuOpen" />
            <Icon name="fluent:dismiss-20-filled" class="h-6 w-6" v-show="isMenuOpen" />
        </button>
    </header>

    <!-- Mobile Menu Drawer -->
    <div
        class="mobile-menu bg-opacity-50 fixed top-14 right-0 z-50 h-full w-64 transform bg-primary backdrop-blur-md transition-transform duration-300 ease-in-out md:hidden"
        :class="isMenuOpen ? 'translate-x-0' : 'translate-x-full'"
    >
        <div class="flex flex-col space-y-4 p-4">
            <NuxtLink @click="isMenuOpen = false" class="hover:text-accent-gold hover:underline" to="/novels">Novels</NuxtLink>
            <NuxtLink v-if="user" @click="isMenuOpen = false" class="hover:text-accent-gold hover:underline" to="/settings">Settings</NuxtLink>
            <NuxtLink
                v-if="user && hasPermission(user, 'novels', 'create')"
                @click="isMenuOpen = false"
                class="hover:text-accent-gold hover:underline"
                to="/novels/import"
            >
                Import Novel
            </NuxtLink>
            <NuxtLink
                v-if="user && hasPermission(user, 'chapters', 'create')"
                @click="isMenuOpen = false"
                class="hover:text-accent-gold hover:underline"
                to="/novels/import/chapters"
            >
                Import Chapters
            </NuxtLink>

            <template v-if="!user">
                <NuxtLink @click="isMenuOpen = false" class="hover:text-accent-gold hover:underline" to="/auth/login">Login</NuxtLink>
                <NuxtLink @click="isMenuOpen = false" class="hover:text-accent-gold hover:underline" to="/auth/sign-up">Register</NuxtLink>
            </template>

            <template v-else>
                <p class="pt-4">Hello, {{ user.username }}</p>
                <p @click="handleLogout()" class="cursor-pointer hover:text-accent-gold hover:underline">Logout</p>
            </template>
        </div>
    </div>
</template>
