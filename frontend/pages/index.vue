<script setup lang="ts">
const runtimeConfig = useRuntimeConfig();
const { data, error } = await useAsyncData("novel", () => myGetFunction(2));

var chapter = data;

function myGetFunction(id: number): Promise<unknown> {
  return fetch(`http://localhost:8080/novels/${id}`).then((res) => res.json());
}
</script>

<template>
  <div>
    <h1 class="text-primary-content bg-primary">Welcome to Nuxt 3!</h1>
    <p class="border-info card">This is a card</p>
    <p class="bg-accent-gold">This is a card</p>
    <NuxtLink to="/about">
      <Button>About Us</Button>
    </NuxtLink>

    <NuxtLink to="/novels/chapters/2">
      <Button>About Us</Button>
    </NuxtLink>
    <p>{{ chapter }}</p>
  </div>
</template>
