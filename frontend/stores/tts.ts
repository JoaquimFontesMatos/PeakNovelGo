import type { Paragraph } from '~/models/Paragraph';
import type { TTSRequest } from '~/models/TTSRequest';

export const useTTSStore = defineStore('TTS', () => {
    const authStore = useAuthStore();
    const runtimeConfig = useRuntimeConfig();
    const url: string = runtimeConfig.public.apiUrl;

    const { user } = storeToRefs(authStore);

    const paragraphs = ref<Paragraph[]>([]);
    const fetchingTTS = ref<boolean>(false);
    const fetchingTTSError = ref<string | null>(null);

    const fetchingTTSVoices = ref<boolean>(false);
    const fetchingTTSVoicesError = ref<string | null>(null);

    const generateTTS = async (ttsRequest: TTSRequest) => {
        fetchingTTS.value = true;
        fetchingTTSError.value = null;

        try {

            if (user.value === null) {
                fetchingTTS.value = false;
                fetchingTTSError.value = 'You\'re not logged in!';
                paragraphs.value = [];
                return;
            }

            const response: Response | undefined = await authStore.authorizedFetch(url + '/novels/tts',
                {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(ttsRequest),
                });

            if (response === null) {
                fetchingTTS.value = false;
                fetchingTTSError.value = 'Unexpected Error';
                paragraphs.value = [];
                return;
            }

            if (response.status === 200) {
                paragraphs.value = await response.json();
            } else {
                const errorResponse = await response.json();
                fetchingTTSError.value = errorResponse.error;
                paragraphs.value = [];
            }
        } catch (error) {
            console.log('Unexpected Error:', error);
            fetchingTTSError.value = 'Unexpected error';
            paragraphs.value = [];
        } finally {
            fetchingTTS.value = false;
        }
    };

    const fetchTTSVoices = async () => {
        fetchingTTSVoices.value = true;
        fetchingTTSVoicesError.value = null;

        try {

            if (user.value === null) {
                fetchingTTSVoices.value = false;
                fetchingTTSVoicesError.value = 'You\'re not logged in!';
                return;
            }

            const response: Response | undefined = await authStore.authorizedFetch(url + '/novels/tts/voices',
                {
                    method: 'GET',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                });

            if (response === null) {
                fetchingTTSVoices.value = false;
                fetchingTTSVoicesError.value = 'Unexpected Error';
                return;
            }

            if (response.status === 200) {
                const result = await response.json();
                console.log(result);
            } else {
                const errorResponse = await response.json();
                fetchingTTSVoicesError.value = errorResponse.error;
            }
        } catch (error) {
            console.log('Unexpected Error:', error);
            fetchingTTSVoicesError.value = 'Unexpected error';
        } finally {
            fetchingTTSVoices.value = false;
        }
    };

    return {
        paragraphs,
        fetchingTTS,
        fetchingTTSError,
        fetchingTTSVoices,
        fetchingTTSVoicesError,
        generateTTS,
        fetchTTSVoices,
    };
});

if (import.meta.hot) {
    import.meta.hot.accept(acceptHMRUpdate(useTTSStore, import.meta.hot));
}