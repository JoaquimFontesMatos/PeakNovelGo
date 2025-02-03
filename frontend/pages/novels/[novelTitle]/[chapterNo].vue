<script setup lang="ts">
import { useScroll } from '@vueuse/core';

const { novelTitle, chapterNo } = useRoute().params;

const chapterStore = useChapterStore()
const ttsStore = useTTSStore()
const bookmarkStore = useBookmarkStore()
const userStore = useUserStore()
const currentChapter = ref(0)

const {
  chapter, fetchingChapters, chapterError, paginatedChapterData
} = storeToRefs(chapterStore)

const {
  bookmark, fetchingBookmarkedNovel, bookmarkedNovelError, updateBookmarkError, updatingBookmark
} = storeToRefs(bookmarkStore)

const {user, isReaderMode} = storeToRefs(userStore)

const {
  currentTime, isPlaying, duration, audioPlayer
} = storeToRefs(ttsStore)

const fetchBookmark = async(novelId: string): Promise<void> => {
  await useBookmarkStore().fetchBookmarkedNovelByUser(novelId)
}

const goToPreviousChapter = async(): Promise<void> => {
  if (!bookmark.value)return ;

  console.log(currentChapter.value + "->" + bookmark.value.currentChapter)

  bookmark.value.currentChapter = currentChapter.value - 1

  await bookmarkStore.updateBookmark(bookmark.value)

  navigateTo('/novels/' + novelTitle as string + '/' + (bookmark.value.currentChapter))
}

const goToNextChapter = async(): Promise<void> => {
  if (!bookmark.value)return ;
  bookmark.value.currentChapter = currentChapter.value + 1

  await bookmarkStore.updateBookmark(bookmark.value)

  navigateTo('/novels/' + novelTitle as string + '/' + (bookmark.value.currentChapter))
}

watchEffect(async() => {
  const novelUpdatesId = novelTitle as string
  const chapterNum = parseInt(chapterNo as string)

  if (chapterNum === undefined || chapterNum === 0) {
    chapterError.value = "Invalid chapter number"
    return ;
  }

  currentChapter.value = chapterNum

  await chapterStore.fetchChapter(novelUpdatesId, currentChapter.value)
})

watchEffect(async() => {
    await Promise.all([
      fetchBookmark(novelTitle as string),
    ]);
  }
);

/**
 * Manage the reading prefs
 */
const drawerOpen = ref<boolean>(false)

// Use the useScroll function to track scroll position reactively
const { y } = useScroll(window)

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

    <VerticalSpacer/>

    <LoadingBar v-show="fetchingChapters"/>
    <div
      class="h-[2px] md:h-1 fixed top-0 left-0 w-full z-50 bg-accent-gold-dark origin-left transition-transform duration-500"
      :style=" {
        transform: 'scaleX(' + scrollProgress + ')'
      }"
    />
    <section v-if="isReaderMode" v-show="!fetchingChapters" @click="drawerOpen = !drawerOpen">
      <div v-if=" chapter">
        <div>
          <div
            class="mb-4 bg-secondary text-secondary-content rounded-md border-accent-gold-dark border-[0.5px] px-4 py-2"
          >
            <h1>Chapter {{ chapter.chapterNo }}</h1>
          </div>
          <p :class="user ? user.readingPreferences.font : ''"
             v-html="user && user.readingPreferences.atomicReading ? convertLineBreaksToHtml(toBionicText(chapter.body)) : convertLineBreaksToHtml(chapter.body)"/>
        </div>

        <VerticalSpacer/>
        <Button @click="goToPreviousChapter" :disabled="currentChapter === 1">
          << Previous Chapter
        </Button>

        <Button @click="goToNextChapter" :disabled="currentChapter && paginatedChapterData && currentChapter === paginatedChapterData.total ">
          Next Chapter >>
        </Button>
      </div>

      <ErrorAlert v-else>Error: {{ chapterError }}</ErrorAlert>
    </section>
    <section v-else v-show="!fetchingChapters" @click="drawerOpen = !drawerOpen">
      <TTSReader :novel-title="novelTitle as string" :chapter="chapter"/>
    </section>

    <section :class="drawerOpen ? 'h-1/3 md:w-1/3 md:h-full py-4' : 'h-0 w-0 m-0 shadow-none border-none backdrop-filter-none'"
             class="w-full fixed bottom-0 left-0 md:bottom-0 md:left-2/3 shadow-md px-4 overflow-y-scroll bg-primary-100/10 backdrop-blur-sm select-none z-50 transition-all">
      <div class="form-container p-4 bg-secondary rounded-lg shadow-lg w-full flex items-center justify-between">
        <button @click="goToPreviousChapter" :disabled="currentChapter === 1" class="flex items-center transition-colors bg-primary/80 rounded-lg p-1 disabled:cursor-not-allowed disabled:bg-bg-primary/20 hover:bg-primary/40 disabled:text-gray-300">
          <Icon name="fluent:previous-28-filled" class="text-accent-gold-dark "/>
        </button>

        <div class="form-group flex items-center space-x-2">
          <input
            id="readerMode"
            name="readerMode"
            type="checkbox"
            v-model="isReaderMode"
            @change="isReaderMode ? audioPlayer = null : console.log('')"
          />
          <label for="readerMode" class="text-sm font-medium text-secondary-content">Reader Mode</label>
        </div>

        <button @click="goToNextChapter" :disabled="currentChapter === paginatedChapterData?.total" class="flex items-center transition-colors bg-primary/80 rounded-lg p-1 disabled:cursor-not-allowed disabled:bg-bg-primary/20 hover:bg-primary/40 disabled:text-gray-300">
          <Icon name="fluent:next-28-filled" class="text-accent-gold-dark"/>
        </button>
      </div>

      <!-- Play/Pause Button -->
      <SmallVerticalSpacer v-if="audioPlayer"/>

      <div v-if="audioPlayer" class="form-container space-y-6 p-4 bg-secondary rounded-lg shadow-lg">
        <fieldset class="border-t border-accent-gold-dark pt-4 flex items-center gap-4">
          <legend class="text-lg font-semibold text-primary-content ml-3.5 px-3.5">Audio Controls</legend>
          <!-- Atomic Reading -->
          <button @click="ttsStore.togglePlayback" class="p-2 bg-primary/80 text-white rounded flex items-center">
            <Icon v-if="isPlaying" name="fluent:pause-28-filled" class="text-accent-gold-dark"/>
            <Icon v-else name="fluent:play-28-filled" class="text-accent-gold-dark"/>
          </button>

          <!-- Progress Display -->
          <div class="flex-1">
            <div class="w-full bg-gray-200 rounded h-2">
              <div
                class="bg-accent-gold-dark h-2 rounded"
                :style=" {
                  width: (currentTime / duration) * 100 + '%'
                }"
              />
            </div>
            <div class="text-sm text-gray-600 mt-1">
              {{ Math.floor(currentTime) }}s / {{ Math.floor(duration) }}s
            </div>
          </div>
        </fieldset>
      </div>

      <SmallVerticalSpacer/>

      <ReadingPreferencesForm/>
    </section>
  </Container>
</template>
