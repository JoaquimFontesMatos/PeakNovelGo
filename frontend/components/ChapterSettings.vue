<script setup lang="ts" xmlns="http://www.w3.org/1999/html">

import TabButton from "~/components/TabButton.vue";
import SectionHeader from "~/components/SectionHeader.vue";
import CircularButton from "~/components/CircularButton.vue";

const authStore = useAuthStore()
const ttsStore = useTTSStore()
const chapterStore = useChapterStore()
const userStore = useUserStore()
const novelStore = useNovelStore()

defineProps<{
  drawerOpen: boolean,
  currentChapter: number,
}>()

const emit = defineEmits(['goToPreviousChapter', 'goToNextChapter']);

const handleGoToPreviousChapter = async () => {
  emit('goToPreviousChapter');
};

const handleGoToNextChapter = async () => {
  emit('goToNextChapter');
};

type Tabs = 'general' | 'display' | 'audio' | 'translate';
const tabs: Tabs[] = ['general', 'display', 'audio', 'translate'];

const sectionHeight = ref('33%');
const isDragging = ref(false);
const initialY = ref(0);
const initialHeight = ref(0);
const isLargeScreen = ref(false);
const currentTab: Ref<Tabs> = ref('general')

const {currentTime, isPlaying, duration, audioPlayer} = storeToRefs(ttsStore);
const {chapter, fetchingChapters, paginatedChapterData, novelProgress} = storeToRefs(chapterStore);
const {user, isReaderMode} = storeToRefs(userStore);

const startDragging = (event: any) => {
  isDragging.value = true;
  initialY.value = event.type === 'mousedown' ? event.clientY : event.touches[0].clientY;
  initialHeight.value = parseFloat(sectionHeight.value); // Store the initial height

  document.addEventListener('mousemove', handleDragging);
  document.addEventListener('mouseup', stopDragging);

  document.addEventListener('touchmove', handleDragging);
  document.addEventListener('touchend', stopDragging);
  document.addEventListener('touchcancel', stopDragging);
};

const handleDragging = (event: any) => {
  if (isDragging.value) {
    const currentY = event.type === 'mousemove' ? event.clientY : event.touches[0].clientY;
    const deltaY = initialY.value - currentY; // Calculate the difference in mouse position
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

  document.removeEventListener('touchmove', handleDragging);
  document.removeEventListener('touchend', stopDragging);
  document.removeEventListener('touchcancel', stopDragging);
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

const handleSelectTab = (tab: Tabs) => {
  currentTab.value = tab
}
</script>

<template>
  <section
      :style="drawerOpen && !isLargeScreen ? { height: sectionHeight } : {}"
      :class="drawerOpen ? 'md:w-1/3 md:h-svh py-4' : 'h-0 w-0 m-0 shadow-none border-none backdrop-filter-none'"
      class="bg-primary-100/10 fixed left-0 z-50 w-full select-none overflow-y-scroll px-4 shadow-md backdrop-blur-sm transition-all duration-0 bottom-0 md:left-2/3"
  >
    <!-- Draggable Divider -->
    <div
        @mousedown="startDragging"
        @touchstart="startDragging"
        class="md:hidden absolute -top-1 left-0 right-0 h-3 cursor-row-resize bg-border hover:bg-secondary-content transition-colors duration-200 m-1 rounded "
    >
      <div class="w-full h-full flex flex-col justify-center items-center gap-[0.1rem] p-0.5">
        <div class="w-[8%] h-auto grow bg-secondary-content"/>
        <div class="w-[8%] h-auto grow bg-secondary-content"/>
        <div class="w-[8%] h-auto grow bg-secondary-content"/>
      </div>
    </div>

    <div class="form-container flex w-full items-center justify-between pt-2 md:pt-0">
      <TabButton
          v-for="tab in tabs"
          :key="tab"
          :name="tab"
          :current-tab="currentTab"
          @select-tab="handleSelectTab(tab)"
      />
    </div>

    <SmallVerticalSpacer/>

    <section v-if="currentTab==='general'">
      <div class="menu-container">
        <div class="w-full flex justify-center align-center">
          <p class="text-xs text-secondary-content w-full flex flex-row gap-1">
            <span>{{ novelProgress }}%</span>
            <div class="w-[3px] h-full py-1 bg-secondary-content"/>
            <span class="line-clamp-1">{{ chapter?.title ?? 0 }}</span>
          </p>
        </div>

        <div class="form-container flex w-full items-center justify-between">
          <CircularButton :disabled="currentChapter === 1"
                          :padding="4"
                          :icon-name="'fluent:previous-28-filled'"
                          :no-background="true"
                          @click="handleGoToPreviousChapter()"/>

          <input type="range" min="0" max="100" class="mx-4" v-model="novelProgress" disabled>

          <CircularButton :disabled="currentChapter === paginatedChapterData?.total"
                          :padding="4"
                          :icon-name="'fluent:next-28-filled'"
                          :no-background="true"
                          @click="handleGoToNextChapter()"/>
        </div>

        <div class="form-container flex w-full items-center justify-between px-6">
          <CircularButton :padding="4" :icon-name="'fluent:book-open-32-regular'" :icon-size="28"
                          @click="navigateTo(`/novels/${novelStore.novel?.novelUpdatesId??''}`)"/>

          <CircularButton :padding="4" :icon-name="'fluent:home-32-regular'" :icon-size="28"
                          @click="navigateTo('/')"/>

          <CircularButton :padding="4" :icon-name="'fluent:comment-multiple-32-regular'" :icon-size="28"
                          @click="navigateTo('/')"/>
        </div>
      </div>

      <SmallVerticalSpacer/>

      <div class="menu-container">
        <SectionHeader :title="'Downloads'" :is-main-header="false">
          <div class="form-container flex flex-col w-full items-start justify-between gap-2 px-6">
            <div class="flex flex-wrap gap-2 items-center">
              <CircularButton :padding="3" :icon-name="'fluent:arrow-download-32-regular'" :icon-size="24"/>
              <p>Download this chapter</p>
            </div>
            <div class="flex flex-wrap gap-2 items-center">
              <CircularButton :padding="3" :icon-name="'fluent:arrow-download-32-regular'" :icon-size="24"/>
              <p>Download all</p>
            </div>
          </div>
        </SectionHeader>
      </div>
    </section>

    <section v-else-if="currentTab==='display'">
      <ReadingPreferencesForm v-if="authStore.isUserLoggedIn()"/>
    </section>

    <section v-else-if="currentTab==='audio'">
      <div class="menu-container form-container">
        <div v-if="authStore.isUserLoggedIn()" class="form-group flex items-center space-x-2">
          <input id="readerMode" name="readerMode" type="checkbox" v-model="isReaderMode"
                 @change="isReaderMode ? (audioPlayer = null) : console.log('')"/>
          <label for="readerMode" class="text-sm font-medium text-secondary-content">Reader Mode</label>
        </div>

        <!-- Play/Pause Button -->
        <fieldset v-if="audioPlayer" class="flex items-center gap-4 border-t border-accent-gold-dark pt-4">
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
    </section>
  </section>
</template>