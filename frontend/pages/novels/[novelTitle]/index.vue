<script setup lang="ts">
import LoginGuard from '~/components/LoginGuard.vue';

const { novelTitle } = useRoute().params;

function openNovelUpdatesUrl(url: string) {
  window.open(url, '_blank');
}

const novelStore = useNovelStore();
const chapterStore = useChapterStore();
const bookmarkStore = useBookmarkStore();
const authStore = useAuthStore();

const { novel, fetchingNovel } = storeToRefs(novelStore);
const { paginatedChapterData, fetchingChapters } = storeToRefs(chapterStore);

const { bookmark, fetchingBookmarkedNovel, bookmarkingNovel, updatingBookmark } = storeToRefs(bookmarkStore);

const isBookmarked = computed(() => {
  return bookmark.value;
});

const fetchBookmark = async (novelId: string): Promise<void> => {
  try {
    await useBookmarkStore().fetchBookmarkedNovelByUser(novelId);
  } catch {}
};

const bookmarkNovel = async (novelId: number): Promise<void> => {
  try {
    await useBookmarkStore().bookmarkNovel(novelId);
  } catch {}
};

const unbookmarkNovel = async (novelId: number): Promise<void> => {
  try {
    await useBookmarkStore().unbookmarkNovel(novelId);
  } catch {}
};

const updateBookmark = async (): Promise<void> => {
  try {
    if (bookmark.value) {
      if (bookmark.value.status === 'Unfollow') {
        await unbookmarkNovel(bookmark.value.ID);
        return;
      }

      await bookmarkStore.updateBookmark(bookmark.value);
    }
  } catch {}
};

const onPageChange = async (newPage: number, limit: number): Promise<void> => {
  try {
    await useChapterStore().fetchChapters(novelTitle as string, newPage, limit);
  } catch {}
};

const setScore = async (score: number) => {
  try {
    if (bookmark.value) {
      bookmark.value.score = score;
      await updateBookmark();
    }
  } catch {}
};

watchEffect(async () => {
  await Promise.all([useNovelStore().fetchNovel(novelTitle as string), onPageChange(1, 10)]);

  if (authStore.isUserLoggedIn()) {
    await fetchBookmark(novelTitle as string), onPageChange(1, 10);
  }
});
</script>

