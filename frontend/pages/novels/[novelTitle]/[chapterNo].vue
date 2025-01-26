<script setup lang="ts">
const { novelTitle, chapterNo } = useRoute().params;

const chapterStore = useChapterStore()
const bookmarkStore = useBookmarkStore()
const currentChapter = ref(0)

const {
  chapter, fetchingChapters, chapterError, paginatedChapterData
} = storeToRefs(chapterStore)

const {
  bookmark, fetchingBookmarkedNovel, bookmarkedNovelError, updateBookmarkError, updatingBookmark
} = storeToRefs(bookmarkStore)

const unbookmarkNovel = async(novelId: number): Promise<void> => {
  await useBookmarkStore().unbookmarkNovel(novelId)
}

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
    <section v-show="!fetchingChapters">
      <div v-if=" chapter">
        <div>
          <div
            class="mb-4 bg-secondary text-secondary-content rounded-md border-accent-gold-dark border-[0.5px] px-4 py-2"
          >
            <h1>Chapter {{ chapter.chapterNo }}</h1>
          </div>
          <p v-html="convertLineBreaksToHtml(toBionicText(chapter.body))"></p>
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
  </main>
</template>
