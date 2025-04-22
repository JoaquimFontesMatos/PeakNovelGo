import type { ErrorHandler } from '~/interfaces/ErrorHandler';
import type { HttpClient } from '~/interfaces/HttpClient';
import type { ResponseParser } from '~/interfaces/ResponseParser';
import type { NovelService } from '~/interfaces/services/NovelService';
import type { Novel, NovelSchema } from '~/schemas/Novel';
import type { PaginatedServerResponse } from '~/schemas/PaginatedServerResponse';
import { BaseNovelService } from '~/services/NovelService';
import { useIndexedDB } from '~/composables/useInitCacheDB';

const RECENTLY_VISITED_NOVEL_STORE = 'recentlyVisitedNovels';

export const useNovelStore = defineStore('Novel', () => {
    const runtimeConfig = useRuntimeConfig();
    const url: string = runtimeConfig.public.apiUrl;
    const httpClient: HttpClient = new FetchHttpClient(useAuthStore());
    const responseParser: ResponseParser = new ZodResponseParser();
    const $novelService: NovelService = new BaseNovelService(url, httpClient, responseParser);
    const $errorHandler: ErrorHandler = new BaseErrorHandler();

    const { initDB } = useIndexedDB();

    const cachedRecentlyVisitedNovels = ref<Novel[]>([]);

    const novel = shallowRef<Novel | null>(null);
    const fetchingNovel = ref(false);

    const paginatedNovelsData = shallowRef<PaginatedServerResponse<typeof NovelSchema> | null>(null);

    const paginatedNovelsDataByTag = shallowRef<PaginatedServerResponse<typeof NovelSchema> | null>(null);

    const paginatedNovelsDataByAuthor = shallowRef<PaginatedServerResponse<typeof NovelSchema> | null>(null);

    const paginatedNovelsDataByGenre = shallowRef<PaginatedServerResponse<typeof NovelSchema> | null>(null);

    const importingNovel = ref(false);

    const updatingNovels = ref(false);
    const novelStatuses = ref<Record<string, string>>({});

    const cacheRecentlyVisitedNovel = async () => {
        try {
            const dbInstance = await initDB();

            if (!novel.value) {
                return;
            }

            const novelUpdatesId = novel.value.novelUpdatesId;

            if (await getCachedNovel(novelUpdatesId)) {
                return;
            }

            const cacheKey = `${novelUpdatesId}`;

            const transaction = dbInstance.transaction(RECENTLY_VISITED_NOVEL_STORE, 'readonly');
            const store = transaction.objectStore(RECENTLY_VISITED_NOVEL_STORE);
            const request = store.get(cacheKey);

            request.onsuccess = event => {
                const cachedNovel = (event.target as IDBRequest<Novel>).result;
                if (cachedNovel) {
                    console.log(`Novel ${novelUpdatesId} loaded from cache`);
                    cachedRecentlyVisitedNovels.value.push(cachedNovel);
                    return;
                }

                const cacheTransaction = dbInstance.transaction(RECENTLY_VISITED_NOVEL_STORE, 'readwrite');
                const cacheStore = cacheTransaction.objectStore(RECENTLY_VISITED_NOVEL_STORE);
                cacheStore.put({ ...novel.value, cacheKey });
                console.log(`Novel ${novelUpdatesId} cached`);
            };
            request.onerror = event => {
                console.error('Error getting novel from IndexedDB:', (event.target as IDBRequest).error);
            };
        } catch (error) {
            console.error('IndexedDB initialization error:', error);
        }
    };

    const getCachedNovel = async (novelUpdatesId: string): Promise<Novel | null> => {
        try {
            const dbInstance = await initDB();

            const cacheKey = `${novelUpdatesId}`;

            const transaction = dbInstance.transaction(RECENTLY_VISITED_NOVEL_STORE, 'readonly');
            const store = transaction.objectStore(RECENTLY_VISITED_NOVEL_STORE);
            const request = store.get(cacheKey);

            return new Promise((resolve, reject) => {
                request.onsuccess = event => {
                    resolve((event.target as IDBRequest<Novel | undefined>).result || null);
                };

                request.onerror = event => {
                    console.error('Error getting novel from IndexedDB:', (event.target as IDBRequest).error);
                    reject(null);
                };
            });
        } catch (error) {
            console.error('Error accessing IndexedDB:', error);
            return null;
        }
    };

    const getCachedNovels = async (): Promise<void> => {
        try {
            const dbInstance = await initDB();

            const transaction = dbInstance.transaction(RECENTLY_VISITED_NOVEL_STORE, 'readonly');
            const store = transaction.objectStore(RECENTLY_VISITED_NOVEL_STORE);

            // Method 1: Using getAll() (Simpler, but potentially less efficient for large datasets)
            const request = store.getAll();
            return new Promise((resolve, reject) => {
                request.onsuccess = event => {
                    cachedRecentlyVisitedNovels.value = (event.target as IDBRequest<Novel[]>).result; // Resolve with the array of novels
                };
                request.onerror = event => {
                    console.error('Error getting novels from IndexedDB:', (event.target as IDBRequest).error);
                    reject([]); // Reject with an empty array if an error occurs
                };
            });
        } catch (error) {
            console.error('Error accessing IndexedDB:', error);
        }
    };

    const fetchNovel = async (novelUpdatesId: string): Promise<void> => {
        fetchingNovel.value = true;

        try {
            novel.value = await $novelService.fetchNovel(novelUpdatesId);
        } catch (error) {
            $errorHandler.handleError(error, { novelUpdatesId: novelUpdatesId, location: 'novel.ts -> fetchNovel' });
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
            $errorHandler.handleError(error, { page: page, limit: limit, location: 'novel.ts -> fetchNovels' });
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
                location: 'novel.ts -> fetchNovelsByTag',
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
                location: 'novel.ts -> fetchNovelsByAuthor',
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
                location: 'novel.ts -> fetchNovelsByGenre',
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
                location: 'novel.ts -> importByNovelUpdatesId',
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
                throw new Error('Failed to read SSE stream');
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
                            console.error('Failed to parse status update:', error);
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
            console.error('SSE stream error:', error);
            useToastStore().addToast('Failed to update novels', 'error', 'novel');
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
        cachedRecentlyVisitedNovels,
        getCachedNovel,
        getCachedNovels,
        cacheRecentlyVisitedNovel,
        fetchNovel,
        fetchNovels,
        fetchNovelsByAuthor,
        fetchNovelsByGenre,
        fetchNovelsByTag,
        importByNovelUpdatesId,
        batchUpdateNovels,
    };
});

if (import.meta.hot) {
    import.meta.hot.accept(acceptHMRUpdate(useNovelStore, import.meta.hot));
}
