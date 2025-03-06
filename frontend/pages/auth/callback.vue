<script setup lang="ts">
import { useRoute } from 'vue-router'; // or '#app' for Nuxt 3
import { useAuthStore } from '~/stores/auth';

const route = useRoute();
const authStore = useAuthStore();
const toastStore = useToastStore();

onMounted(() => {
  // Extract query parameters
  const accessToken = route.query.accessToken as string;
  const refreshToken = route.query.refreshToken as string;
  const user = route.query.user ? JSON.parse(decodeURIComponent(route.query.user as string)) : null;

  if (!accessToken || !refreshToken || !user) {
    toastStore.addToast('Error: Failed to log in', 'error');
    navigateTo('/auth/login');

    return;
  }

  // Store the session data
  authStore.setSession({ user: user, accessToken: accessToken, refreshToken: refreshToken });

  // Optionally, redirect the user to another page
  navigateTo('/'); // Or another route
});
</script>

<template>
  <Container>
    <p>Logging in...</p>
    <LoadingBar />
  </Container>
</template>
