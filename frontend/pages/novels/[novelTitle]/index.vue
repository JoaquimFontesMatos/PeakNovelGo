<script setup lang="ts">
const { novelTitle } = useRoute().params;

function openNovelUpdatesUrl(url: string) {
  window.open(url, "_blank");
}

const { novel, novelError, fetchingNovel } = storeToRefs(useNovelStore());
const { paginatedChapterData, fetchingChapters, chapterError } = storeToRefs(
  useChapterStore()
);

const onPageChange = async (newPage: number, limit: number): Promise<void> => {
  await useChapterStore().fetchChapters(novelTitle as string, newPage, limit);
};

watch(
  () => novelTitle,
  async (newTitle) => {
    if (newTitle !== undefined) {
      await useNovelStore().fetchNovel(novelTitle as string);
      await onPageChange(1, 10);
    }
  },
  { immediate: true }
);
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

    <LoadingBar v-show="fetchingNovel" />

    <div v-show="!fetchingNovel">
      <ErrorAlert v-if="novelError !== '' && novelError !== null"
        >Error: {{ novelError }}</ErrorAlert
      >
      <div v-else-if="novel && novel !== null && novel.title">
        <div
          class="mb-4 bg-secondary text-secondary-content rounded-md border-accent-gold-dark border-[0.5px] px-4 py-2"
        >
          <h1>{{ novel.title }}</h1>
        </div>

        <div class="details-grid">
          <div class="figure">
            <img
              :src="novel.coverUrl"
              alt="Cover Image"
              class="aspect-auto h-full object-cover"
            />
          </div>
          <div class="general-details">
            <h1 class="text-secondary-content text-bold text-4xl">
              {{ novel.title }}
            </h1>
            <SmallVerticalSpacer />

            <DetailsLabel>Genre(s):</DetailsLabel>

            <ul
              class="flex flex-wrap max-h-28 overflow-scroll lg:max-h-64 xl:h-auto"
            >
              <li
                v-for="({ name, id }, index) in novel.genres"
                :key="id"
                class="text-accent-gold"
              >
                <NuxtLink :to="`/novels/genre/${name}`" class="hover:underline">
                  {{ name }}
                </NuxtLink>
                <span v-if="index !== novel.genres.length - 1" class="mr-2"
                  >,</span
                >
              </li>
            </ul>
            <SmallVerticalSpacer />

            <DetailsLabel>Author(s):</DetailsLabel>
            <ul
              class="flex flex-wrap max-h-28 overflow-scroll lg:max-h-64 xl:h-auto"
            >
              <li
                v-for="({ name, id }, index) in novel.authors"
                :key="id"
                class="text-accent-gold"
              >
                <NuxtLink
                  :to="`/novels/author/${name}`"
                  class="hover:underline"
                >
                  {{ name }}
                </NuxtLink>
                <span v-if="index !== novel.authors.length - 1" class="mr-2"
                  >,</span
                >
              </li>
            </ul>
            <SmallVerticalSpacer />

            <DetailsLabel>Tag(s):</DetailsLabel>
            <ul
              class="flex flex-wrap max-h-28 overflow-scroll lg:max-h-64 xl:max-h-full"
            >
              <li
                v-for="({ name, id }, index) in novel.tags"
                :key="id"
                class="text-accent-gold"
              >
                <NuxtLink :to="`/novels/tag/${name}`" class="hover:underline">
                  {{ name }}
                </NuxtLink>
                <span v-if="index !== novel.tags.length - 1" class="mr-2"
                  >,</span
                >
              </li>
            </ul>
          </div>
          <div class="extra-details">
            <DetailsLabel>Added At:</DetailsLabel>
            <DetailsInfo>
              {{ new Date(novel.CreatedAt).toLocaleString() }}
            </DetailsInfo>
            <SmallVerticalSpacer />

            <DetailsLabel>Updated At:</DetailsLabel>
            <DetailsInfo>
              {{ new Date(novel.UpdatedAt).toLocaleString() }}
            </DetailsInfo>
            <SmallVerticalSpacer />

            <DetailsLabel>Released In:</DetailsLabel>
            <DetailsInfo>
              {{ novel.year }}
            </DetailsInfo>
            <SmallVerticalSpacer />

            <DetailsLabel>Status:</DetailsLabel>
            <DetailsInfo>
              {{ novel.status }}
            </DetailsInfo>
            <SmallVerticalSpacer />

            <DetailsLabel>Language:</DetailsLabel>
            <DetailsInfo>
              {{ novel.language }}
            </DetailsInfo>
            <SmallVerticalSpacer />

            <DetailsLabel>Release Frequency:</DetailsLabel>
            <DetailsInfo>
              {{ novel.releaseFrequency }}
            </DetailsInfo>
            <SmallVerticalSpacer />

            <DetailsLabel>Novel Updates URL:</DetailsLabel>
            <img
              @click="openNovelUpdatesUrl(novel.novelUpdatesUrl)"
              class="w-5 hover:cursor-pointer hover:scale-110 hover:brightness-125"
              src="@img/novel_updates_logo.png"
              alt="Novel Updates Logo"
            />
          </div>
          <div class="buttons">
            <LoadingIndicator v-show="fetchingChapters" />
            <section v-show="!fetchingChapters">
              <div v-if="paginatedChapterData?.data?.length">
                <DetailsLabel>Controls:</DetailsLabel>

                <section class="flex gap-4">
                  <Button
                    v-if="paginatedChapterData && !fetchingChapters"
                    class="flex-grow"
                  >
                    <NuxtLink
                      class="block w-full h-full text-center"
                      :to="`/novels/${novelTitle}/1`"
                      >First</NuxtLink
                    >
                  </Button>

                  <Button
                    v-if="paginatedChapterData && !fetchingChapters"
                    class="flex-grow"
                  >
                    <NuxtLink
                      class="block w-full h-full text-center"
                      :to="`/novels/${novelTitle}/${paginatedChapterData!.total}`"
                      >Last -> {{ paginatedChapterData!.total }}</NuxtLink
                    ></Button
                  >
                </section>
              </div>
            </section>
          </div>
        </div>
        <DetailsLabel>Description:</DetailsLabel>
        <p v-html="convertLineBreaksToHtml(novel.synopsis)" />

        <VerticalSpacer />

        <DetailsLabel>Chapters:</DetailsLabel>

        <LoadingBar v-show="fetchingChapters" />

        <PaginatedChapterList
          v-show="!fetchingChapters"
          :paginated-data="paginatedChapterData"
          :error-message="chapterError"
          :on-page-change="onPageChange"
        />
      </div>
      <div v-else>
        <ErrorAlert
          >Error: An error occurred while fetching the novel.</ErrorAlert
        >
      </div>
    </div>
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
