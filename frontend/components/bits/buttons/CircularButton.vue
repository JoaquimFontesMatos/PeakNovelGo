<script setup lang="ts">
const props = withDefaults(
  defineProps<{
    iconName: string;
    iconSize?: number;
    padding: number;
    noBackground?: boolean;
    onClick?: () => void;
  }>(),
  {
    noBackground: false,
  }
);

const emit = defineEmits(['click']);

const handleClick = () => {
  if (props.onClick) {
    props.onClick();
  }
  emit('click');
};

const btnClass = `p-${props.padding} bg-${props.noBackground ? 'secondary' : 'primary'} rounded-full flex justify-center
items-center transition-colors hover:bg-gray-500 hover:animate-spin disabled:animate-none disabled:cursor-not-allowed
disabled:opacity-50`;
</script>

<template>
  <button @click="handleClick" :class="btnClass">
    <Icon
      :name="props.iconName"
      :size="props.iconSize"
      class="text-accent-gold-dark"
    />
  </button>
</template>
