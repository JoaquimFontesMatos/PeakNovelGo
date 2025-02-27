<script setup lang="ts">
import { useScroll } from '@vueuse/core';
import { BaseTextUtils } from '~/composables/textUtils';

const { novelTitle, chapterNo } = useRoute().params as { novelTitle: string; chapterNo: string };

const textUtils = new BaseTextUtils();

const chapterStore = useChapterStore();
const ttsStore = useTTSStore();
const bookmarkStore = useBookmarkStore();
const userStore = useUserStore();
const authStore = useAuthStore();
const currentChapter = ref(0);

const { chapter, fetchingChapters, paginatedChapterData } = storeToRefs(chapterStore);

const { bookmark, fetchingBookmarkedNovel, updatingBookmark } = storeToRefs(bookmarkStore);

const { user, isReaderMode } = storeToRefs(userStore);

const { currentTime, isPlaying, duration, audioPlayer } = storeToRefs(ttsStore);

const fetchBookmark = async (novelId: string): Promise<void> => {
  try {
    await useBookmarkStore().fetchBookmarkedNovelByUser(novelId);
  } catch {}
};

const goToPreviousChapter = async (): Promise<void> => {
  if (!bookmark.value) return;

  bookmark.value.currentChapter = currentChapter.value - 1;

  try {
    if (authStore.isUserLoggedIn()) {
      await bookmarkStore.updateBookmark(bookmark.value);
    }
  } catch {}

  navigateTo((('/novels/' + novelTitle) as string) + '/' + bookmark.value.currentChapter);
};

const goToNextChapter = async (): Promise<void> => {
  if (!bookmark.value) return;
  bookmark.value.currentChapter = currentChapter.value + 1;

  try {
    if (authStore.isUserLoggedIn()) {
      await bookmarkStore.updateBookmark(bookmark.value);
    }
  } catch {}

  navigateTo((('/novels/' + novelTitle) as string) + '/' + bookmark.value.currentChapter);
};

watchEffect(async () => {
  const novelUpdatesId = novelTitle as string;
  const chapterNum = parseInt(chapterNo as string);

  if (chapterNum === undefined || chapterNum === 0) {
    return;
  }

  currentChapter.value = chapterNum;

  try {
    chapterStore.fetchChapter(novelUpdatesId, currentChapter.value);
  } catch {}

  if (authStore.isUserLoggedIn()) {
    await fetchBookmark(novelTitle as string);
  }
});

/**
 * Manage the reading prefs
 */
const drawerOpen = ref<boolean>(false);

// Use the useScroll function to track scroll position reactively
const { y } = useScroll(window);

const scrollProgress = computed(() => {
  if (import.meta.client) {
    // Ensure this code runs only on the client-side
    const scrollTop = y.value || 0;
    const scrollHeight = document.documentElement.scrollHeight - window.innerHeight;
    return scrollTop / scrollHeight;
  }
  return 0; // Return a default value on the server-side
});

