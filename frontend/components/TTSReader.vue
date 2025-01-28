<script setup lang="ts">
import type { TTSRequest } from '~/models/TTSRequest';

const chapterStore = useChapterStore();
const ttsStore = useTTSStore();
const userStore = useUserStore();

const { chapter } = storeToRefs(chapterStore);
const { user } = storeToRefs(userStore);
const {
    paragraphs, fetchingTTS, fetchingTTSError
} = storeToRefs(ttsStore);

const currentParagraph = ref<number>(0);
const audioPlayer = ref<HTMLAudioElement | null>(null);
const isPlaying = ref<boolean>(false);
const currentTime = ref<number>(0);
const duration = ref<number>(0);

// Play or pause the audio
const togglePlayback = () => {
    if (audioPlayer.value) {
        if (isPlaying.value) {
            audioPlayer.value.pause();
        } else {
            audioPlayer.value.play();
        }
        isPlaying.value = !isPlaying.value;
    }
};

// Update progress when the audio time changes
const updateProgress = () => {
    if (audioPlayer.value) {
        currentTime.value = audioPlayer.value.currentTime;
        duration.value = audioPlayer.value.duration;
    }
};

// Play audio for the current paragraph
const playAudio = () => {
    if (audioPlayer.value) {
        // Pause and reset the current audio
        audioPlayer.value.pause();
        audioPlayer.value.currentTime = 0;

        // Set the new audio source with a cache-busting query parameter
        const url = paragraphs.value[currentParagraph.value].url + '?t=' + Date.now();
        audioPlayer.value.src = url;

        // Wait for the audio to load before playing
        audioPlayer.value.addEventListener('canplay', () => {
            audioPlayer.value?.play();
            isPlaying.value = true;
        }, { once: true });

        // Update progress when the audio time changes
        audioPlayer.value.addEventListener('timeupdate', updateProgress);

        // Update duration when the audio metadata is loaded
        audioPlayer.value.addEventListener('loadedmetadata', () => {
            duration.value = audioPlayer.value?.duration || 0;
        });
    }
};

// Move to the next paragraph when the current audio ends
const playNextParagraph = () => {
    if (currentParagraph.value < paragraphs.value.length - 1) {
        currentParagraph.value++;
        playAudio();
    }
};

// Watch for changes to the current paragraph
watch(currentParagraph, (newValue) => {
    playAudio();
});

// Fetch TTS data and start playing the first paragraph
watchEffect(async() => {
    if (chapter.value && user.value) {
        const ttsRequest: TTSRequest = {
            text: chapter.value.body,
            novelId: chapter.value.novelId,
            chapterNo: chapter.value.chapterNo,
            voice: user.value.readingPreferences.tts.voice,
        };

        await ttsStore.generateTTS(ttsRequest);

        // Start playing audio for the first paragraph
        if (paragraphs.value.length > 0) {
            currentParagraph.value = 0;
            playAudio();
        }
    }
});
</script>

<template>
    <div>
        <!-- Audio element for playing paragraph audio -->
        <audio ref="audioPlayer" @ended="playNextParagraph"/>

        <!-- Play/Pause Button -->
        <div class="flex items-center gap-4 my-4">
            <button @click="togglePlayback" class="p-2 bg-primary text-white rounded">
                {{ isPlaying ? 'Pause' : 'Play' }}
            </button>

            <!-- Progress Display -->
            <div class="flex-1">
                <div class="w-full bg-gray-200 rounded h-2">
                    <div
                        class="bg-primary h-2 rounded"
                        :style=" {
                            width: `$ {(currentTime / duration) * 100}%`
                        }"
                    ></div>
                </div>
                <div class="text-sm text-gray-600 mt-1">
                    {{ Math.floor(currentTime) }}s / {{ Math.floor(duration) }}s
                </div>
            </div>
        </div>

        <!-- Paragraphs -->
        <div
            class="transition-colors"
            v-for="(paragraph, index) in paragraphs"
            :key="index"
            :id="`paragraph-$ {index}`"
            :style=" {
                opacity: 1 - Math.min(Math.abs(index - currentParagraph) / 5, 0.8)
            }"
            :class="index === currentParagraph ? 'text-primary p-[0.75rem] bg-primary-content' : 'text-secondary'"
        >
            <p class="font-bold font-mono my-9 md:text-lg lg:text-xl xl:text-2xl 2xl:text:3xl">
                {{ paragraph.text }}
            </p>
        </div>
    </div>
</template>