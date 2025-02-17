<script setup lang="ts">
import { loginFormSchema, type LoginForm } from '~/schemas/Forms';

// Use the Vee-Validate form hook
const { handleSubmit } = useForm<LoginForm>({
  validationSchema: toTypedSchema(loginFormSchema),
});

const { value: email, errorMessage: emailError } = useField('email');
const { value: password, errorMessage: passwordError } = useField('password');

// Reactive object for handling form data
const authStore = useAuthStore();
const { loadingLogin } = storeToRefs(authStore);

const showPassword = ref(false);

const toggleShowPassword = () => {
  showPassword.value = !showPassword.value;
};

const onSubmit = handleSubmit(async (values: LoginForm) => {
  try {
    await authStore.login(values);

    navigateTo('/');
  } catch (error) {}
});

const handleEnterSignUp = () => {
  navigateTo('/auth/sign-up');
};

const loginWithGoogle = () => {
  window.location.href = 'http://localhost:8081/auth/google';
};
</script>

<template>
  <main class="my-10 bg-gradient-to-r from-primary to-secondary px-5 py-2.5 md:px-20 md:py-10">
    <div class="flex flex-col justify-center gap-10 md:flex-row">
      <!-- Login Section -->
      <section class="flex w-full flex-col items-center justify-center rounded-lg bg-secondary p-4 text-secondary-content shadow-lg md:w-2/4 md:p-8">
        <h1 class="text-center text-4xl font-bold text-primary-content">Login to Your Account</h1>

        <VerticalSpacer />

        <!-- Email Input -->

        <div class="w-full md:w-2/3">
          <label for="email" class="block after:text-sm after:text-error after:content-['*']"> Email </label>
          <input id="email" name="email" type="email" v-model="email" autocomplete="email" placeholder="you@example.com" />
          <span v-if="emailError" class="mt-1 text-sm text-error">
            {{ emailError }}
          </span>
        </div>

        <SmallVerticalSpacer />

        <!-- Password Input -->
        <div class="w-full md:w-2/3">
          <label for="password" class="block after:text-sm after:text-error after:content-['*']"> Password </label>
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
            <button type="button" @click="toggleShowPassword" class="absolute right-2 top-1/2 -translate-y-1/2 transform text-gray-500 focus:outline-none">
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
        <Button :disabled="loadingLogin" @click="onSubmit">
          <div v-if="loadingLogin" class="flex items-center justify-center rounded">
            <LoadingSpinner />
            <span>Signing in...</span>
          </div>
          <span v-else>Login</span>
        </Button>

        <SmallVerticalSpacer />
      </section>

      <section class="flex w-full flex-col items-center justify-center rounded-lg bg-secondary p-4 text-secondary-content shadow-lg md:w-1/4 md:p-8">
        <h1 class="text-center text-4xl font-bold text-primary-content">Other Sign-up Options</h1>

        <VerticalSpacer />

        <Button :disabled="loadingLogin" @click="loginWithGoogle">
          <div v-if="loadingLogin" class="flex items-center justify-center rounded">
            <LoadingSpinner />
            <span>Signing in...</span>
          </div>
          <span v-else>Login with Google</span>
        </Button>
      </section>

      <!-- Sign Up Section -->
      <section class="flex w-full flex-col items-center justify-center rounded-lg bg-accent-gold p-4 text-primary shadow-lg md:w-1/4 md:p-8">
        <h1 class="text-center text-3xl font-bold">New Here?</h1>

        <VerticalSpacer />

        <p class="text-center text-lg">Join now and start reading!</p>

        <SmallVerticalSpacer />

        <!-- Sign Up Button -->
        <Button @click="handleEnterSignUp" class="w-full">Sign-up</Button>
      </section>
    </div>
  </main>
</template>
