<script setup lang="ts">
import type { Novel } from "~/models/Novel";
import type { PaginatedServerResponse } from "~/models/PaginatedServerResponse";

const runtimeConfig = useRuntimeConfig();
const { tagName } = useRoute().params;

var errorMessage = "";

const url = runtimeConfig.public.apiUrl;

const paginatedData = ref<PaginatedServerResponse<Novel> | null>(null);

const { data, error } = await useAsyncData("novelsByTag", () =>
  fetchNovels(tagName as string, url, 1, 10)
);

if (data.value == undefined || data.value.data.length == 0) {
  errorMessage = "Novels not found";
}

if (error.value) {
  errorMessage = error.value.message;
}

function fetchNovels(
  tagName: string,
  url: string,
  page: number,
  limit: number
): Promise<PaginatedServerResponse<Novel>> {
  return fetch(
    `${url}/novels/tags/${encodeURIComponent(
      tagName as string
    )}?page=${page}&limit=${limit}`
  ).then((res) => res.json());
}

// Handle page changes
async function onPageChange(newPage: number, limit: number) {
  try {
    const response = await fetchNovels(tagName as string, url, newPage, limit);
    paginatedData.value = response;
  } catch (err) {
    errorMessage = "Failed to fetch new page";
    console.error(err);
  }
}
</script>

<template>
  <main class="px-20 py-10">
    <RouteTree
      class="mb-10"
      :routes="[
        { path: '/', name: 'Home' },
        { path: '/novels', name: 'Novels' },
        { path: `/tag/${tagName}`, name: tagName as string },
      ]"
    />

    <PaginatedNovelGallery
      :errorMessage="errorMessage"
      :paginatedData="data"
      @page-change="onPageChange"
    />
  </main>
</template>
