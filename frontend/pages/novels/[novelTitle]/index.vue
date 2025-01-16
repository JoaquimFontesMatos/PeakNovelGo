<script setup lang="ts">
import type { Novel } from "~/models/Novel";

const runtimeConfig = useRuntimeConfig();
const { novelTitle } = useRoute().params;

var errorMessage = "";

const url = runtimeConfig.public.apiUrl;

const { data, error } = await useAsyncData("novel", () =>
  fetchNovel(novelTitle as string, url)
);

if (data.value == undefined || data.value.title == undefined) {
  errorMessage = "Novel not found";
}

if (error.value) {
  errorMessage = error.value.message;
}

function fetchNovel(
  novelTitle: string,
  url: string
): Promise<Novel | undefined> {
  return fetch(
    `${url}/novels/chapters/novel/title/${encodeURIComponent(
      novelTitle as string
    )}`
  ).then((res) => res.json());
}
</script>

<template>
  <main class="px-20 py-10">
    <RouteTree
      class="mb-10"
      :routes="[
        { path: '/', name: 'Home' },
        { path: '/novels', name: 'Novels' },
        { path: `/novels/${novelTitle}`, name: novelTitle as string },
      ]"
    />

    <div v-if="errorMessage === '' && data && data.title">
      <div>
        <div
          class="mb-4 bg-secondary text-secondary-content rounded-md border-accent-gold-dark border-[0.5px] px-4 py-2"
        >
          <h1>{{ data.title }}</h1>
        </div>

        <div class="flex flex-row justify-between h-[20rem] space-x-10">
          <div class="bg-secondary h-full my-auto px-5 py-10 grow-0">
            <img
              :src="data.coverUrl"
              alt="Cover Image"
              class="aspect-auto h-full object-cover"
            />
          </div>
          <div class="bg-secondary flex flex-col justify-between grow p-5">
            <div class="flex justify-between">
              <div class="flex items-center">
                <div class="text-secondary-content">Created At:</div>
                <div class="ml-2 text-accent-gold">
                  {{ new Date(data.CreatedAt).toLocaleString() }}
                </div>
              </div>
              <div class="flex items-center">
                <div class="text-secondary-content">Updated At:</div>
                <div class="ml-2 text-accent-gold">
                  {{ new Date(data.UpdatedAt).toLocaleString() }}
                </div>
              </div>
            </div>
            <div class="flex justify-between">
              <div class="flex items-center">
                <div class="text-secondary-content">Status:</div>
                <div class="ml-2 text-accent-gold">
                  {{ data.status }}
                </div>
              </div>
              <div class="flex items-center">
                <div class="text-secondary-content">Language:</div>
                <div class="ml-2 text-accent-gold">
                  {{ data.language }}
                </div>
              </div>
            </div>
            <div class="flex justify-between">
              <div class="flex items-center">
                <div class="text-secondary-content">URL:</div>
                <a
                  :href="data.url"
                  class="ml-2 text-accent-gold hover:underline"
                >
                  {{ data.url }}
                </a>
              </div>
              <div class="flex items-center">
                <div class="text-secondary-content">Novel Updates URL:</div>
                <a
                  :href="data.novelUpdatesUrl"
                  class="ml-2 text-accent-gold hover:underline"
                >
                  {{ data.novelUpdatesUrl }}
                </a>
              </div>
            </div>
            <div class="flex">
              <span class="text-secondary-content mr-2">Genre(s):</span>
              <div class="flex flex-wrap">
                <div
                  v-for="({ name, id }, index) in data.genres"
                  :key="id"
                  class="text-accent-gold"
                >
                  <NuxtLink :to="`/novels/genre/${name}`" class="hover:underline">
                    {{ name }}
                  </NuxtLink>
                  <span v-if="index !== data.genres.length - 1" class="mr-2"
                    >,</span
                  >
                </div>
              </div>
            </div>
            <div class="flex">
              <span class="text-secondary-content mr-2">Author(s):</span>
              <div class="flex flex-wrap">
                <div
                  v-for="({ name, id }, index) in data.authors"
                  :key="id"
                  class="text-accent-gold"
                >
                  <NuxtLink :to="`/novels/author/${name}`" class="hover:underline">
                    {{ name }}
                  </NuxtLink>
                  <span v-if="index !== data.authors.length - 1" class="mr-2"
                    >,</span
                  >
                </div>
              </div>
            </div>
            <div class="flex">
              <span class="text-secondary-content mr-2">Tag(s):</span>
              <div class="flex flex-wrap">
                <div
                  v-for="({ name, id }, index) in data.tags"
                  :key="id"
                  class="text-accent-gold"
                >
                  <NuxtLink :to="`/novels/tag/${name}`" class="hover:underline">
                    {{ name }}
                  </NuxtLink>
                  <span v-if="index !== data.tags.length - 1" class="mr-2"
                    >,</span
                  >
                </div>
              </div>
            </div>
          </div>
        </div>
        <p v-html="convertLineBreaksToHtml(data.synopsis)"></p>
      </div>
    </div>

    <ErrorAlert v-else>Error: {{ errorMessage }}</ErrorAlert>
  </main>
</template>
