<script setup lang="ts">
const { novelTitle, chapterNo } = useRoute().params;

const chapterStore = useChapterStore()
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

</script>

<template>
  <main class="px-20 py-10">
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
      <TTSReader/>
    </section>

    <section :class="drawerOpen ? 'h-1/3 md:w-1/3 md:h-full py-4' : 'h-0 w-0 m-0 shadow-none border-none backdrop-filter-none'" class="w-full fixed bottom-0 left-0 md:bottom-0 md:left-2/3 shadow-md px-4 overflow-y-scroll bg-primary-100/10 backdrop-blur-sm select-none z-50 transition-all">
      <Button @click="goToPreviousChapter" :disabled="currentChapter === 1">
        <<
      </Button>

      <Button @click="goToNextChapter" :disabled="currentChapter && paginatedChapterData && currentChapter === paginatedChapterData.total ">
        >>
      </Button>

      <SmallVerticalSpacer/>

      <div class="form-group flex items-center space-x-2">
        <input
          id="readerMode"
          name="readerMode"
          type="checkbox"
          v-model="isReaderMode"
        />
        <label for="readerMode" class="text-sm font-medium text-secondary-content">Reader Mode</label>
      </div>

      <SmallVerticalSpacer/>

      <ReadingPreferencesForm/>
    </section>
  </main>
</template>
