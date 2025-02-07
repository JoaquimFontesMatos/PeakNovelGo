<script setup lang="ts">
import type { PaginatedServerResponse } from "~/models/PaginatedServerResponse";
import NumberedPaginator from "./NumberedPaginator.vue";
import type { Chapter } from "~/models/Chapter";

const { novelTitle } = useRoute().params;

defineProps<{
  errorMessage: string | null;
  paginatedData: PaginatedServerResponse<Chapter> | null;
  onPageChange: (newPage: number, limit: number) => Promise<void>;
}>();

const [parent] = useAutoAnimate({ duration: 150 });
</script>

<template>
  <ErrorAlert
    v-if="
    !paginatedData || !paginatedData.data || paginatedData.data.length === 0
    "
  >
    Error:
    {{
    errorMessage === "" || errorMessage === null
    ? "No Chapters Found"
    : errorMessage
    }}
  </ErrorAlert>

  <div v-else-if="paginatedData?.data?.length > 0" class="min-h-72">
    <NumberedPaginator
      :totalPages="paginatedData.totalPages"
      :total="paginatedData.total"
      @page-change="(page, limit) => onPageChange(page, limit)"
    />
    <VerticalSpacer/>
    <ul ref="parent" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-10 justify-center gap-y-2">
      <li v-for="chapter in paginatedData.data" :key="chapter.chapterNo">
        <NuxtLink
          class="block bg-secondary border-2 border-transparent h-full rounded-md transition-all duration-150 hover:border-accent-gold hover:scale-[1.01] hover:brightness-105 hover:drop-shadow-md"
          :to="'/novels/' + novelTitle + '/' + chapter.chapterNo"
        >
          <span class="ml-2">Chapter {{ chapter.chapterNo }}</span>
        </NuxtLink>
      </li>
    </ul>
    <VerticalSpacer/>
    <NumberedPaginator
      v-show="paginatedData.limit > 25"
      :totalPages="paginatedData.totalPages"
      :total="paginatedData.total"
      @page-change="(page, limit) => onPageChange(page, limit)"
    />
  </div>

  <ErrorAlert v-else>Error: Failed to Fetch Chapters</ErrorAlert>
</template>
