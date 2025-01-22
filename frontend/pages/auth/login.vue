<script setup lang="ts">
import type { LoginForm } from "~/models/Forms";

const form = reactive<LoginForm>({
  email: "",
  password: "",
});

const handleSubmit = async () => {
  await useAuthStore().login(form);
  navigateTo("/");
};
</script>

<template>
  <main class="px-20 py-10 my-10 bg-gradient-to-r from-primary to-secondary">
    <div class="flex flex-col md:flex-row justify-center gap-10">
      <!-- Login Section -->
      <section
        class="w-full md:w-2/3 bg-secondary text-secondary-content rounded-lg shadow-lg p-8 flex flex-col items-center justify-center space-y-6"
      >
        <h1 class="text-4xl font-bold text-center text-primary-content">
          Login to Your Account
        </h1>

        <!-- Email Input -->
        <input
          type="email"
          placeholder="Email"
          v-model="form.email"
          class="w-full p-3 border-2 border-accent rounded-md focus:outline-none focus:ring-2 focus:ring-primary"
        />

        <!-- Password Input -->
        <input
          type="password"
          placeholder="Password"
          v-model="form.password"
          class="w-full p-3 border-2 border-accent rounded-md focus:outline-none focus:ring-2 focus:ring-primary"
        />

        <!-- Submit Button -->
        <Button
          @click="handleSubmit"
          class="w-full py-3 bg-primary text-white rounded-md hover:bg-primary-dark transition-all"
        >
          Login
        </Button>
      </section>

      <!-- Sign Up Section -->
      <section
        class="w-full md:w-1/3 bg-accent-gold text-primary rounded-lg shadow-lg p-8 flex flex-col items-center justify-center space-y-6"
      >
        <h1 class="text-3xl font-bold text-center">New Here?</h1>

        <p class="text-center text-lg">Join now and start reading!</p>

        <!-- Sign In Button -->
        <Button
          class="w-full py-3 bg-primary text-white rounded-md hover:bg-primary-dark transition-all"
        >
          <NuxtLink class="block w-full h-full text-center" to="/auth/sign-up"
            >Sign-up</NuxtLink
          >
        </Button>
      </section>
    </div>
  </main>
</template>
