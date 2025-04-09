<script setup lang="ts">
import type {PaginatedServerResponse} from "~/schemas/PaginatedServerResponse";
import NumberedPaginator from "./NumberedPaginator.vue";
import type {NovelSchema} from "~/schemas/Novel";

defineProps<{
  errorMessage: string | null;
  paginatedData: PaginatedServerResponse<typeof NovelSchema> | null;
  onPageChange: (newPage: number, limit: number) => Promise<void>;
}>();

const statusLabels: Record<string, string> = {
  "Ongoing": "blue-500",
  "Completed": "green-500",
  "On-Hold": "yellow-500",
  "Hiatus": "yellow-500",
  "Dropped": "red-500",
};

const statusBorderClass = (status: string) => {
  switch (status) {
    case 'Completed':
      return 'border-green-500'
    case 'Ongoing':
      return 'border-blue-500'
    case 'Hiatus':
    case 'On-Hold':
      return 'border-yellow-500'
    case 'Cancelled':
    case 'Dropped':
      return 'border-red-500'
    default:
      return 'border-gray-500'
  }
}

const statusClass = (status: string): string => {
  return `text-${statusLabels[status] || 'gray-500'} ${statusBorderClass(status)}`;
};
</script>

<template>
  <ErrorAlert v-if="errorMessage !== '' && errorMessage !== null"
  >Error:
    {{ errorMessage == "" ? "No Novels Found" : errorMessage }}
  </ErrorAlert
  >

  <div v-else-if="paginatedData && paginatedData.data.length > 0">
    <ul v-auto-animate
        class="@container/list grid grid-cols-1 @md/list:grid-cols-fluid gap-2 @md/list:gap-10 justify-center">
      <li
          v-for="novel in paginatedData.data"
          :key="novel.novelUpdatesId"
          class="bg-secondary h-full rounded-md border-2 border-transparent transition-all duration-150 hover:border-accent-gold hover:scale-[1.01] hover:brightness-105 hover:drop-shadow-md"
      >
        <NuxtLink :to="'/novels/' + novel.novelUpdatesId">
          <div class="@container/item flex flex-row text-secondary-content w-full h-24 @md/item:h-[9rem]">
            <img
                :src="novel.coverUrl"
                alt="Cover Image"
                class="h-full w-1/4 object-cover rounded-s-md"
            />
            <div class="float-right w-3/4 p-4">
              <p class="truncate">{{ novel.title }}</p>
              <div class="flex flex-row gap-2">
                <!-- Status Label -->
                <span
                    class="truncate rounded p-0.5 text-xs border"
                    :class="statusClass(novel.status)"
                >
                  {{ novel.status }}
                </span>
                <!-- Total Chapters -->
                <div class="flex items-center gap-1">
                  <Icon name="fluent:book-open-16-regular"/>
                  <p class="text-xs">{{ novel.latestChapter }} Chapters</p>
                </div>
              </div>
              <div>
                <!-- TODO: Add Rating Label -->
                <span>

                </span>
                <!-- TODO: Add General Info, like Rank, Views, Bookmarks and Comments -->
                <div>

                </div>
              </div>
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
