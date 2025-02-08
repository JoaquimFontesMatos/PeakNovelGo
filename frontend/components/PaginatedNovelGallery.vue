<script setup lang="ts">
import type { PaginatedServerResponse } from "~/models/PaginatedServerResponse";
import NumberedPaginator from "./NumberedPaginator.vue";
import type { Novel } from "~/schemas/Novel";

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
    <ul v-auto-animate class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-2 md:gap-10 justify-center">
      <li
        v-for="novel in paginatedData.data"
        :key="novel.novelUpdatesId"
        class="bg-secondary h-full rounded-md border-2 border-transparent transition-all duration-150 hover:border-accent-gold hover:scale-[1.01] hover:brightness-105 hover:drop-shadow-md"
      >
        <NuxtLink :to="'/novels/' + novel.novelUpdatesId">
          <div class="flex flex-row text-secondary-content w-full h-24 md:h-[9rem]">
            <img
              :src="novel.coverUrl"
              alt="Cover Image"
              class="h-full w-1/3 object-cover rounded-s-md"
            />
            <div class="float-right w-2/3 p-4 line-clamp-3 md:line-clamp-5">
              <p>{{ novel.title }}</p>
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
