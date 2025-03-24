<script setup lang="ts">

import {useScroll} from '@vueuse/core';
import {BaseTextUtils} from '~/composables/textUtils';

const {novelTitle, chapterNo} = useRoute().params as { novelTitle: string; chapterNo: string };
definePageMeta({
  ssr: false,
});

useHead({
  title: `📖 ${novelTitle.replace(/-/g, ' ').replace(/\b\w/g, (char) => char.toUpperCase())}: Chapter ${chapterNo}`,
  meta: [
    {
      name: 'description',
      content: `Read ${novelTitle}, Chapter ${chapterNo} on our platform.`,
    },
  ],
});

const textUtils = new BaseTextUtils();

const avgReadingSpeed = 238;

const chapterStore = useChapterStore();
const ttsStore = useTTSStore();
const bookmarkStore = useBookmarkStore();
const userStore = useUserStore();
const authStore = useAuthStore();
const currentChapter = ref(0);

const {chapter, fetchingChapters, paginatedChapterData} = storeToRefs(chapterStore);

const {bookmark, fetchingBookmarkedNovel, updatingBookmark} = storeToRefs(bookmarkStore);

const {user, isReaderMode} = storeToRefs(userStore);

const {currentTime, isPlaying, duration, audioPlayer} = storeToRefs(ttsStore);

const fetchBookmark = async (novelId: string): Promise<void> => {
  try {
    await useBookmarkStore().fetchBookmarkedNovelByUser(novelId);
  } catch {
  }
};

const goToPreviousChapter = async (): Promise<void> => {
  if (authStore.isUserLoggedIn()) {
    if (!bookmark.value) return;

    bookmark.value.currentChapter = currentChapter.value - 1;

    try {
      await bookmarkStore.updateBookmark(bookmark.value);
    } catch {
    }
  }
  navigateTo((('/novels/' + novelTitle) as string) + '/' + (currentChapter.value - 1));
};

const goToNextChapter = async (): Promise<void> => {
  if (user.value) {
    if (!bookmark.value) return;
    bookmark.value.currentChapter = currentChapter.value + 1;

    try {
      if (authStore.isUserLoggedIn()) {
        await bookmarkStore.updateBookmark(bookmark.value);
      }
    } catch {
    }
  }
  navigateTo((('/novels/' + novelTitle) as string) + '/' + (currentChapter.value + 1));
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
  } catch {
  }

  if (authStore.isUserLoggedIn()) {
    await fetchBookmark(novelTitle as string);
  }
});

/**
 * Manage the reading prefs
 */
const drawerOpen = ref<boolean>(false);

// Use the useScroll function to track scroll position reactively
const {y} = useScroll(window);

const scrollProgress = computed(() => {
  if (import.meta.client) {
    // Ensure this code runs only on the client-side
    const scrollTop = y.value || 0;
    const scrollHeight = document.documentElement.scrollHeight - window.innerHeight;
    return scrollTop / scrollHeight;
  }
  return 0;
});

const keepSessionAlive = () => {
  setInterval(async () => {
    await authStore.keepAlive()
  }, 5 * 60 * 1000);
};

onMounted(() => {
  keepSessionAlive();
});

// In your script setup
const calculateReadingTime = (text: string): number => {
  if (!text) return 0;

  const wordCount = (text.match(/\b\w+\b/g)?.length || 0);
  return Math.ceil(wordCount / avgReadingSpeed);
};

const sectionHeight = ref('33%'); // Default height
const isDragging = ref(false);
const initialY = ref(0);
const initialHeight = ref(0);
const isLargeScreen = ref(false);

const startDragging = (event: MouseEvent) => {
  isDragging.value = true;
  initialY.value = event.clientY;
  initialHeight.value = parseFloat(sectionHeight.value); // Store the initial height

  document.addEventListener('mousemove', handleDragging);
  document.addEventListener('mouseup', stopDragging);
};

const handleDragging = (event: MouseEvent) => {
  if (isDragging.value) {
    const deltaY = initialY.value - event.clientY; // Calculate the difference in mouse position
    const newHeight = initialHeight.value + deltaY;

    // Clamp the height between 100px and 75vh
    const maxHeight = window.innerHeight * 0.75; // 75% of viewport height
    sectionHeight.value = `${Math.max(100, Math.min(newHeight, maxHeight))}px`;
  }
};

const stopDragging = () => {
  isDragging.value = false;
  document.removeEventListener('mousemove', handleDragging);
  document.removeEventListener('mouseup', stopDragging);
};

const checkScreenSize = () => {
  isLargeScreen.value = window.innerWidth >= 768; // 768px is the default breakpoint for `md` in Tailwind
};

onMounted(() => {
  checkScreenSize();
  window.addEventListener('resize', checkScreenSize);
});

