<script setup lang="ts">
import { onMounted } from 'vue';

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
  <div>
    <AppHeader />
    <VerticalSpacer />
    <div class="min-h-svh"><slot /></div>
    <AppFooter />
    <Toast />
  </div>
</template>
