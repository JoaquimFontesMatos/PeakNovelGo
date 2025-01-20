<script setup lang="ts">
import type { Chapter } from "~/models/Chapter";
import type { Novel } from "~/models/Novel";
import type { PaginatedServerResponse } from "~/models/PaginatedServerResponse";

const runtimeConfig = useRuntimeConfig();
const { novelTitle } = useRoute().params;

const novelErrorMessage = ref("");
const chapterErrorMessage = ref("");

const url = runtimeConfig.public.apiUrl;

const { data, error } = await useAsyncData("novel", () =>
  fetchNovel(novelTitle as string, url)
);

if (data.value == undefined || data.value.title == undefined) {
  novelErrorMessage.value = "Novel not found";
}

if (error.value) {
  novelErrorMessage.value = error.value.message;
}

function fetchNovel(
  novelTitle: string,
  url: string
): Promise<Novel | undefined> {
  return fetch(
    `${url}/novels/title/${encodeURIComponent(novelTitle as string)}`
  ).then((res) => res.json());
}

function openNovelUpdatesUrl(url: string) {
  window.open(url, "_blank");
}

const paginatedChapterData = ref<PaginatedServerResponse<Chapter> | null>(null);

async function fetchChapters(
  url: string,
  page: number,
  limit: number
): Promise<PaginatedServerResponse<Chapter> | null> {
  if (data.value === null || data.value === undefined) {
    return null;
  }
  const res = await fetch(
    `${url}/novels/chapters/novel/${data.value.novelUpdatesId}/chapters?page=${page}&limit=${limit}`
  );
  return await res.json();
}

async function onPageChange(newPage: number, limit: number) {
  try {
    const response = await fetchChapters(url, newPage, limit);
    paginatedChapterData.value = response;
  } catch (err) {
    chapterErrorMessage.value = "Failed to fetch new page";
  }
}

onMounted(async () => {
  onPageChange(1, 10);
});
</script>

