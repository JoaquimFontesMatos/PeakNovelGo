import {acceptHMRUpdate, defineStore} from 'pinia';
import type {ErrorHandler} from '~/interfaces/ErrorHandler';
import type {HttpClient} from '~/interfaces/HttpClient';
import type {ResponseParser} from '~/interfaces/ResponseParser';
import type {ChapterService} from '~/interfaces/services/ChapterService';
import type {Chapter, ChapterSchema} from '~/schemas/Chapter';
import type {PaginatedServerResponse} from '~/schemas/PaginatedServerResponse';
import {BaseChapterService} from '~/services/ChapterService';

export const useChapterStore = defineStore('Chapter', () => {
    const runtimeConfig = useRuntimeConfig();
    const url: string = runtimeConfig.public.apiUrl;
    const httpClient: HttpClient = new FetchHttpClient(useAuthStore());
    const responseParser: ResponseParser = new ZodResponseParser();
    const $chapterService: ChapterService = new BaseChapterService(url, httpClient, responseParser);
    const $errorHandler: ErrorHandler = new BaseErrorHandler();

    const paginatedChapterData = shallowRef<PaginatedServerResponse<typeof ChapterSchema> | null>(null);
    const chapter: Ref<Chapter | null> = ref<Chapter | null>(null);
    const fetchingChapters = ref(true);

    const importingChapters = ref(false);
    let currentNovelUpdatesId: Ref<string | null> = ref(null);
    const chapterStatuses = ref<Record<string, string>>({});

    const novelProgress = computed((): number => {
        if (!chapter.value || !paginatedChapterData.value) return 0.0

        const x = chapter.value?.chapterNo ?? 1
        const y = paginatedChapterData.value?.total ?? 1

        return parseFloat(((x / y) * 100).toFixed(2))
    })

    const fetchChapter = async (novelUpdatesId: string, chaptNo: number) => {
        fetchingChapters.value = true;

        try {
            chapter.value = await $chapterService.fetchChapter(novelUpdatesId, chaptNo);
        } catch (error) {
            $errorHandler.handleError(error, {
                novelUpdatesId: novelUpdatesId,
                chapterNo: chaptNo,
                location: 'chapter.ts -> fetchChapter'
            });
            chapter.value = null;
            throw error;
        } finally {
            fetchingChapters.value = false;
        }
    };

    const fetchChapters = async (novelUpdatesId: string, page: number, limit: number): Promise<void> => {
        fetchingChapters.value = true;

        try {
            paginatedChapterData.value = await $chapterService.fetchChapters(novelUpdatesId, page, limit);
        } catch (error) {
            $errorHandler.handleError(error, {
                novelUpdatesId: novelUpdatesId,
                page: page,
                limit: limit,
                location: 'chapter.ts -> fetchChapters'
            });
            paginatedChapterData.value = null;
            throw error;
        } finally {
            fetchingChapters.value = false;
        }
    };

    const updateChapterStatuses = httpClient.throttleWithFlush((statuses: Record<string, string>) => {
        chapterStatuses.value = {...chapterStatuses.value, ...statuses};
    }, 100)


    const importChapters = async (novelUpdatesId: string) => {
        if (importingChapters.value && currentNovelUpdatesId.value === novelUpdatesId) {
            console.warn(`Import already in progress for novelUpdatesId: ${novelUpdatesId}`);
            return;
        }

        if (importingChapters.value) {
            console.warn('Import already in progress.');
            return;
        }

        currentNovelUpdatesId.value = novelUpdatesId; // Track the current id.
        importingChapters.value = true; // Set importing state to true
        chapterStatuses.value = {}; // Reset statuses

        const eventSourceUrl = `${url}/novels/chapters/${novelUpdatesId}/scrape`;

        try {
            // Use your authorizedRequest function to fetch the SSE stream
            const response = await httpClient.authorizedRequest(eventSourceUrl, {
                headers: {
                    Accept: 'text/event-stream', // Ensure the server knows you want SSE
                },
            });

            if (!response.ok) {
                console.error('Failed to establish SSE connection:', response.statusText);
                importingChapters.value = false; // Reset importing state
                currentNovelUpdatesId.value = null; // Reset after failure
                return;
            }

            // Read the stream
            const reader = response.body?.getReader();
            const decoder = new TextDecoder();

            if (!reader) {
                console.error('Failed to read SSE stream');
                importingChapters.value = false; // Reset importing state
                currentNovelUpdatesId.value = null; // Reset after failure
                return;
            }

            const processStream = async () => {
                try {
                    let buffer = ''; // Buffer to accumulate incomplete data

                    while (true) {
                        const {done, value} = await reader.read();

                        if (done) {
                            console.log('SSE stream closed');
                            // Flush the debounced function to apply the last updates
                            updateChapterStatuses.flush();
                            importingChapters.value = false;
                            currentNovelUpdatesId.value = null;
                            return;
                        }

                        const chunk = decoder.decode(value); // Decode the chunk
                        buffer += chunk; // Append the chunk to the buffer

                        // Split the buffer by double newlines (SSE events are separated by \n\n)
                        const events = buffer.split('\n\n');

                        // Process all complete events (leave the last one in the buffer if incomplete)
                        for (let i = 0; i < events.length - 1; i++) {
                            const event = events[i].trim(); // Remove any leading/trailing whitespace

                            if (!event) {
                                // Skip empty events
                                continue;
                            }

                            // Split the event into individual lines
                            const lines = event.split('\n');

                            let eventType = ''; // To store the event type (e.g., "status")
                            let eventData = ''; // To store the JSON data

                            // Process each line in the event
                            for (const line of lines) {
                                if (line.startsWith('event:')) {
                                    // Extract the event type
                                    eventType = line.replace('event:', '').trim();
                                } else if (line.startsWith('data:')) {
                                    // Extract the JSON data
                                    eventData = line.replace('data:', '').trim();
                                }
                            }

                            // Only process if the event type is "status" and data is present
                            if (eventType === 'status' && eventData) {
                                try {
                                    const statuses = JSON.parse(eventData) as Record<string, string>;

                                    updateChapterStatuses(statuses);
                                } catch (error) {
                                    console.error('Failed to parse status update:', error);
                                    console.error('Invalid JSON data:', eventData); // Log the invalid data
                                }
                            } else if (eventType === 'error' && eventData) {
                                useToastStore().addToast(eventData, 'error', 'novel');
                            } else if (eventType === 'complete' && eventData) {
                                useToastStore().addToast(eventData, 'success', 'novel');
                            }
                        }

                        // Keep the last (possibly incomplete) event in the buffer
                        buffer = events[events.length - 1];
                    }
                } catch (streamError) {
                    console.error('Error processing SSE stream:', streamError);
                } finally {
                    updateChapterStatuses.flush();
                    importingChapters.value = false;
                    currentNovelUpdatesId.value = null; // Reset after completion or error
                    reader.releaseLock();  // Release the lock, important for cleaning up resources
                }
            };

            processStream();
        } catch (error) {
            console.error('Error during importChapters:', error);
        } finally {
            updateChapterStatuses.flush();
            importingChapters.value = false;
            currentNovelUpdatesId.value = null; // Ensure it's reset even if other parts fail
        }
    };

    return {
        chapter,
        fetchingChapters,
        paginatedChapterData,
        importingChapters,
        chapterStatuses,
        novelProgress,
        fetchChapter,
        fetchChapters,
        importChapters,
    };
});

if (import.meta.hot) {
    import.meta.hot.accept(acceptHMRUpdate(useChapterStore, import.meta.hot));
}
