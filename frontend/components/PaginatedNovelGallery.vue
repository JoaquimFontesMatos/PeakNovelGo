<script setup lang="ts">
import type { PaginatedServerResponse } from "~/models/PaginatedServerResponse";
import NumberedPaginator from "./NumberedPaginator.vue";
import type { Novel } from "~/models/Novel";

defineProps<{
  errorMessage: string | null;
  paginatedData: PaginatedServerResponse<Novel> | null;
  onPageChange: (newPage: number, limit: number) => Promise<void>;
}>();
</script>

<template>
  <ErrorAlert v-if="errorMessage !== '' && errorMessage !== null"
  >Error:
    {{ errorMessage == "" ? "No Novels Found" : errorMessage }}</ErrorAlert
  >

  <div v-else-if="paginatedData && paginatedData.data.length > 0">
    <ul v-auto-animate class="grid grid-cols-2 gap-10 justify-center">
      <li
        v-for="novel in paginatedData.data"
        :key="novel.novelUpdatesId"
        class="bg-secondary h-full rounded-md border-2 border-transparent transition-all duration-150 hover:border-accent-gold hover:scale-[1.01] hover:brightness-105 hover:drop-shadow-md"
      >
        <NuxtLink :to="'/novels/' + novel.novelUpdatesId">
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
      </li>
    </ul>

    <VerticalSpacer/>

    <NumberedPaginator
      :totalPages="paginatedData.totalPages"
      :total="paginatedData.total"
      @page-change="(page, limit) => onPageChange(page, limit)"
    />
  </div>
  <div v-else>
    <ErrorAlert>No novels found</ErrorAlert>
  </div>
</template>
