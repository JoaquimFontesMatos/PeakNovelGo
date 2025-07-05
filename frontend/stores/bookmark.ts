import type { ShallowRef } from 'vue';
import { AuthError } from '~/errors/AuthError';
import type { ErrorHandler } from '~/interfaces/ErrorHandler';
import type { HttpClient } from '~/interfaces/HttpClient';
import type { ResponseParser } from '~/interfaces/ResponseParser';
import type { BookmarkService } from '~/interfaces/services/BookmarkService';
import type { BookmarkedNovel, BookmarkedNovelSchema } from '~/schemas/BookmarkedNovel';
import type { PaginatedServerResponse } from '~/schemas/PaginatedServerResponse';
import { BaseBookmarkService } from '~/services/BookmarkService';
import type { NovelSchema } from '~/schemas/Novel';
import { useIndexedDB } from '~/composables/useInitCacheDB';

const BOOKMARK_STORE = 'bookmarks';

export const useBookmarkStore = defineStore('Bookmark', () => {
    const runtimeConfig = useRuntimeConfig();
    const url: string = runtimeConfig.public.apiUrl;
    const httpClient: HttpClient = new FetchHttpClient(useAuthStore());
    const responseParser: ResponseParser = new ZodResponseParser();
    const $bookmarkService: BookmarkService = new BaseBookmarkService(url, httpClient, responseParser);
    const $errorHandler: ErrorHandler = new BaseErrorHandler();
    const { initDB } = useIndexedDB();

    const authStore = useAuthStore();

    const { user } = storeToRefs(authStore);

    const bookmark: Ref<BookmarkedNovel | null> = ref<BookmarkedNovel | null>(null);

    const cachedBookmarks = ref<BookmarkedNovel[]>([]);

    const paginatedBookmarkedNovels: ShallowRef<PaginatedServerResponse<typeof NovelSchema> | null> = shallowRef<PaginatedServerResponse<
        typeof NovelSchema
    > | null>(null);
    const fetchingBookmarkedNovel: Ref<boolean | null> = ref<boolean>(false);
    const fetchingBookmarkedNovels: Ref<boolean | null> = ref<boolean>(false);

    const bookmarkingNovel: Ref<boolean | null> = ref<boolean>(false);
    const updatingBookmark: Ref<boolean | null> = ref<boolean>(false);
    const unbookmarkMessage: Ref<string | null> = ref<string | null>(null);

    const cacheBookmark = async () => {
        try {
            const dbInstance = await initDB();

            if (!bookmark.value) {
                return;
            }

            const novelId = bookmark.value.novelId;

            if (await getCachedBookmark(novelId)) {
                return;
            }

            const cacheKey = `${novelId}`;

            const transaction = dbInstance.transaction(BOOKMARK_STORE, 'readonly');
            const store = transaction.objectStore(BOOKMARK_STORE);
            const request = store.get(cacheKey);

            request.onsuccess = event => {
                const cachedBookmark = (event.target as IDBRequest<BookmarkedNovel>).result;
                if (cachedBookmark) {
                    console.log(`Novel ${novelId} loaded from cache`);
                    cachedBookmarks.value.push(cachedBookmark);
                    return;
                }

                const cacheTransaction = dbInstance.transaction(BOOKMARK_STORE, 'readwrite');
                const cacheStore = cacheTransaction.objectStore(BOOKMARK_STORE);
                cacheStore.put({ ...bookmark.value, cacheKey });
                console.log(`Bookmark ${novelId} cached`);
            };
            request.onerror = event => {
                console.error('Error getting bookmark from IndexedDB:', (event.target as IDBRequest).error);
            };
        } catch (error) {
            console.error('IndexedDB initialization error:', error);
        }
    };

    const getCachedBookmark = async (novelId: number): Promise<BookmarkedNovel | null> => {
        try {
            const dbInstance = await initDB();

            const cacheKey = `${novelId}`;

            const transaction = dbInstance.transaction(BOOKMARK_STORE, 'readonly');
            const store = transaction.objectStore(BOOKMARK_STORE);
            const request = store.get(cacheKey);

            return new Promise((resolve, reject) => {
                request.onsuccess = event => {
                    resolve((event.target as IDBRequest<BookmarkedNovel | undefined>).result || null);
                };

                request.onerror = event => {
                    console.error('Error getting bookmark from IndexedDB:', (event.target as IDBRequest).error);
                    reject(null);
                };
            });
        } catch (error) {
            console.error('Error accessing IndexedDB:', error);
            return null;
        }
    };

    const bookmarkNovel = async (novelId: number): Promise<void> => {
        bookmarkingNovel.value = true;

        try {
            if (user.value === null) {
                throw new AuthError({
                    type: 'UNAUTHORIZED_ERROR',
                    message: "You're not logged in!",
                    cause: 'User tried to bookmark novel without being logged in.',
                });
            }

            bookmark.value = await $bookmarkService.bookmarkNovel(novelId, user.value.ID);
            await cacheBookmark();
        } catch (error) {
            $errorHandler.handleError(error, {
                user: user,
                novelId: novelId,
                location: 'bookmark.ts -> bookmarkNovel',
            });
            bookmark.value = null;
            throw error;
        } finally {
            bookmarkingNovel.value = false;
        }
    };

    const updateBookmark = async (updatedBookmark: BookmarkedNovel): Promise<void> => {
        updatingBookmark.value = true;

        try {
            if (user.value === null) {
                throw new AuthError({
                    type: 'UNAUTHORIZED_ERROR',
                    message: "You're not logged in!",
                    cause: 'User tried to update bookmark without being logged in.',
                });
            }

            bookmark.value = await $bookmarkService.updateBookmark(updatedBookmark);
            await cacheBookmark();
        } catch (error) {
            $errorHandler.handleError(error, {
                user: user,
                updatedBookmark: updatedBookmark,
                location: 'bookmark.ts -> updateBookmark',
            });
            throw error;
        } finally {
            updatingBookmark.value = false;
        }
    };

    const unbookmarkNovel = async (novelId: number): Promise<void> => {
        bookmarkingNovel.value = true;
        unbookmarkMessage.value = null;

        try {
            if (user.value === null) {
                throw new AuthError({
                    type: 'UNAUTHORIZED_ERROR',
                    message: "You're not logged in!",
                    cause: 'User tried to unbookmark without being logged in.',
                });
            }

            const message = await $bookmarkService.unbookmarkNovel(novelId, user.value.ID);
            unbookmarkMessage.value = message.message;
            bookmark.value = null;
        } catch (error) {
            $errorHandler.handleError(error, {
                user: user,
                novelId: novelId,
                location: 'bookmark.ts -> unbookmarkNovel',
            });
            throw error;
        } finally {
            bookmarkingNovel.value = false;
        }
    };

    const fetchBookmarkedNovelByUser = async (novelId: string): Promise<void> => {
        fetchingBookmarkedNovel.value = true;

        try {
            if (user.value === null) {
                throw new AuthError({
                    type: 'UNAUTHORIZED_ERROR',
                    message: "You're not logged in!",
                    cause: 'User tried to fetch bookmark without being logged in.',
                });
            }

            bookmark.value = await $bookmarkService.fetchBookmarkedNovelByUser(novelId, user.value.ID);
        } catch (error) {
            $errorHandler.handleError(error, {
                user: user,
                novelId: novelId,
                location: 'bookmark.ts -> fetchBookmarkedNovelByUser',
            });
            bookmark.value = null;
            throw error;
        } finally {
            fetchingBookmarkedNovel.value = false;
        }
    };

    const fetchBookmarkedNovelsByUser = async (page: number, limit: number): Promise<void> => {
        fetchingBookmarkedNovel.value = true;

        try {
            if (user.value === null) {
                throw new AuthError({
                    type: 'UNAUTHORIZED_ERROR',
                    message: "You're not logged in!",
                    cause: 'User tried to fetch bookmarks without being logged in.',
                });
            }

            paginatedBookmarkedNovels.value = await $bookmarkService.fetchBookmarkedNovelsByUser(user.value.ID, page, limit);
        } catch (error) {
            paginatedBookmarkedNovels.value = null;
            $errorHandler.handleError(error, {
                user: user,
                page: page,
                limit: limit,
                location: 'bookmark.ts -> fetchBookmarkedNovelsByUser',
            });
            throw error;
        } finally {
            fetchingBookmarkedNovel.value = false;
        }
    };

    return {
        bookmark,
        paginatedBookmarkedNovels,
        fetchingBookmarkedNovel,
        fetchingBookmarkedNovels,
        bookmarkingNovel,
        updatingBookmark,
        cachedBookmarks,
        cacheBookmark,
        getCachedBookmark,
        bookmarkNovel,
        unbookmarkNovel,
        updateBookmark,
        fetchBookmarkedNovelByUser,
        fetchBookmarkedNovelsByUser,
    };
});

if (import.meta.hot) {
    import.meta.hot.accept(acceptHMRUpdate(useBookmarkStore, import.meta.hot));
}
