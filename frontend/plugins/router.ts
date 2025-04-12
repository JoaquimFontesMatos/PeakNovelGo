export default defineNuxtPlugin(nuxtApp => {
    const router = useRouter(); // Access the router instance

    // Store the previous route before navigating to the login page
    router.afterEach((to, from) => {
        if (from.name && from.name !== 'auth-login' && from.name !== 'auth-callback') {
            // Don't store if coming from login
            sessionStorage.setItem('previousRoute', from.fullPath);
        }
    });
});
