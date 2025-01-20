<script setup lang="ts">
import type { PaginatedServerResponse } from "~/models/PaginatedServerResponse";
import NumberedPaginator from "./NumberedPaginator.vue";
import type { Novel } from "~/models/Novel";

defineProps<{
  errorMessage: string;
  paginatedData: PaginatedServerResponse<Novel> | null;
  onPageChange: (newPage: number, limit: number) => void;
}>();
</script>

<template>
  <ErrorAlert
    v-if="
      errorMessage !== '' ||
      (paginatedData && paginatedData.data.length == 0) ||
      paginatedData == null
    "
    >Error:
    {{ errorMessage == "" ? "No Novels Found" : errorMessage }}</ErrorAlert
  >

  <div v-else-if="paginatedData && paginatedData.data.length > 0">
    <div class="grid grid-cols-2 gap-10 justify-center">
      <div
        v-for="novel in paginatedData.data"
        :key="novel.ID"
        class="bg-secondary h-full rounded-md"
      >
        <NuxtLink :to="`/novels/${novel.novelUpdatesId}`">
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
        </NuxtLink>
      </div>
    </div>

    <VerticalSpacer />

    <NumberedPaginator
      :totalPages="paginatedData.totalPages"
      :total="paginatedData.total"
      @page-change="(page, limit) => onPageChange(page, limit)"
    />
  </div>
</template>
