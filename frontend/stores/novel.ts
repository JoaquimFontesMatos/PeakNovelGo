import type {ErrorHandler} from '~/interfaces/ErrorHandler';
import type {HttpClient} from '~/interfaces/HttpClient';
import type {ResponseParser} from '~/interfaces/ResponseParser';
import type {NovelService} from '~/interfaces/services/NovelService';
import type {Novel, NovelSchema} from '~/schemas/Novel';
import type {PaginatedServerResponse} from '~/schemas/PaginatedServerResponse';
import {BaseNovelService} from '~/services/NovelService';

export const useNovelStore = defineStore('Novel', () => {
    const runtimeConfig = useRuntimeConfig();
    const url: string = runtimeConfig.public.apiUrl;
    const httpClient: HttpClient = new FetchHttpClient(useAuthStore());
    const responseParser: ResponseParser = new ZodResponseParser();
    const $novelService: NovelService = new BaseNovelService(url, httpClient, responseParser);
    const $errorHandler: ErrorHandler = new BaseErrorHandler();

    const novel = shallowRef<Novel | null>(null);
    const fetchingNovel = ref(true);

    const paginatedNovelsData = shallowRef<PaginatedServerResponse<typeof NovelSchema> | null>(null);

    const paginatedNovelsDataByTag = shallowRef<PaginatedServerResponse<typeof NovelSchema> | null>(null);

    const paginatedNovelsDataByAuthor = shallowRef<PaginatedServerResponse<typeof NovelSchema> | null>(null);

    const paginatedNovelsDataByGenre = shallowRef<PaginatedServerResponse<typeof NovelSchema> | null>(null);

    const importingNovel = ref(false);

    const updatingNovels = ref(false);
    const novelStatuses = ref<Record<string, string>>({});

    const fetchNovel = async (novelUpdatesId: string): Promise<void> => {
        fetchingNovel.value = true;

        try {
            novel.value = await $novelService.fetchNovel(novelUpdatesId);
        } catch (error) {
            $errorHandler.handleError(error, {novelUpdatesId: novelUpdatesId, location: 'novel.ts -> fetchNovel'});
            novel.value = null;
            throw error;
        } finally {
            fetchingNovel.value = false;
        }
    };

    const fetchNovels = async (page: number, limit: number): Promise<void> => {
        fetchingNovel.value = true;

        try {
            paginatedNovelsData.value = await $novelService.fetchNovels(page, limit);
        } catch (error) {
            $errorHandler.handleError(error, {page: page, limit: limit, location: 'novel.ts -> fetchNovels'});
            paginatedNovelsData.value = null;
            throw error;
        } finally {
            fetchingNovel.value = false;
        }
    };

    const fetchNovelsByTag = async (tag: string, page: number, limit: number): Promise<void> => {
        fetchingNovel.value = true;

        try {
            paginatedNovelsDataByTag.value = await $novelService.fetchNovelsByTag(tag, page, limit);
        } catch (error) {
            $errorHandler.handleError(error, {
                tag: tag,
                page: page,
                limit: limit,
                location: 'novel.ts -> fetchNovelsByTag'
            });
            paginatedNovelsDataByTag.value = null;
            throw error;
        } finally {
            fetchingNovel.value = false;
        }
    };

    const fetchNovelsByAuthor = async (author: string, page: number, limit: number): Promise<void> => {
        fetchingNovel.value = true;

        try {
            paginatedNovelsDataByAuthor.value = await $novelService.fetchNovelsByAuthor(author, page, limit);
        } catch (error) {
            $errorHandler.handleError(error, {
                author: author,
                page: page,
                limit: limit,
                location: 'novel.ts -> fetchNovelsByAuthor'
            });
            paginatedNovelsDataByAuthor.value = null;
            throw error;
        } finally {
            fetchingNovel.value = false;
        }
    };

    const fetchNovelsByGenre = async (genre: string, page: number, limit: number): Promise<void> => {
        fetchingNovel.value = true;

        try {
            paginatedNovelsDataByGenre.value = await $novelService.fetchNovelsByGenre(genre, page, limit);
        } catch (error) {
            $errorHandler.handleError(error, {
                genre: genre,
                page: page,
                limit: limit,
                location: 'novel.ts -> fetchNovelsByGenre'
            });
            paginatedNovelsDataByGenre.value = null;
            throw error;
        } finally {
            fetchingNovel.value = false;
        }
    };

    const importByNovelUpdatesId = async (novelUpdatesId: string): Promise<void> => {
        importingNovel.value = true;

        try {
            novel.value = await $novelService.importByNovelUpdatesId(novelUpdatesId);
        } catch (error) {
            $errorHandler.handleError(error, {
                novelUpdatesId: novelUpdatesId,
                location: 'novel.ts -> importByNovelUpdatesId'
            });
            novel.value = null;
            throw error;
        } finally {
            importingNovel.value = false;
        }
    };

    const batchUpdateNovels = async () => {
        updatingNovels.value = true;
        novelStatuses.value = {};

        const eventSourceUrl = `${url}/novels/update`;

        try {
            const response = await httpClient.authorizedRequest(eventSourceUrl, {
                headers: { Accept: 'text/event-stream' },
            });

            if (!response.ok) {
                throw new Error(`SSE connection failed: ${response.statusText}`);
            }

            const reader = response.body?.getReader();
            const decoder = new TextDecoder();

            if (!reader) {
                throw new Error("Failed to read SSE stream");
            }

            let buffer = '';
            while (true) {
                const { done, value } = await reader.read();
                if (done) break;

                buffer += decoder.decode(value, { stream: true });
                const events = buffer.split('\n\n');

                // Process all complete events
                for (let i = 0; i < events.length - 1; i++) {
                    const event = events[i].trim();
                    if (!event) continue;

                    let eventType = '';
                    let eventData = '';

                    event.split('\n').forEach(line => {
                        if (line.startsWith('event:')) {
                            eventType = line.substring(6).trim();
                        } else if (line.startsWith('data:')) {
                            eventData = line.substring(5).trim();
                        }
                    });

                    if (eventType === 'status' && eventData) {
                        try {
                            novelStatuses.value = JSON.parse(eventData);
                        } catch (error) {
                            console.error("Failed to parse status update:", error);
                        }
                    } else if (eventType === 'error' && eventData) {
                        useToastStore().addToast(eventData, 'error', 'novel');
                    } else if (eventType === 'complete' && eventData) {
                        useToastStore().addToast(eventData, 'success', 'novel');
                    }
                }

                buffer = events[events.length - 1]; // Keep incomplete data
            }
        } catch (error) {
            console.error("SSE stream error:", error);
            useToastStore().addToast("Failed to update novels", 'error', 'novel');
        } finally {
            updatingNovels.value = false;
        }
    };

    return {
        novel,
        fetchingNovel,
        paginatedNovelsData,
        paginatedNovelsDataByAuthor,
        paginatedNovelsDataByGenre,
        paginatedNovelsDataByTag,
        importingNovel,
        updatingNovels,
        novelStatuses,
        fetchNovel,
        fetchNovels,
        fetchNovelsByAuthor,
        fetchNovelsByGenre,
        fetchNovelsByTag,
        importByNovelUpdatesId,
        batchUpdateNovels
    };
});

if (import.meta.hot) {
    import.meta.hot.accept(acceptHMRUpdate(useNovelStore, import.meta.hot));
}
