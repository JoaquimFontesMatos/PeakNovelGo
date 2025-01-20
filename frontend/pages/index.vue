<script setup lang="ts"></script>

<template>
  <div class="w-screen">
    <h1 class="text-primary-content bg-primary">Welcome to Nuxt 3!</h1>
    <p class="border-info card">This is a card</p>
    <p class="bg-accent-gold">This is a card</p>
    <NuxtLink to="/about">
      <Button>About Us</Button>
    </NuxtLink>

    <NuxtLink to="/novels/reverend-insanity/2">
      <Button>Go to Reverend Insanity Chapter 2</Button>
    </NuxtLink>
  </div>
</template>
