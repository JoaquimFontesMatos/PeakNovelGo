export default defineNuxtRouteMiddleware((to, from) => {
    const { user } = storeToRefs(useAuthStore());

    // isAuthenticated() is an example method verifying if a user is authenticated
    if (!user.value) {
        console.log('hey');

        return navigateTo('/auth/login');
    }
});
