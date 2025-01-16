<script setup lang="ts">
const props = defineProps<{
  currentPage: number;
  totalPages: number;
  limit: number;
  total: number;
  onPageChange: (newPage: number, limit: number) => void;
}>();

const selectedPageSize = ref(props.limit);

function goToNextPage() {
  if (props.currentPage < props.totalPages) {
    props.onPageChange(props.currentPage + 1, props.limit);
  }
}

function goToPreviousPage() {
  if (props.currentPage > 1) {
    props.onPageChange(props.currentPage - 1, props.limit);
  }
}

function goToSelectedPage(newPage: number, limit: number = props.limit) {
  props.onPageChange(newPage, limit);
}

function gotToFirstPage() {
  props.onPageChange(1, props.limit);
}

function gotToLastPage() {
  props.onPageChange(props.totalPages, props.limit);
}
</script>

<template>
  <nav class="flex flex-row justify-center gap-2 items-center">
    <button
      class="rounded-full p-2 py-1 border-2 border-transparent hover:enabled:text-primary hover:enabled:bg-accent-gold-light hover:enabled:border-accent-gold disabled:text-secondary-content disabled:cursor-not-allowed"
      :disabled="currentPage === totalPages"
      @click="gotToFirstPage"
    >
      <<
    </button>

    <button
      class="rounded-full p-2 py-1 border-2 border-transparent hover:enabled:text-primary hover:enabled:bg-accent-gold-light hover:enabled:border-accent-gold disabled:text-secondary-content disabled:cursor-not-allowed"
      :disabled="currentPage === 1"
      @click="goToPreviousPage"
    >
      <
    </button>

    <button
      v-for="page in totalPages"
      :key="page"
      @click="goToSelectedPage(page)"
      class="rounded-full p-2 py-1 border-2 border-transparent hover:enabled:text-primary hover:enabled:bg-accent-gold-light hover:enabled:border-accent-gold disabled:text-secondary-content disabled:cursor-not-allowed"
      :class="currentPage === page ? 'bg-accent-gold' : ''"
      :disabled="currentPage === page"
    >
      {{ page }}
    </button>

    <button
      class="rounded-full p-2 py-1 border-2 border-transparent hover:enabled:text-primary hover:enabled:bg-accent-gold-light hover:enabled:border-accent-gold disabled:text-secondary-content disabled:cursor-not-allowed"
      :disabled="currentPage === totalPages"
      @click="goToNextPage"
    >
      >
    </button>

    <button
      class="rounded-full p-2 py-1 border-2 border-transparent hover:enabled:text-primary hover:enabled:bg-accent-gold-light hover:enabled:border-accent-gold disabled:text-secondary-content disabled:cursor-not-allowed"
      :disabled="currentPage === totalPages"
      @click="gotToLastPage"
    >
      >>
    </button>

    <label for="page-size-select" class="text-secondary-content">
      Page Size:
    </label>
    <select
      id="page-size-select"
      class="rounded-full bg-secondary p-2 py-1 border-2 border-accent-gold"
      v-model="selectedPageSize"
      @change="goToSelectedPage(1, selectedPageSize)"
    >
      <option value="10">10</option>
      <option value="25">25</option>
      <option value="50">50</option>
      <option value="100">100</option>
    </select>
  </nav>
</template>
