<script setup lang="ts">
defineProps<{
  name: string,
  currentTab: string
}>()

const emit = defineEmits(['selectTab']);

const handleSelectTab = () => {
  emit('selectTab');
};
</script>
<template>
  <button @click="handleSelectTab()"
          :class="name===currentTab?'text-accent-gold-dark border-x-accent-gold-dark border-x-2':''"
          class="border border-gray-500/[50%] text-primary-content rounded-full px-2 py-1 flex items-center transition-colors hover:bg-gray-500/50 hover:border-x-accent-gold-dark hover:border-x-2 hover:transition-colors hover:duration-150 hover:scale-105">
    {{ name.toWellFormed() }}
  </button>
</template>