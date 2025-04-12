<script setup lang="ts">
    import type { TTSRequest } from '~/schemas/TTSRequest';

    const props = defineProps<{
        novelTitle: string;
    }>();

    const ttsStore = useTTSStore();
    const userStore = useUserStore();
    const bookmarkStore = useBookmarkStore();
    const chapterStore = useChapterStore();

    const { user } = storeToRefs(userStore);
    const { bookmark } = storeToRefs(bookmarkStore);
    const { audioPlayer, isPlaying, duration, paragraphs, fetchingTTS } = storeToRefs(ttsStore);
    const { fetchingChapters } = storeToRefs(chapterStore);

    const currentParagraph = ref<number>(0);

    // Clean up audio player event listeners
    const cleanupAudioPlayer = () => {
        if (audioPlayer.value) {
            audioPlayer.value.pause();
            audioPlayer.value.currentTime = 0;
            audioPlayer.value.src = '';
            audioPlayer.value.removeEventListener('canplay', handleCanPlay);
            audioPlayer.value.removeEventListener('timeupdate', ttsStore.updateProgress);
            audioPlayer.value.removeEventListener('loadedmetadata', handleLoadedMetadata);
            audioPlayer.value.removeEventListener('error', handleAudioError);
        }
    };

    // Handle audio playback when the audio is ready
    const handleCanPlay = () => {
        if (audioPlayer.value) {
            audioPlayer.value.play();
            isPlaying.value = true;
        }
    };

    // Handle metadata loading
    const handleLoadedMetadata = () => {
        if (audioPlayer.value) {
            duration.value = audioPlayer.value.duration || 0;
        }
    };

    // Handle audio errors
    const handleAudioError = async () => {
        const errorCode = audioPlayer.value?.error?.code || 'unknown';
        console.error('Error playing audio for paragraph ' + currentParagraph.value, 'Error Code: ' + errorCode);

        // Only skip if the error is unrecoverable
        if (errorCode === 4) {
            // Code 4: MEDIA_ERR_SRC_NOT_SUPPORTED
            console.warn('Invalid source. Skipping paragraph...');
            await playNextParagraph();
        }
    };

    // Play audio for the current paragraph
    const playAudio = () => {
        cleanupAudioPlayer();

        if (!audioPlayer.value || paragraphs.value.length === 0 || !paragraphs.value[currentParagraph.value]) return;
        scrollToParagraph();

        // Set the new audio source with a cache-busting query parameter
        audioPlayer.value.src = paragraphs.value[currentParagraph.value].url + '?t=' + Date.now();

        // Add new event listeners
        audioPlayer.value.addEventListener('canplay', handleCanPlay, { once: true });
        audioPlayer.value.addEventListener('timeupdate', ttsStore.updateProgress);
        audioPlayer.value.addEventListener('loadedmetadata', handleLoadedMetadata);
        audioPlayer.value.addEventListener('error', handleAudioError);
    };

    const scrollToParagraph = () => {
        const paragraphElement = document.getElementById('paragraph-' + currentParagraph.value);
        if (paragraphElement) {
            paragraphElement.scrollIntoView({ behavior: 'smooth', block: 'center' });
        }
    };

    // Move to the next paragraph when the current audio ends
    const playNextParagraph = async () => {
        if (currentParagraph.value < paragraphs.value.length - 1) {
            currentParagraph.value++;
        } else {
            if (!chapterStore.chapter) return;

            const nextChapter = chapterStore?.chapter.chapterNo + 1;
            // Stop playback if we've reached the end
            isPlaying.value = false;
            fetchingTTS.value = true;
            paragraphs.value = [];

            // Update user currentChapter
            if (bookmark.value) {
                bookmark.value.currentChapter = nextChapter;
                try {
                    await bookmarkStore.updateBookmark(bookmark.value);
                } catch {}
            }

            // Navigate to new current chapter
            if (user.value && user.value.readingPreferences.tts.autoplay) {
                await navigateTo('/novels/' + props.novelTitle + '/' + nextChapter);
            }
        }
    };

    // Watch for changes to the current paragraph
    watch(currentParagraph, _ => {
        if (!fetchingTTS.value && isPlaying.value) {
            playAudio();
        }
    });

    watch(fetchingChapters, async () => {
        fetchingTTS.value = true;
        paragraphs.value = [];
        cleanupAudioPlayer();
        isPlaying.value = false;

        if (!chapterStore.chapter || !user.value) return;

        if (!user.value.readingPreferences.tts.voice || user.value.readingPreferences.tts.voice === '' || user.value.readingPreferences.tts.voice === ' ') {
            user.value.readingPreferences.tts.voice = 'en-US-AriaNeural';
        }

        const ttsRequest: TTSRequest = {
            text: chapterStore.chapter.body,
            novelId: chapterStore.chapter.novelId,
            chapterNo: chapterStore.chapter.chapterNo,
            voice: user.value.readingPreferences.tts.voice,
            rate: user.value.readingPreferences.tts.rate || 0,
        };

        try {
            await ttsStore.generateTTS(ttsRequest);
        } catch {}
        // Reset state for the new chapter
        currentParagraph.value = 0;

        // Start playing audio for the first paragraph
        if (paragraphs.value.length > 0) {
            playAudio();
        }
        {
            flush: 'post';
        }
    });
</script>

<template>
    <LoadingBar v-show="fetchingTTS || fetchingChapters" />
    <div v-show="!fetchingTTS && paragraphs.length > 0">
        <!-- Audio element for playing paragraph audio -->
        <audio ref="audioPlayer" @ended="playNextParagraph" />

        <!-- Paragraphs -->
        <div
            class="transition-colors"
            v-for="(paragraph, index) in paragraphs"
            :key="index"
            :id="'paragraph-' + index"
            :style="{
                opacity: 1 - Math.min(Math.abs(index - currentParagraph) / 5, 0.8),
            }"
            :class="index === currentParagraph ? 'bg-primary-content p-[0.75rem] text-primary' : 'text-secondary'"
            @dblclick="currentParagraph = index"
        >
            <p class="2xl:text:3xl my-9 font-bold md:text-lg lg:text-xl xl:text-2xl" :class="user ? user.readingPreferences.font : 'font-mono'">
                {{ paragraph.text }}
            </p>
        </div>
    </div>
</template>
