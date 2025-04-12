<script setup lang="ts">
    import { useAuthStore } from '~/stores/auth';

    const route = useRoute();
    const router = useRouter();
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
        authStore.setSession({
            user: user,
            accessToken: accessToken,
            refreshToken: refreshToken,
        });

        const previousRoute = sessionStorage.getItem('previousRoute') || '/';

        // If the previous route exists and is not the login page, navigate there
        router.push(previousRoute); // Redirect to the previous route

        // Clear the previous route after use
        sessionStorage.removeItem('previousRoute');
    });
</script>

<template>
    <Container>
        <p>Logging in...</p>
        <LoadingBar />
    </Container>
</template>