/**
 * TODO: add slider in the settings menu so i can increase height
 */
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
        {
          path: '/novels/' + novelTitle + '/' + chapterNo,
          name: 'Chapter ' + chapterNo as string
        },
      ]"
    />

    <VerticalSpacer />

    <LoadingBar v-show="fetchingChapters" />
    <div
      class="fixed left-0 top-0 z-50 h-[2px] w-full origin-left bg-accent-gold-dark transition-transform duration-500 md:h-1"
      :style="{
        transform: 'scaleX(' + scrollProgress + ')',
      }"
    />
    <section v-if="isReaderMode" v-show="!fetchingChapters" @click="drawerOpen = !drawerOpen">
      <div v-if="chapter">
        <div>
          <div class="mb-4 rounded-md border-[0.5px] border-accent-gold-dark bg-secondary px-4 py-2 text-secondary-content">
            <h1>Chapter {{ chapter.chapterNo }}</h1>
          </div>
          <p
            :class="user ? user.readingPreferences.font : ''"
            v-html="user && user.readingPreferences.atomicReading ? textUtils.convertLineBreaksToHtml(textUtils.toBionicText(chapter.body)) : textUtils.convertLineBreaksToHtml(chapter.body)"
          />
        </div>

        <VerticalSpacer />
        <Button @click="goToPreviousChapter" :disabled="currentChapter === 1"> << Previous Chapter </Button>

        <Button @click="goToNextChapter" :disabled="currentChapter && paginatedChapterData && currentChapter === paginatedChapterData.total">
          Next Chapter >>
        </Button>
      </div>

      <ErrorAlert v-else>Error: {{ chapter === null ? 'Invalid Chapter Number' : null }}</ErrorAlert>
    </section>
    <section v-else v-show="!fetchingChapters" @click="drawerOpen = !drawerOpen">
      <TTSReader :novel-title="novelTitle as string" :chapter="chapter" />
    </section>

    <section
      :class="drawerOpen ? 'h-1/3 md:w-1/3 md:h-full py-4' : 'h-0 w-0 m-0 shadow-none border-none backdrop-filter-none'"
      class="bg-primary-100/10 fixed bottom-0 left-0 z-50 w-full select-none overflow-y-scroll px-4 shadow-md backdrop-blur-sm transition-all md:bottom-0 md:left-2/3"
    >
      <div class="form-container flex w-full items-center justify-between rounded-lg bg-secondary p-4 shadow-lg">
        <button
          @click="goToPreviousChapter"
          :disabled="currentChapter === 1"
          class="bg-primary/80 disabled:bg-bg-primary/20 hover:bg-primary/40 flex items-center rounded-lg p-1 transition-colors disabled:cursor-not-allowed disabled:text-gray-300"
        >
          <Icon name="fluent:previous-28-filled" class="text-accent-gold-dark" />
        </button>

        <div v-if="authStore.isUserLoggedIn()" class="form-group flex items-center space-x-2">
          <input id="readerMode" name="readerMode" type="checkbox" v-model="isReaderMode" @change="isReaderMode ? (audioPlayer = null) : console.log('')" />
          <label for="readerMode" class="text-sm font-medium text-secondary-content">Reader Mode</label>
        </div>

        <button
          @click="goToNextChapter"
          :disabled="currentChapter === paginatedChapterData?.total"
          class="bg-primary/80 disabled:bg-bg-primary/20 hover:bg-primary/40 flex items-center rounded-lg p-1 transition-colors disabled:cursor-not-allowed disabled:text-gray-300"
        >
          <Icon name="fluent:next-28-filled" class="text-accent-gold-dark" />
        </button>
      </div>

      <!-- Play/Pause Button -->
      <SmallVerticalSpacer v-if="audioPlayer" />

      <div v-if="audioPlayer" class="form-container space-y-6 rounded-lg bg-secondary p-4 shadow-lg">
        <fieldset class="flex items-center gap-4 border-t border-accent-gold-dark pt-4">
          <legend class="ml-3.5 px-3.5 text-lg font-semibold text-primary-content">Audio Controls</legend>
          <!-- Atomic Reading -->
          <button @click="ttsStore.togglePlayback" class="bg-primary/80 flex items-center rounded p-2 text-white">
            <Icon v-if="isPlaying" name="fluent:pause-28-filled" class="text-accent-gold-dark" />
            <Icon v-else name="fluent:play-28-filled" class="text-accent-gold-dark" />
          </button>

          <!-- Progress Display -->
          <div class="flex-1">
            <div class="h-2 w-full rounded bg-gray-200">
              <div
                class="h-2 rounded bg-accent-gold-dark"
                :style="{
                  width: (currentTime / duration) * 100 + '%',
                }"
              />
            </div>
            <div class="mt-1 text-sm text-gray-600">{{ Math.floor(currentTime) }}s / {{ Math.floor(duration) }}s</div>
          </div>
        </fieldset>
      </div>

      <SmallVerticalSpacer />

      <ReadingPreferencesForm v-if="authStore.isUserLoggedIn()" />
    </section>
  </Container>
</template>
