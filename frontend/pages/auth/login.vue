<script setup lang="ts">
    import { loginFormSchema } from '~/schemas/Forms';
    import { z } from 'zod';

    const router = useRouter();

    useHead({
        title: '🔐 Login',
    });

    const runtimeConfig = useRuntimeConfig();
    const url: string = runtimeConfig.public.apiUrl;

    const email = ref('');
    const password = ref('');

    const emailError = ref('');
    const passwordError = ref('');

    // Reactive object for handling form data
    const authStore = useAuthStore();
    const { loadingLogin } = storeToRefs(authStore);

    const showPassword = ref(false);

    const toggleShowPassword = () => {
        showPassword.value = !showPassword.value;
    };

    const onSubmit = async () => {
        // reset error messages
        emailError.value = '';
        passwordError.value = '';

        const result = loginFormSchema.safeParse({ email: email.value, password: password.value });
        if (!result.success) {
            const tree = z.treeifyError(result.error);

            emailError.value = tree.properties?.email?.errors[0] || '';
            passwordError.value = tree.properties?.password?.errors[0] || '';

            return;
        }

        try {
            await authStore.login(result.data);

            const previousRoute = sessionStorage.getItem('previousRoute') || '/';

            // If the previous route exists and is not the login page, navigate there
            await router.push(previousRoute); // Redirect to the previous route

            // Clear the previous route after use
            sessionStorage.removeItem('previousRoute');
        } catch (error) {}
    };

    const handleEnterSignUp = () => {
        navigateTo('/auth/sign-up');
    };

    const loginWithGoogle = () => {
        window.location.href = url + '/auth/oauth2/google';
    };
</script>

<template>
    <main class="my-10 bg-linear-to-r from-primary to-secondary px-5 py-2.5 md:px-20 md:py-10">
        <div class="flex flex-col justify-center gap-10 md:flex-row">
            <!-- Login Section -->
            <section class="flex w-full flex-col items-center justify-center rounded-lg bg-secondary p-4 text-secondary-content shadow-lg md:w-2/4 md:p-8">
                <h1 class="text-center text-4xl font-bold text-primary-content">Login to Your Account</h1>

                <VerticalSpacer />

                <!-- Email Input -->

                <div class="w-full md:w-2/3">
                    <label for="email" class="block after:text-sm after:text-error after:content-['*']">Email</label>
                    <input id="email" name="email" type="email" v-model="email" autocomplete="email" placeholder="you@example.com" />
                    <span v-if="emailError" class="mt-1 text-sm text-error">
                        {{ emailError }}
                    </span>
                </div>

                <SmallVerticalSpacer />

                <!-- Password Input -->
                <div class="w-full md:w-2/3">
                    <label for="password" class="block after:text-sm after:text-error after:content-['*']">Password</label>
                    <div class="relative">
                        <!-- Password Input -->
                        <input
                            id="password"
                            name="password"
                            :type="showPassword ? 'text' : 'password'"
                            v-model="password"
                            autocomplete="current-password"
                            placeholder="********"
                            class="pr-12"
                        />
                        <!-- Toggle Button -->
                        <button
                            type="button"
                            @click="toggleShowPassword"
                            class="absolute top-1/2 right-2 -translate-y-1/2 transform text-gray-500 focus:outline-hidden"
                        >
                            <Icon name="fluent:eye-16-regular" :size="'1.5em'" v-if="showPassword" />
                            <!-- Show Icon -->
                            <Icon name="fluent:eye-off-16-regular" :size="'1.5em'" v-else />
                            <!-- Hide Icon -->
                        </button>
                    </div>
                    <!-- Error Message -->
                    <span v-if="passwordError" class="mt-1 text-sm text-red-500">
                        {{ passwordError }}
                    </span>
                </div>

                <VerticalSpacer />

                <!-- Submit Button -->
                <MainButton :disabled="loadingLogin" @click="onSubmit">
                    <div v-if="loadingLogin" class="flex items-center justify-center rounded-sm">
                        <LoadingSpinner />
                        <span>Signing in...</span>
                    </div>
                    <span v-else>Login</span>
                </MainButton>

                <SmallVerticalSpacer />
            </section>

            <section class="flex w-full flex-col items-center justify-center rounded-lg bg-secondary p-4 text-secondary-content shadow-lg md:w-1/4 md:p-8">
                <h1 class="text-center text-4xl font-bold text-primary-content">Other Sign-up Options</h1>

                <VerticalSpacer />

                <MainButton :disabled="loadingLogin" @click="loginWithGoogle">
                    <div v-if="loadingLogin" class="flex items-center justify-center rounded-sm">
                        <LoadingSpinner />
                        <span>Signing in...</span>
                    </div>
                    <span v-else>Login with Google</span>
                </MainButton>
            </section>

            <!-- Sign Up Section -->
            <section class="flex w-full flex-col items-center justify-center rounded-lg bg-accent-gold p-4 text-primary shadow-lg md:w-1/4 md:p-8">
                <h1 class="text-center text-3xl font-bold">New Here?</h1>

                <VerticalSpacer />

                <p class="text-center text-lg">Join now and start reading!</p>

                <SmallVerticalSpacer />

                <!-- Sign Up Button -->
                <MainButton @click="handleEnterSignUp" class="w-full">Sign-up</MainButton>
            </section>
        </div>
    </main>
</template>
