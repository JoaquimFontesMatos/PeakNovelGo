export default defineNuxtRouteMiddleware((to, from) => {
  const authStore = useAuthStore();

  if (!authStore.isUserLoggedIn()) {
    return navigateTo('/auth/login');
  }
});
