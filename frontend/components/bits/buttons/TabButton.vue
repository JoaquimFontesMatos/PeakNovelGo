<script setup lang="ts">
defineProps<{
  name: string;
  currentTab: string;
}>();

const emit = defineEmits(['selectTab']);

const handleSelectTab = () => {
  emit('selectTab');
};
</script>
<template>
  <button
    @click="handleSelectTab()"
    :class="
      name === currentTab
        ? 'border-x-2 border-x-accent-gold-dark text-accent-gold-dark'
        : ''
    "
    class="flex items-center rounded-full border border-gray-500/[50%] px-2 py-1 text-primary-content transition-colors hover:scale-105 hover:border-x-2 hover:border-x-accent-gold-dark hover:bg-gray-500/50 hover:transition-colors hover:duration-150"
  >
    {{ name.toWellFormed() }}
  </button>
</template>
