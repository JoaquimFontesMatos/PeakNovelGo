<script setup lang="ts">
defineProps<{
  routes: { path: string; name: string }[];
}>();
</script>

<template>
  <div class="flex items-center">
    <div v-for="({ path, name }, index) in routes" :key="path">
      <div>
        <div
          v-if="index !== routes.length - 1"
          class="text-lg md:text-xl text-secondary hover:text-accent-gold"
        >
          <NuxtLink :to="path" class="hover:underline">
            {{ name }}
          </NuxtLink>
          <span class="mx-2">/</span>
        </div>
        <div v-else class="h-full">
          <span class="md:text-md md:text-lg text-accent-gold line-clamp-1">{{ name }}</span>
        </div>
      </div>
    </div>
  </div>
</template>
