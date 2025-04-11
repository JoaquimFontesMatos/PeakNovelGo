<script setup lang="ts">
const toastStore = useToastStore();
const { toasts } = storeToRefs(toastStore);
</script>

<template>
  <div
    class="fixed top-2 right-2 z-50 flex w-2/3 flex-col gap-2 md:top-4 md:right-4 md:w-96 md:gap-4"
  >
    <div
      v-for="toast in toasts"
      :key="toast.id"
      :class="[
        'rounded-sm p-2 shadow-sm outline outline-2 outline-offset-[-5px] outline-accent-gold-dark transition-opacity duration-300 md:p-4',
        toast.type === 'error'
          ? 'bg-red-500 text-white'
          : toast.type === 'warning'
            ? 'bg-yellow-500 text-white'
            : 'bg-green-500 text-white',
      ]"
    >
      <button
        class="float-right"
        type="button"
        @click="toastStore.removeToast(toast.id)"
      >
        <Icon
          name="fluent:dismiss-24-filled"
          class="text-primary"
          :size="'1.5em'"
        />
      </button>
      <div class="flex flex-row gap-2">
        <Icon
          v-if="toast.icon !== 'none'"
          :name="toast.icon"
          class="stroke-primary stroke-1 text-accent-gold-dark"
          :size="'2em'"
          :mode="'svg'"
        />
        <p>{{ toast.message }}</p>
      </div>
    </div>
  </div>
</template>
