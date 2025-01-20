<script setup lang="ts">
import type { PaginatedServerResponse } from "~/models/PaginatedServerResponse";
import NumberedPaginator from "./NumberedPaginator.vue";
import type { Chapter } from "~/models/Chapter";

const { novelTitle } = useRoute().params;

defineProps<{
  errorMessage: string;
  paginatedData: PaginatedServerResponse<Chapter> | null;
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
    {{ errorMessage == "" ? "No Chapters Found" : errorMessage }}</ErrorAlert
  >

  <div v-else-if="paginatedData && paginatedData.data.length > 0">
    <div class="grid grid-cols-2 gap-10 justify-center gap-y-2">
      <NuxtLink
        v-for="chapter in paginatedData.data"
        :key="chapter.chapterNo"
        class="bg-secondary border-2 border-transparent h-full rounded-md transition-all duration-150 hover:border-accent-gold hover:scale-105 hover:brightness-105"
        :to="`/novels/${novelTitle}/${chapter.chapterNo}`"
      >
        <span class="ml-2">Chapter {{ chapter.chapterNo }}</span>
      </NuxtLink>
    </div>

    <VerticalSpacer />

    <NumberedPaginator
      :totalPages="paginatedData.totalPages"
      :total="paginatedData.total"
      @page-change="(page, limit) => onPageChange(page, limit)"
    />
  </div>
</template>