<template>
  <Container>
    <RouteTree
      :routes="[
        { path: '/', name: 'Home' },
        { path: '/novels', name: 'Novels' },
        {
          path: '/novels/' + novelTitle,
          name: novelTitle as string
        },
      ]"
    />

    <VerticalSpacer />

    <LoadingBar v-show="fetchingNovel" />

    <div v-show="!fetchingNovel">
      <ErrorAlert v-if="novel === null">Error: Novel not Found.</ErrorAlert>
      <div v-else>
        <div class="mb-4 rounded-md border-[0.5px] border-accent-gold-dark bg-secondary px-4 py-2 text-secondary-content">
          <h1>{{ novel.title }}</h1>
        </div>

        <div class="details-grid">
          <div class="figure">
            <img :src="novel.coverUrl" alt="Cover Image" class="aspect-auto h-full object-cover" />
          </div>
          <div class="general-details">
            <h1 class="text-bold text-4xl text-secondary-content">
              {{ novel.title }}
            </h1>
            <SmallVerticalSpacer />

            <DetailsLabel>Genre(s):</DetailsLabel>

            <ul class="flex max-h-28 flex-wrap overflow-scroll lg:max-h-64 xl:h-auto">
              <li v-for="({ name, id }, index) in novel.genres" :key="id" class="text-accent-gold">
                <NuxtLink :to="'/novels/genre/' + name" class="hover:underline">
                  {{ name }}
                </NuxtLink>
                <span v-if="index !== novel.genres.length - 1" class="mr-2">,</span>
              </li>
            </ul>
            <SmallVerticalSpacer />

            <DetailsLabel>Author(s):</DetailsLabel>
            <ul class="flex max-h-28 flex-wrap overflow-scroll lg:max-h-64 xl:h-auto">
              <li v-for="({ name, id }, index) in novel.authors" :key="id" class="text-accent-gold">
                <NuxtLink :to="'/novels/author/' + name" class="hover:underline">
                  {{ name }}
                </NuxtLink>
                <span v-if="index !== novel.authors.length - 1" class="mr-2">,</span>
              </li>
            </ul>
            <SmallVerticalSpacer />

            <DetailsLabel>Tag(s):</DetailsLabel>
            <ul class="flex max-h-28 flex-wrap overflow-scroll lg:max-h-64 xl:max-h-full">
              <li v-for="({ name, id }, index) in novel.tags" :key="id" class="text-accent-gold">
                <NuxtLink :to="'/novels/tag/' + name" class="hover:underline">
                  {{ name }}
                </NuxtLink>
                <span v-if="index !== novel.tags.length - 1" class="mr-2">,</span>
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
              class="w-5 hover:scale-110 hover:cursor-pointer hover:brightness-125"
              src="@img/novel_updates_logo.png"
              alt="Novel Updates Logo"
            />
          </div>
          <div class="buttons">
            <div v-if="authStore.isUserLoggedIn()">
              <LoadingBar v-show="bookmarkingNovel || fetchingBookmarkedNovel" />

              <div v-if="!bookmarkingNovel && bookmark" class="w-full justify-center">
                <div class="flex flex-wrap gap-3.5">
                  <div class="flex grow flex-col gap-3.5">
                    <label for="status" class="Status"> Status </label>
                    <select id="status" name="status" v-model="bookmark.status" @change="updateBookmark()" class="">
                      <option value="Unfollow">Unfollow</option>
                      <option value="Reading">Reading</option>
                      <option value="Completed">Completed</option>
                      <option value="On-Hold">On-Hold</option>
                      <option value="Dropped">Dropped</option>
                      <option value="Plan to Read">Plan to Read</option>
                    </select>
                  </div>

                  <div class="flex grow flex-col gap-3.5">
                    <label for="rating" class="block"> Rating </label>
                    <div class="flex flex-row gap-1">
                      <Icon
                        v-for="index in 5"
                        :key="index"
                        @click="setScore(index)"
                        :class="index <= bookmark.score ? 'text-accent-gold-dark' : ''"
                        name="fluent:star-12-filled"
                      />
                    </div>
                  </div>
                </div>
                <SmallVerticalSpacer />
              </div>
            </div>

            <VerticalSpacer />

            <div v-show="!bookmarkingNovel">
              <LoginGuard v-slot="{ handleClick }">
                <button
                  v-if="!isBookmarked && novel?.ID"
                  @click="handleClick($event, () => bookmarkNovel(novel!.ID))"
                  class="active:bg-gradient- flex w-min content-center justify-center rounded-full p-3 hover:bg-secondary hover:text-accent-gold-light hover:transition-colors active:bg-primary active:transition-colors"
                >
                  <Icon name="fluent:bookmark-16-regular" class="text-accent-gold-dark" :size="'1.5em'" />
                </button>
              </LoginGuard>
            </div>
            <LoadingIndicator v-show="fetchingChapters" />
            <section v-show="!fetchingChapters">
              <div v-if="paginatedChapterData?.data?.length">
                <section class="flex gap-4">
                  <Button v-if="paginatedChapterData && !fetchingChapters && bookmark?.currentChapter === 1" class="flex-grow">
                    <NuxtLink class="block h-full w-full text-center" :to="'/novels/' + novelTitle + '/' + 1"> Start Reading </NuxtLink>
                  </Button>
                  <Button v-if="paginatedChapterData && !fetchingChapters && bookmark && bookmark.currentChapter > 1" class="flex-grow">
                    <NuxtLink class="block h-full w-full text-center" :to="'/novels/' + novelTitle + '/' + bookmark.currentChapter">
                      Continue {{ bookmark.currentChapter }}
                    </NuxtLink>
                  </Button>
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
          :error-message="paginatedChapterData === null ? 'No Chapters Found' : null"
          :on-page-change="onPageChange"
        />
      </div>
    </div>
  </Container>
</template>

<style scoped>
.details-grid {
  display: grid;
  gap: 1.5rem;
  grid-auto-columns: 1fr;

  grid-template-areas:
    'image'
    'buttons'
    'general-info'
    'extra-info';

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
      'image image'
      'general-info general-info'
      'buttons extra-info';
  }
}

@media (min-width: 40em) {
  .details-grid {
    grid-template-areas:
      'image general-info'
      'buttons extra-info';
  }
}

@media (min-width: 50em) {
  .details-grid {
    grid-template-areas:
      'image general-info general-info extra-info'
      'image buttons buttons extra-info';
  }
}
</style>
