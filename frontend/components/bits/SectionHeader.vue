<script setup lang="ts">
defineProps<{
  title: string;
  isMainHeader: boolean;
}>();
</script>

<template>
  <fieldset class="border-t border-accent-gold-dark pt-4">
    <legend
      :class="isMainHeader ? 'text-xl' : 'text-lg'"
      class="ml-3.5 px-3.5 font-semibold text-primary-content md:w-2"
    >
      {{ title }}
    </legend>
    <slot />
  </fieldset>
</template>
