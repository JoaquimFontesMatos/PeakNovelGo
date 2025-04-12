<script setup lang="ts" xmlns="http://www.w3.org/1999/html">
    const authStore = useAuthStore();
    const ttsStore = useTTSStore();
    const chapterStore = useChapterStore();
    const userStore = useUserStore();
    const novelStore = useNovelStore();

    defineProps<{
        drawerOpen: boolean;
        onGoToPreviousChapter: () => Promise<void>;
        onGoToNextChapter: () => Promise<void>;
    }>();

    type Tabs = 'general' | 'display' | 'audio' | 'translate';
    const tabs: Tabs[] = ['general', 'display', 'audio', 'translate'];

    const currentTab: Ref<Tabs> = ref('general');

    const { currentTime, isPlaying, duration, audioPlayer } = storeToRefs(ttsStore);
    const { chapter, fetchingChapters, paginatedChapterData, novelProgress } = storeToRefs(chapterStore);
    const { user, isReaderMode } = storeToRefs(userStore);

    const handleSelectTab = (tab: Tabs) => {
        currentTab.value = tab;
    };

    const handleChangeReadingMode = async () => {
        if (isReaderMode) audioPlayer.value = null;
    };
</script>

<template>
    <section v-show="drawerOpen">
        <SlidingDrawer v-if="drawerOpen">
            <!-- Drawer content -->
            <div class="form-container flex w-full items-center justify-between pt-2 md:pt-0">
                <TabButton v-for="tab in tabs" :key="tab" :name="tab" :current-tab="currentTab" @select-tab="handleSelectTab(tab)" />
            </div>

            <SmallVerticalSpacer />

            <section v-if="currentTab === 'general'">
                <div class="menu-container">
                    <div class="align-center flex w-full justify-center">
                        <div class="flex w-full flex-row gap-1 text-xs text-secondary-content">
                            <span>{{ novelProgress }}%</span>
                            <div class="h-full w-[3px] bg-secondary-content py-1" />
                            <span class="line-clamp-1">{{ chapter?.title ?? 0 }}</span>
                        </div>
                    </div>

                    <div class="form-container flex w-full items-center justify-between">
                        <CircularButton
                            :disabled="chapterStore.chapter?.chapterNo === 1"
                            :padding="4"
                            :icon-name="'fluent:previous-28-filled'"
                            :no-background="true"
                            @click="onGoToPreviousChapter()"
                        />

                        <input type="range" min="0" max="100" class="mx-4" v-model="novelProgress" disabled />

                        <CircularButton
                            :disabled="chapterStore.chapter?.chapterNo === paginatedChapterData?.total"
                            :padding="4"
                            :icon-name="'fluent:next-28-filled'"
                            :no-background="true"
                            @click="onGoToNextChapter()"
                        />
                    </div>

                    <div class="form-container flex w-full items-center justify-between px-6">
                        <CircularButton
                            :padding="4"
                            :icon-name="'fluent:book-open-32-regular'"
                            :icon-size="28"
                            @click="navigateTo(`/novels/${novelStore.novel?.novelUpdatesId ?? ''}`)"
                        />

                        <CircularButton :padding="4" :icon-name="'fluent:home-32-regular'" :icon-size="28" @click="navigateTo('/')" />

                        <CircularButton :padding="4" :icon-name="'fluent:comment-multiple-32-regular'" :icon-size="28" @click="navigateTo('/')" />
                    </div>
                </div>

                <SmallVerticalSpacer />

                <div class="menu-container">
                    <SectionHeader :title="'Downloads'" :is-main-header="false">
                        <div class="form-container flex w-full flex-col items-start justify-between gap-2 px-6">
                            <div class="flex flex-wrap items-center gap-2">
                                <CircularButton :padding="3" :icon-name="'fluent:arrow-download-32-regular'" :icon-size="24" />
                                <p>Download this chapter</p>
                            </div>
                            <div class="flex flex-wrap items-center gap-2">
                                <CircularButton :padding="3" :icon-name="'fluent:arrow-download-32-regular'" :icon-size="24" />
                                <p>Download all</p>
                            </div>
                        </div>
                    </SectionHeader>
                </div>
            </section>

            <section v-else-if="currentTab === 'display'">
                <ReadingPreferencesSettings v-if="authStore.isUserLoggedIn()" />
            </section>

            <section v-else-if="currentTab === 'audio'">
                <div class="menu-container form-container">
                    <div v-if="authStore.isUserLoggedIn()" class="form-group flex items-center space-x-2">
                        <input id="readerMode" name="readerMode" type="checkbox" v-model="isReaderMode" @change="handleChangeReadingMode()" />
                        <label for="readerMode" class="text-sm font-medium text-secondary-content">Reader Mode</label>
                    </div>

                    <!-- Play/Pause Button -->
                    <fieldset v-if="audioPlayer" class="flex items-center gap-4 border-t border-accent-gold-dark pt-4">
                        <legend class="ml-3.5 px-3.5 text-lg font-semibold text-primary-content">Audio Controls</legend>
                        <!-- Atomic Reading -->
                        <button @click="ttsStore.togglePlayback" class="flex items-center rounded-sm bg-primary/80 p-2 text-white">
                            <Icon v-if="isPlaying" name="fluent:pause-28-filled" class="text-accent-gold-dark" />
                            <Icon v-else name="fluent:play-28-filled" class="text-accent-gold-dark" />
                        </button>

                        <!-- Progress Display -->
                        <div class="flex-1">
                            <div class="h-2 w-full rounded-sm bg-gray-200">
                                <div
                                    class="h-2 rounded-sm bg-accent-gold-dark"
                                    :style="{
                                        width: (currentTime / duration) * 100 + '%',
                                    }"
                                />
                            </div>
                            <div class="mt-1 text-sm text-gray-600">{{ Math.floor(currentTime) }}s / {{ Math.floor(duration) }}s</div>
                        </div>
                    </fieldset>
                </div>
            </section>
        </SlidingDrawer>
    </section>
</template>
