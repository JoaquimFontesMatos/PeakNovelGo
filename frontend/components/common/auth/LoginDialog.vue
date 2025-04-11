<script setup>
import { ref } from 'vue';

const showDialog = ref(false);

defineExpose({ showDialog });
</script>

<template>
  <Teleport to="body">
    <div
      v-if="showDialog"
      class="fixed inset-0 flex items-center justify-center bg-primary/50"
    >
      <div class="w-[400px] rounded-lg bg-border p-6 shadow-lg">
        <h2 class="text-xl font-semibold">Login Required</h2>
        <p class="mt-2 text-secondary-content">
          You need to log in to continue.
        </p>
        <div class="mt-4 flex justify-end gap-2">
          <button
            @click="showDialog = false"
            class="rounded-sm bg-secondary px-4 py-2"
          >
            Cancel
          </button>
          <NuxtLink
            to="/auth/login"
            class="rounded-sm bg-accent-gold-dark px-4 py-2 text-primary"
          >
            Login
          </NuxtLink>
        </div>
      </div>
    </div>
  </Teleport>
</template>
