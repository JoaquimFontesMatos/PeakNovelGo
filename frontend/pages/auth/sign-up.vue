<script setup lang="ts">
import { signUpFormSchema, type SignUpForm } from '~/schemas/Forms';

useHead({
  title: "🔐 Sign Up",
});

// Use the Vee-Validate form hook
const { handleSubmit } = useForm<SignUpForm>({
  validationSchema: toTypedSchema(signUpFormSchema),
});

const { value: email, errorMessage: emailError } = useField('email');
const { value: password, errorMessage: passwordError } = useField('password');
const { value: username, errorMessage: usernameError } = useField('username');
const { value: dateOfBirth, errorMessage: dateOfBirthError } = useField('dateOfBirth');

// Reactive object for handling form data
const authStore = useAuthStore();

const { loadingSignUp, signUpMessage } = storeToRefs(authStore);

const onSubmit = handleSubmit(async (values: SignUpForm) => {
  try {
    values.dateOfBirth = formatDateToYYYYMMDD(values.dateOfBirth);

    await authStore.signUp(values);
  } catch (error) {}
});

// Helper function to format date to YYYY-MM-DD
const formatDateToYYYYMMDD = (date: string): string => {
  const parsedDate = new Date(date);
  const year = parsedDate.getFullYear();
  const month = String(parsedDate.getMonth() + 1).padStart(2, '0');
  const day = String(parsedDate.getDate()).padStart(2, '0');
  return year + '-' + month + '-' + day;
};

const handleEnterLogin = () => {
  navigateTo('/auth/login');
};

const showPassword = ref(false);

const toggleShowPassword = () => {
  showPassword.value = !showPassword.value;
};
</script>

<template>
  <main class="my-10 bg-linear-to-r from-primary to-secondary px-5 py-2.5 md:px-20 md:py-10">
    <div class="flex flex-col justify-center gap-10 md:flex-row">
      <!-- Sign-Up Section -->
      <section class="flex w-full flex-col items-center justify-center rounded-lg bg-secondary p-4 text-secondary-content shadow-lg md:w-2/3 md:p-8">
        <h1 class="text-center text-4xl font-bold text-primary-content">Sign-up to PeakNovelGo</h1>

        <VerticalSpacer />

        <!-- Username Input -->
        <div class="w-full md:w-2/3">
          <label for="username" class="block after:text-sm after:text-error after:content-['*']"> Username </label>
          <input id="username" name="username" type="text" v-model="username" placeholder="Username" />
          <p v-if="usernameError" class="mt-1 text-sm text-red-500">
            {{ usernameError }}
          </p>
        </div>

        <SmallVerticalSpacer />

        <!-- Date of Birth Input -->
        <div class="w-full md:w-2/3">
          <label for="dateOfBirth" class="block after:text-sm after:text-error after:content-['*']"> Birthdate </label>
          <input id="dateOfBirth" name="dateOfBirth" type="date" v-model="dateOfBirth" placeholder="Date of Birth" />
          <p v-if="dateOfBirthError" class="mt-1 text-sm text-red-500">
            {{ dateOfBirthError }}
          </p>
        </div>

        <SmallVerticalSpacer />

        <!-- Email Input -->
        <div class="w-full md:w-2/3">
          <label for="email" class="block after:text-sm after:text-error after:content-['*']"> Email </label>
          <input id="email" name="email" type="email" v-model="email" placeholder="you@example.com" />
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
            <input id="password" name="password" :type="showPassword ? 'text' : 'password'" v-model="password" placeholder="********" class="pr-12" />
            <!-- Toggle Button -->
            <button type="button" @click="toggleShowPassword" class="absolute right-2 top-1/2 -translate-y-1/2 transform text-gray-500 focus:outline-hidden">
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
        <Button :disabled="loadingSignUp" @click="onSubmit">
          <div v-if="loadingSignUp" class="flex items-center justify-center rounded-sm">
            <LoadingSpinner />
            <span>Signing Up...</span>
          </div>
          <span v-else>Sign-up</span>
        </Button>

        <SmallVerticalSpacer />

        <p v-if="signUpMessage && !loadingSignUp">
          {{ signUpMessage }}
        </p>
      </section>

      <!-- Login Section -->
      <section class="flex w-full flex-col items-center justify-center rounded-lg bg-accent-gold p-8 text-primary shadow-lg md:w-1/3">
        <h1 class="text-center text-3xl font-bold">Already have an account?</h1>

        <VerticalSpacer />

        <p class="text-center text-lg">Access your account here.</p>

        <SmallVerticalSpacer />

        <!-- Login Button -->
        <Button @click="handleEnterLogin" class="hover:bg-primary-dark w-full rounded-md bg-primary py-3 text-white transition-all"> Login </Button>
      </section>
    </div>
  </main>
</template>
