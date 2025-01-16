<script setup lang="ts">
import type { Novel } from "~/models/Novel";
import type { PaginatedServerResponse } from "~/models/PaginatedServerResponse";

const runtimeConfig = useRuntimeConfig();

var errorMessage = "";
const url = runtimeConfig.public.apiUrl;

// Reactive variable for holding the data
const paginatedData = ref<PaginatedServerResponse<Novel> | null>(null);

// Fetch data on initial load
const { data, error } = await useAsyncData("novels", () =>
  fetchNovels(url, 1, 10)
);
paginatedData.value = data.value;

if (paginatedData.value == undefined || paginatedData.value.data.length == 0) {
  errorMessage = "Novels not found";
}

if (error.value) {
  errorMessage = error.value.message;
}

// Fetch novels function with support for pagination
function fetchNovels(
  url: string,
  page: number,
  limit: number
): Promise<PaginatedServerResponse<Novel>> {
  return fetch(`${url}/novels?page=${page}&limit=${limit}`).then((res) =>
    res.json()
  );
}

// Handle page changes
async function onPageChange(newPage: number, limit: number) {
  try {
    const response = await fetchNovels(url, newPage, limit);
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
      ]"
    />

    <PaginatedNovelGallery
      :errorMessage="errorMessage"
      :paginatedData="paginatedData"
      @page-change="onPageChange"
    />
  </main>
</template>
