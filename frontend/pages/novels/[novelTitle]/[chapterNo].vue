<script setup lang="ts">
import type { Chapter } from "~/models/Chapter";

const runtimeConfig = useRuntimeConfig();
const { chapterNo, novelTitle } = useRoute().params;

var errorMessage = "";

var id = parseInt(chapterNo as string);

if (id == 0 || id == undefined) {
  errorMessage = "Invalid ID";
}

const url = runtimeConfig.public.apiUrl;

const { data, error } = await useAsyncData("chapter", () =>
  fetchChapter(id, url)
);

if (data.value == undefined || data.value.chapterNo == undefined) {
  errorMessage = "Chapter not found";
}

if (error.value) {
  errorMessage = error.value.message;
}

function fetchChapter(id: number, url: string): Promise<Chapter | undefined> {
  return fetch(
    `${url}/novels/chapters/novel/${encodeURIComponent(
      novelTitle as string
    )}/chapter/${id}`
  ).then((res) => res.json());
}
</script>

<template>
  <main class="px-20 py-10">
    <RouteTree
      :routes="[
      { path: '/', name: 'Home' },
      { path: '/novels', name: 'Novels' },
      { path: `/novels/${novelTitle}`, name: novelTitle as string },
      { path: `/novels/${novelTitle}/${chapterNo}`, name: `Chapter ${ chapterNo as string}` },
    ]"
    />

    <VerticalSpacer />

    <div v-if="errorMessage === '' && data && data.body">
      <div>
        <div
          class="mb-4 bg-secondary text-secondary-content rounded-md border-accent-gold-dark border-[0.5px] px-4 py-2"
        >
          <h1>Chapter {{ data.chapterNo }}</h1>
        </div>
        <p v-html="convertLineBreaksToHtml(toBionicText(data.body))"></p>
      </div>
    </div>

    <ErrorAlert v-else>Error: {{ errorMessage }}</ErrorAlert>
  </main>
</template>
