import { AuthError } from '~/errors/AuthError';
import type { Paragraph } from '~/schemas/Paragraph';
import type { TTSRequest } from '~/schemas/TTSRequest';
import { TtsService } from '~/services/TTSService';

export const useTTSStore = defineStore('TTS', () => {
  const authStore = useAuthStore();
  const runtimeConfig = useRuntimeConfig();
  const url: string = runtimeConfig.public.apiUrl;

  const ttsService: TtsService = new TtsService(url);

  const { user } = storeToRefs(authStore);

  const paragraphs = ref<Paragraph[]>([]);
  const fetchingTTS = ref<boolean>(false);

  const fetchingTTSVoices = ref<boolean>(false);

  const currentTime = ref<number>(0);
  const duration = ref<number>(0);
  const audioPlayer = ref<HTMLAudioElement | null>(null);
  const isPlaying = ref<boolean>(false);

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

  const generateTTS = async (ttsRequest: TTSRequest) => {
    fetchingTTS.value = true;

    try {
      if (user.value === null) {
        throw new AuthError({
          name: 'UNAUTHORIZED_ERROR',
          message: "You're not logged in!",
          cause: 'User tried to generate TTS without being logged in.',
        });
      }

      paragraphs.value = await ttsService.generateTTS(ttsRequest);
    } catch (error) {
      handleError(error, { user: user, ttsRequest: ttsRequest, location: 'tts.ts -> generateTTS' });
      throw error;
    } finally {
      fetchingTTS.value = false;
    }
  };

  const fetchTTSVoices = async () => {
    fetchingTTSVoices.value = true;

    try {
      if (user.value === null) {
        throw new AuthError({
          name: 'UNAUTHORIZED_ERROR',
          message: "You're not logged in!",
          cause: 'User tried to fetch TTS voices without being logged in.',
        });
      }

      const voices = await ttsService.fetchTTSVoices();
      console.log(voices);
    } catch (error) {
      handleError(error, { user: user, location: 'tts.ts -> fetchTTSVoices' });
      throw error;
    } finally {
      fetchingTTS.value = false;
    }
  };

  return {
    isPlaying,
    audioPlayer,
    currentTime,
    duration,
    paragraphs,
    fetchingTTS,
    fetchingTTSVoices,
    togglePlayback,
    updateProgress,
    generateTTS,
    fetchTTSVoices,
  };
});

if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useTTSStore, import.meta.hot));
}
