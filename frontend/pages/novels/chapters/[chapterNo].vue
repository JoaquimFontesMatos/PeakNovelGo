<script setup lang="ts">
import type { Chapter } from "~/models/Chapter";

const runtimeConfig = useRuntimeConfig();
const { chapterNo } = useRoute().params;

var errorMessage = "";

var id = parseInt(chapterNo as string);

if (id == 0 || id == undefined) {
  errorMessage = "Invalid ID";
}

const url = runtimeConfig.public.apiUrl;

const { data, error } = await useAsyncData("chapter", () =>
  myGetFunction(id, url)
);

if (error.value) {
  errorMessage = error.value.message;
}

function myGetFunction(id: number, url: string): Promise<Chapter> {
  return fetch(`${url}/novels/chapters/chapter/${id}`).then((res) =>
    res.json()
  );
}
</script>

<template>
  <NuxtLink to="/">
    <Button>Home</Button>
  </NuxtLink>
  <main class="px-20 py-10">
    <div v-if="errorMessage == '' && data">
      <div>
        <div
          class="mb-4 bg-secondary text-secondary-content rounded-md border-accent-gold-dark border-[0.5px] px-4 py-2"
        >
          <h1>{{ data.chapter_no }} - {{ data.title }}</h1>
        </div>
        <p v-html="convertLineBreaksToHtml(toBionicText(data.body))"></p>
      </div>
    </div>

    <ErrorAlert v-else>Error: {{ errorMessage }}</ErrorAlert>
  </main>
</template>
