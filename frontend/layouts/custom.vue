<script setup lang="ts">
    const authStore = useAuthStore();
    const { user } = storeToRefs(authStore);
    const colorMode = useColorMode();

    onMounted(async () => {
        await authStore.initSession();

        if (!user.value || user.value.readingPreferences.theme === undefined || !user.value.readingPreferences.theme) return;

        colorMode.preference = user.value.readingPreferences.theme;
    });
</script>
<template>
    <meta name="viewport" content="width=device-width, initial-scale=1, viewport-fit=cover" />

    <NuxtPwaAssets />
    <div>
        <div class="min-h-svh">
            <slot />
        </div>
        <Toast />
    </div>
</template>
