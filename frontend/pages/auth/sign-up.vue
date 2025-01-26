<script setup lang="ts">
import {signUpFormSchema, type SignUpForm} from "~/veeSchemas/Forms";

// Use the Vee-Validate form hook
const {handleSubmit} = useForm<SignUpForm>({
  validationSchema: signUpFormSchema,
});

const {value: email, errorMessage: emailError} = useField("email");
const {value: password, errorMessage: passwordError} = useField("password");
const {value: username, errorMessage: usernameError} = useField("username");
const {value: dateOfBirth, errorMessage: dateOfBirthError} =
  useField("dateOfBirth");

// Reactive object for handling form data
const authStore = useAuthStore();

const onSubmit = handleSubmit(async (values: SignUpForm) => {
  values.dateOfBirth = formatDateToYYYYMMDD(values.dateOfBirth);

  await authStore.signUp(values);
});

// Helper function to format date to YYYY-MM-DD
const formatDateToYYYYMMDD = (date: string): string => {
  const parsedDate = new Date(date);
  const year = parsedDate.getFullYear();
  const month = String(parsedDate.getMonth() + 1).padStart(2, "0");
  const day = String(parsedDate.getDate()).padStart(2, "0");
  return `${year}-${month}-${day}`;
};

const handleEnterLogin = () => {
  navigateTo("/auth/login");
};

const showPassword = ref(false);

const toggleShowPassword = () => {
  showPassword.value = !showPassword.value;
};
</script>

<template>
  <main class="px-20 py-10 my-10 bg-gradient-to-r from-primary to-secondary">
    <div class="flex flex-col md:flex-row justify-center gap-10">
      <!-- Sign-Up Section -->
      <section
        class="w-full md:w-2/3 bg-secondary text-secondary-content rounded-lg shadow-lg p-8 flex flex-col items-center justify-center"
      >
        <h1 class="text-4xl font-bold text-center text-primary-content">
          Sign-Up to our Website
        </h1>

        <VerticalSpacer/>

        <!-- Username Input -->
        <div class="w-2/3">
          <label for="username" class="block  after:content-['*'] after:text-sm after:text-error">
            Username
          </label>
          <input
            id="username"
            name="username"
            type="text"
            v-model="username"
            placeholder="Username"
          />
          <p v-if="usernameError" class="text-red-500 text-sm mt-1">
            {{ usernameError }}
          </p>
        </div>

        <SmallVerticalSpacer/>

        <!-- Date of Birth Input -->
        <div class="w-2/3">
          <label for="dateOfBirth" class="block  after:content-['*'] after:text-sm after:text-error">
            Birthdate
          </label>
          <input
            id="dateOfBirth"
            name="dateOfBirth"
            type="date"
            v-model="dateOfBirth"
            placeholder="Date of Birth"
          />
          <p v-if="dateOfBirthError" class="text-red-500 text-sm mt-1">
            {{ dateOfBirthError }}
          </p>
        </div>

        <SmallVerticalSpacer/>

        <!-- Email Input -->
        <div class="w-2/3">
          <label for="email" class="block  after:content-['*'] after:text-sm after:text-error">
            Email
          </label>
          <input
            id="email"
            name="email"
            type="email"
            v-model="email"
            placeholder="you@example.com"
          />
          <span v-if="emailError" class="mt-1 text-sm text-error">
            {{ emailError }}
          </span>
        </div>

        <SmallVerticalSpacer/>

        <!-- Password Input -->
        <div class="w-2/3">
          <label for="password" class="block after:content-['*'] after:text-sm after:text-error">
            Password
          </label>
          <div class="relative">
            <!-- Password Input -->
            <input
              id="password"
              name="password"
              :type="showPassword ? 'text' : 'password'"
              v-model="password"
              placeholder="********"
              class="pr-12"
            />
            <!-- Toggle Button -->
            <button
              type="button"
              @click="toggleShowPassword"
              class="absolute right-2 top-1/2 transform -translate-y-1/2 text-gray-500 focus:outline-none"
            >
              <Icon name="fluent:eye-16-regular" :size="'1.5em'" v-if="showPassword"/> <!-- Show Icon -->
              <Icon name="fluent:eye-off-16-regular" :size="'1.5em'" v-else/> <!-- Hide Icon -->
            </button>
          </div>
          <!-- Error Message -->
          <span v-if=" passwordError" class="mt-1 text-sm text-red-500">
              {{ passwordError }}
          </span>
        </div>

        <VerticalSpacer/>

        <!-- Submit Button -->
        <Button
          @click="onSubmit"
          class="w-2/3 py-3 bg-primary text-white rounded-md hover:bg-primary-dark transition-all"
        >
          Sign-up
        </Button>
      </section>

      <!-- Login Section -->
      <section
        class="w-full md:w-1/3 bg-accent-gold text-primary rounded-lg shadow-lg p-8 flex flex-col items-center justify-center"
      >
        <h1 class="text-3xl font-bold text-center">Already have an account?</h1>

        <VerticalSpacer/>

        <p class="text-center text-lg">Access your account here.</p>

        <SmallVerticalSpacer/>

        <!-- Login Button -->
        <Button
          @click="handleEnterLogin"
          class="w-full py-3 bg-primary text-white rounded-md hover:bg-primary-dark transition-all"
        >
          Login
        </Button>
      </section>
    </div>
  </main>
</template>