onUnmounted(() => {
  window.removeEventListener('resize', checkScreenSize);
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
        {
          path: '/novels/' + novelTitle + '/' + chapterNo,
          name: 'Chapter ' + chapterNo as string
        },
      ]"
    />

    <VerticalSpacer/>

    <LoadingBar v-show="fetchingChapters"/>
    <client-only>
      <div
          class="fixed left-0 top-0 z-50 h-[2px] w-full origin-left bg-accent-gold-dark transition-transform duration-500 md:h-1"
          :style="{
        transform: 'scaleX(' + scrollProgress + ')',
      }"
      />
      <section v-if="isReaderMode" v-show="!fetchingChapters" @click="drawerOpen = !drawerOpen">
        <div v-if="chapter">
          <div>
            <div
                class="mb-4 rounded-md border-[0.5px] border-accent-gold-dark bg-secondary px-4 py-2 text-secondary-content">
              <h1>Chapter {{ chapter.chapterNo }}</h1>
              <h3 class="opacity-50 text-xs">Approx. {{ calculateReadingTime(chapter.body) }} minutes</h3>
            </div>
            <p
                :class="user ? user.readingPreferences.font : ''"
                v-html="user && user.readingPreferences.atomicReading ? textUtils.convertLineBreaksToHtml(textUtils.toBionicText(chapter.body)) : textUtils.convertLineBreaksToHtml(chapter.body)"
            />
          </div>

          <VerticalSpacer/>
          <Button @click="goToPreviousChapter" :disabled="currentChapter === 1"> << Previous Chapter</Button>

          <Button @click="goToNextChapter"
                  :disabled="currentChapter && paginatedChapterData && currentChapter === paginatedChapterData.total">
            Next Chapter >>
          </Button>
        </div>

        <ErrorAlert v-else>Error: {{ chapter === null ? 'Invalid Chapter Number' : null }}</ErrorAlert>
      </section>
      <section v-else v-show="!fetchingChapters" @click="drawerOpen = !drawerOpen">
        <TTSReader :novel-title="novelTitle as string" :chapter="chapter"/>
      </section>

      <section
          :style="drawerOpen && !isLargeScreen ? { height: sectionHeight } : {}"
          :class="drawerOpen ? 'md:w-1/3 md:h-svh py-4' : 'h-0 w-0 m-0 shadow-none border-none backdrop-filter-none'"
          class="bg-primary-100/10 fixed left-0 z-50 w-full select-none overflow-y-scroll px-4 shadow-md backdrop-blur-sm transition-all duration-0 bottom-0 md:left-2/3"
      >
        <!-- Draggable Divider -->
        <div
            @mousedown="startDragging"
            class="md:hidden absolute -top-1 left-0 right-0 h-3 cursor-row-resize bg-border hover:bg-secondary-content transition-colors duration-200 m-1 rounded "
        >
          <div class="w-full h-full flex flex-col justify-center items-center gap-[0.1rem] p-0.5">
            <div class="w-[8%] h-auto grow bg-secondary-content"/>
            <div class="w-[8%] h-auto grow bg-secondary-content"/>
            <div class="w-[8%] h-auto grow bg-secondary-content"/>
          </div>
          <SmallVerticalSpacer/>
        </div>

        <section>
          <div>
            <h1>Chapter {{ chapter?.chapterNo ?? 0 }}</h1>
          </div>

          <SmallVerticalSpacer/>

          <div class="form-container flex w-full items-center justify-between rounded-lg bg-secondary p-4 shadow-lg">
            <button
                @click="goToPreviousChapter"
                :disabled="currentChapter === 1"
                class="bg-primary/80 disabled:bg-bg-primary/20 hover:bg-primary/40 flex items-center rounded-lg p-1 transition-colors disabled:cursor-not-allowed disabled:text-gray-300"
            >
              <Icon name="fluent:previous-28-filled" class="text-accent-gold-dark"/>
            </button>

            <div v-if="authStore.isUserLoggedIn()" class="form-group flex items-center space-x-2">
              <input id="readerMode" name="readerMode" type="checkbox" v-model="isReaderMode"
                     @change="isReaderMode ? (audioPlayer = null) : console.log('')"/>
              <label for="readerMode" class="text-sm font-medium text-secondary-content">Reader Mode</label>
            </div>

            <button
                @click="goToNextChapter"
                :disabled="currentChapter === paginatedChapterData?.total"
                class="bg-primary/80 disabled:bg-bg-primary/20 hover:bg-primary/40 flex items-center rounded-lg p-1 transition-colors disabled:cursor-not-allowed disabled:text-gray-300"
            >
              <Icon name="fluent:next-28-filled" class="text-accent-gold-dark"/>
            </button>
          </div>

          <!-- Play/Pause Button -->
          <SmallVerticalSpacer v-if="audioPlayer"/>

          <div v-if="audioPlayer" class="form-container space-y-6 rounded-lg bg-secondary p-4 shadow-lg">
            <fieldset class="flex items-center gap-4 border-t border-accent-gold-dark pt-4">
              <legend class="ml-3.5 px-3.5 text-lg font-semibold text-primary-content">Audio Controls</legend>
              <!-- Atomic Reading -->
              <button @click="ttsStore.togglePlayback" class="bg-primary/80 flex items-center rounded p-2 text-white">
                <Icon v-if="isPlaying" name="fluent:pause-28-filled" class="text-accent-gold-dark"/>
                <Icon v-else name="fluent:play-28-filled" class="text-accent-gold-dark"/>
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
                <div class="mt-1 text-sm text-gray-600">{{ Math.floor(currentTime) }}s / {{
                    Math.floor(duration)
                  }}s
                </div>
              </div>
            </fieldset>
          </div>

          <SmallVerticalSpacer/>

          <ReadingPreferencesForm v-if="authStore.isUserLoggedIn()"/>
        </section>
      </section>
    </client-only>
  </Container>
</template>
