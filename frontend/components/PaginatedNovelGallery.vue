<script setup lang="ts">
import type { Novel } from "~/models/Novel";
import type { PaginatedServerResponse } from "~/models/PaginatedServerResponse";
import NumberedPaginator from "./NumberedPaginator.vue";

defineProps<{
  errorMessage: string;
  paginatedData: PaginatedServerResponse<Novel> | null;
  onPageChange: (newPage: number, limit: number) => void;
}>();
</script>

<template>
  <ErrorAlert
    v-if="
      errorMessage !== '' || (paginatedData && paginatedData.data.length == 0)
    "
    >Error: {{ errorMessage }}</ErrorAlert
  >

  <div v-else-if="paginatedData && paginatedData.data.length > 0">
    <p>
      total {{ paginatedData.total }}, page {{ paginatedData.page }}, limit
      {{ paginatedData.limit }}
    </p>

    <div class="grid grid-cols-2 gap-10 justify-center">
      <div
        v-for="novel in paginatedData.data"
        :key="novel.ID"
        class="bg-secondary h-full rounded-md"
      >
        <div class="flex flex-row text-secondary-content w-full h-full">
          <div>
            <img
              :src="novel.coverUrl"
              alt="Cover Image"
              class="h-[9rem] w-[6rem] object-cover rounded-s-md"
            />
          </div>
          <div class="grow p-4">
            <h1>{{ novel.title }}</h1>
          </div>
        </div>
      </div>
    </div>

    <NumberedPaginator
      :totalPages="paginatedData.totalPages"
      :limit="paginatedData.limit"
      :total="paginatedData.total"
      :currentPage="paginatedData.page"
      @page-change="(page, limit) => onPageChange(page, limit)"
    />
  </div>
</template>
