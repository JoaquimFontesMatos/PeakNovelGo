<script setup lang="ts">
    import { useScroll } from '@vueuse/core';
    import { useIndexedDB } from '~/composables/useInitCacheDB';

    const { novelTitle, chapterNo } = useRoute().params as {
        novelTitle: string;
        chapterNo: string;
    };
    definePageMeta({
        ssr: false,
        layout: 'custom',
    });

    useHead({
        title: `📖 ${novelTitle.replace(/-/g, ' ').replace(/\b\w/g, char => char.toUpperCase())}: Chapter ${chapterNo}`,
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
    const novelStore = useNovelStore();
    const bookmarkStore = useBookmarkStore();
    const userStore = useUserStore();
    const authStore = useAuthStore();
    const currentChapter = ref(0);

    const { cachedChapters, chapter, fetchingChapters, paginatedChapterData } = storeToRefs(chapterStore);

    const { bookmark, fetchingBookmarkedNovel, updatingBookmark } = storeToRefs(bookmarkStore);

    const { user, isReaderMode } = storeToRefs(userStore);

    const { novel } = storeToRefs(novelStore);

    const fetchBookmark = async (novelId: string): Promise<void> => {
        if (!novel.value) {
            return;
        }

        if (import.meta.client) {
            const cachedBookmark = await bookmarkStore.getCachedBookmark(novel.value.ID);

            if (cachedBookmark) {
                bookmark.value = cachedBookmark;
                console.log(`Bookmark from novel with id ${bookmark.value.novelId} loaded from cache`);
            } else {
                // If not cached, fetch from the server
                await useBookmarkStore().fetchBookmarkedNovelByUser(novelId);
            }

            await bookmarkStore.cacheBookmark();
        }
    };

    const goToPreviousChapter = async (): Promise<void> => {
        const previousChapter = currentChapter.value - 1;

        if (authStore.isUserLoggedIn()) {
            if (!bookmark.value) {
                await navigateTo((('/novels/' + novelTitle) as string) + '/' + previousChapter);
                return;
            }

            bookmark.value.currentChapter = previousChapter;

            try {
                await bookmarkStore.updateBookmark(bookmark.value);
            } catch {}
        }
        await navigateTo((('/novels/' + novelTitle) as string) + '/' + previousChapter);
    };

    const goToNextChapter = async (): Promise<void> => {
        const nextChapter = currentChapter.value + 1;

        if (user.value) {
            if (!bookmark.value) {
                await navigateTo((('/novels/' + novelTitle) as string) + '/' + nextChapter);
                return;
            }
            bookmark.value.currentChapter = nextChapter;

            try {
                if (authStore.isUserLoggedIn()) {
                    await bookmarkStore.updateBookmark(bookmark.value);
                }
            } catch {}
        }
        await navigateTo((('/novels/' + novelTitle) as string) + '/' + nextChapter);
    };

    watchEffect(async () => {
        const novelUpdatesId = novelTitle as string;
        const chapterNum = parseInt(chapterNo as string);

        if (chapterNum === undefined || chapterNum === 0) {
            return;
        }

        currentChapter.value = chapterNum;

        try {
            if (import.meta.client) {
                const cachedChapter = await chapterStore.getCachedChapter(novelUpdatesId, currentChapter.value);

                if (cachedChapter) {
                    chapter.value = cachedChapter;
                    console.log(`Chapter ${currentChapter.value} loaded from cache`);
                } else {
                    // If not cached, fetch from the server
                    await chapterStore.fetchChapter(novelUpdatesId, currentChapter.value);
                }
            }
        } catch {}

        if (authStore.isUserLoggedIn()) {
            await fetchBookmark(novelTitle as string);
        }

        if (paginatedChapterData.value === null) {
            await chapterStore.fetchChapters(novelUpdatesId, 1, 10);
        }

        if (import.meta.client) {
            await chapterStore.cacheNextChapters(novelUpdatesId, currentChapter.value, 5);
        }
    });

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
        return 0;
    });

    const keepSessionAlive = () => {
        setInterval(
            async () => {
                await authStore.keepAlive();
            },
            5 * 60 * 1000
        );
    };

    onMounted(() => {
        keepSessionAlive();
    });

    // In your script setup
    const calculateReadingTime = (text: string): number => {
        if (!text) return 0;

        const wordCount = text.match(/\b\w+\b/g)?.length || 0;
        return Math.ceil(wordCount / avgReadingSpeed);
    };
</script>

<template>
    <Container>
        <RouteTree
            :routes="[
                { path: '/', name: 'Home' },
                { path: '/novels', name: 'Novels' },
                {
                    path: '/novels/' + novelTitle,
                    name: novelTitle as string,
                },
                {
                    path: '/novels/' + novelTitle + '/' + chapterNo,
                    name: ('Chapter ' + chapterNo) as string,
                },
            ]"
        />

        <VerticalSpacer />

        <LoadingBar v-show="fetchingChapters" />
        <client-only>
            <div
                class="fixed top-0 left-0 z-50 h-[2px] w-full origin-left bg-accent-gold-dark transition-transform duration-500 md:h-1"
                :style="{
                    transform: 'scaleX(' + scrollProgress + ')',
                }"
            />
            <section v-if="isReaderMode" v-show="!fetchingChapters">
                <div v-if="chapter">
                    <div @click="drawerOpen = !drawerOpen">
                        <div class="mb-4 rounded-md border-[0.5px] border-accent-gold-dark bg-secondary px-4 py-2 text-secondary-content">
                            <h1>Chapter {{ chapter.chapterNo }}</h1>
                            <h3 class="text-xs opacity-50">Approx. {{ calculateReadingTime(chapter.body) }} minutes</h3>
                        </div>
                        <p
                            :class="user ? user.readingPreferences.font : ''"
                            v-html="
                                user && user.readingPreferences.atomicReading
                                    ? textUtils.convertLineBreaksToHtml(textUtils.toBionicText(chapter.body))
                                    : textUtils.convertLineBreaksToHtml(chapter.body)
                            "
                        />
                    </div>
                    <VerticalSpacer />
                    <div class="flex w-full justify-between px-6 md:px-[25%]">
                        <MainButton @click="goToPreviousChapter()" :disabled="currentChapter === 1" class="flex flex-row items-center justify-center gap-1">
                            <span><<</span>
                            <span class="hidden md:block">Previous Chapter</span>
                        </MainButton>

                        <MainButton
                            @click="goToNextChapter()"
                            :disabled="currentChapter && paginatedChapterData && currentChapter === paginatedChapterData.total"
                            class="flex flex-row items-center justify-center gap-1"
                        >
                            <span class="hidden md:block">Previous Chapter</span>
                            <span>>></span>
                        </MainButton>
                    </div>
                </div>

                <ErrorAlert v-else>Error: {{ chapter === null ? 'Invalid Chapter Number' : null }}</ErrorAlert>
            </section>

            <section v-else v-show="!fetchingChapters" @click="drawerOpen = !drawerOpen">
                <TTSReader :novel-title="novelTitle as string" />
            </section>

            <ChapterSettings :drawer-open="drawerOpen" @goToPreviousChapter="goToPreviousChapter" @goToNextChapter="goToNextChapter" />
        </client-only>
    </Container>
</template>
