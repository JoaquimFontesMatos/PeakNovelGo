export default defineNuxtPlugin(nuxtApp => {
    const router = useRouter();

    router.afterEach((to, from) => {
        if (to.name && to.name !== 'auth-login' && to.name !== 'auth-callback' && to.name !== 'auth-sign-up') {
            // Don't store if coming from login
            sessionStorage.setItem('previousRoute', to.fullPath);
        }
    });
});