<template>
  <main class="px-20 py-10">
    <RouteTree
      :routes="[
        { path: '/', name: 'Home' },
        { path: '/novels', name: 'Novels' },
        { path: `/novels/${novelTitle}`, name: novelTitle as string },
      ]"
    />

    <VerticalSpacer />

    <div v-if="novelErrorMessage === '' && data && data.title">
      <div
        class="mb-4 bg-secondary text-secondary-content rounded-md border-accent-gold-dark border-[0.5px] px-4 py-2"
      >
        <h1>{{ data.title }}</h1>
      </div>

      <div class="details-grid">
        <div class="figure">
          <img
            :src="data.coverUrl"
            alt="Cover Image"
            class="aspect-auto h-full object-cover"
          />
        </div>
        <div class="general-details">
          <h1 class="text-secondary-content text-bold text-4xl">
            {{ data.title }}
          </h1>
          <SmallVerticalSpacer />

          <DetailsLabel>Genre(s):</DetailsLabel>

          <div
            class="flex flex-wrap max-h-28 overflow-scroll lg:max-h-64 xl:h-auto"
          >
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
          <SmallVerticalSpacer />

          <DetailsLabel>Author(s):</DetailsLabel>
          <div
            class="flex flex-wrap max-h-28 overflow-scroll lg:max-h-64 xl:h-auto"
          >
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
          <SmallVerticalSpacer />

          <DetailsLabel>Tag(s):</DetailsLabel>
          <div
            class="flex flex-wrap max-h-28 overflow-scroll lg:max-h-64 xl:h-auto"
          >
            <div
              v-for="({ name, id }, index) in data.tags"
              :key="id"
              class="text-accent-gold"
            >
              <NuxtLink :to="`/novels/tag/${name}`" class="hover:underline">
                {{ name }}
              </NuxtLink>
              <span v-if="index !== data.tags.length - 1" class="mr-2">,</span>
            </div>
          </div>
        </div>
        <div class="extra-details">
          <DetailsLabel>Added At:</DetailsLabel>
          <DetailsInfo>
            {{ new Date(data.CreatedAt).toLocaleString() }}
          </DetailsInfo>
          <SmallVerticalSpacer />

          <DetailsLabel>Updated At:</DetailsLabel>
          <DetailsInfo>
            {{ new Date(data.UpdatedAt).toLocaleString() }}
          </DetailsInfo>
          <SmallVerticalSpacer />

          <DetailsLabel>Released In:</DetailsLabel>
          <DetailsInfo>
            {{ data.year }}
          </DetailsInfo>
          <SmallVerticalSpacer />

          <DetailsLabel>Status:</DetailsLabel>
          <DetailsInfo>
            {{ data.status }}
          </DetailsInfo>
          <SmallVerticalSpacer />

          <DetailsLabel>Language:</DetailsLabel>
          <DetailsInfo>
            {{ data.language }}
          </DetailsInfo>
          <SmallVerticalSpacer />

          <DetailsLabel>Release Frequency:</DetailsLabel>
          <DetailsInfo>
            {{ data.releaseFrequency }}
          </DetailsInfo>
          <SmallVerticalSpacer />

          <DetailsLabel>Novel Updates URL:</DetailsLabel>
          <img
            @click="openNovelUpdatesUrl(data.novelUpdatesUrl)"
            class="w-5 hover:cursor-pointer hover:scale-110 hover:brightness-125"
            src="@img/novel_updates_logo.png"
            alt="Novel Updates Logo"
          />
        </div>
        <div class="buttons">
          <DetailsLabel>Controls:</DetailsLabel>
          <div
            v-if="paginatedChapterData && paginatedChapterData.data.length > 0"
            class="flex gap-4"
          >
            <Button class="flex-grow">
              <NuxtLink
                class="block w-full h-full text-center"
                :to="`/novels/${novelTitle}/1`"
                >First</NuxtLink
              ></Button
            >

            <Button class="flex-grow">
              <NuxtLink
                class="block w-full h-full text-center"
                :to="`/novels/${novelTitle}/${paginatedChapterData.total}`"
                >Last -> {{ paginatedChapterData.total }}</NuxtLink
              ></Button
            >
          </div>
          <div v-else>
            <ErrorAlert>Loading Chapters...</ErrorAlert>
          </div>
        </div>
      </div>
      <DetailsLabel>Description:</DetailsLabel>
      <p v-html="convertLineBreaksToHtml(data.synopsis)" />

      <VerticalSpacer />

      <DetailsLabel>Chapters:</DetailsLabel>

      <PaginatedChapterList
        :paginated-data="paginatedChapterData"
        :error-message="chapterErrorMessage"
        :on-page-change="onPageChange"
      />
    </div>

    <ErrorAlert v-else>Error: {{ novelErrorMessage }}</ErrorAlert>
  </main>
</template>

<style scoped>
.details-grid {
  display: grid;
  gap: 1.5rem;
  grid-auto-columns: 1fr;

  grid-template-areas:
    "image"
    "buttons"
    "general-info"
    "extra-info";

  padding-block: 2rem;
  padding-inline: 2rem;
  margin-inline: auto;
}

.figure {
  grid-area: image;
}

.general-details {
  grid-area: general-info;
}

.extra-details {
  grid-area: extra-info;
}

.comments {
  grid-area: comments;
}

.related {
  grid-area: related;
}

.buttons {
  grid-area: buttons;
}

@media (min-width: 30em) {
  .details-grid {
    grid-template-areas:
      "image image"
      "general-info general-info"
      "buttons extra-info";
  }
}

@media (min-width: 40em) {
  .details-grid {
    grid-template-areas:
      "image general-info"
      "buttons extra-info";
  }
}

@media (min-width: 50em) {
  .details-grid {
    grid-template-areas:
      "image general-info general-info extra-info"
      "image buttons buttons extra-info";
  }
}
</style>
