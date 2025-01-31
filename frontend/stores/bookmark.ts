import type { BookmarkedNovel } from '~/models/BookmarkedNovel';
import type { PaginatedServerResponse } from '~/models/PaginatedServerResponse';
import type { ShallowRef } from 'vue';

export const useBookmarkStore = defineStore('Bookmark', () => {
    const runtimeConfig = useRuntimeConfig();
    const url: string = runtimeConfig.public.apiUrl;
    const authStore = useAuthStore();

    const { user } = storeToRefs(authStore);

    const bookmark: Ref<BookmarkedNovel | null> = ref<BookmarkedNovel | null>(null);
    const paginatedBookmarkedNovels: ShallowRef<PaginatedServerResponse<BookmarkedNovel> | null> =
        shallowRef<PaginatedServerResponse<BookmarkedNovel> | null>(null);
    const fetchingBookmarkedNovel: Ref<boolean | null> = ref<boolean>(true);
    const bookmarkedNovelError: Ref<string | null> = ref<string | null>(null);
    const bookmarkNovelError: Ref<string | null> = ref<string | null>(null);
    const bookmarkingNovel: Ref<boolean | null> = ref<boolean>(false);
    const updateBookmarkError: Ref<string | null> = ref<string | null>(null);
    const updatingBookmark: Ref<boolean | null> = ref<boolean>(false);

    const bookmarkNovel = async (novelId: number): Promise<void> => {
        bookmarkingNovel.value = true;
        bookmarkNovelError.value = null;


        if (user.value === null) {
            bookmarkingNovel.value = false;
            bookmarkNovelError.value = 'You\'re not logged in!';
            bookmark.value = null;
            return;
        }

        const createBookmarkedNovel = {
            novelId: novelId,
            userId: user.value.ID,
            score: 0,
            status: 'Plan to Read',
            currentChapter: 1,
        };

        try {
            const response: Response | undefined = await authStore.authorizedFetch(`${url}/novels/bookmarked`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(createBookmarkedNovel),
            });

            if (response === undefined) {
                bookmark.value = null;
                bookmarkingNovel.value = false;
                return;
            }
            if (response.status === 200) {
                bookmark.value = await response.json();
            } else {
                const errorResponse = await response.json();

                bookmarkNovelError.value = errorResponse.error;
                bookmark.value = null;
            }
        } catch (error) {
            console.error('Bookmarking Novel error:', error);
            bookmark.value = null;
        } finally {
            bookmarkingNovel.value = false;
        }
    };

    const updateBookmark = async (updatedBookmark: BookmarkedNovel): Promise<void> => {
        updatingBookmark.value = true;
        updateBookmarkError.value = null;

        try {
            const response: Response | undefined = await authStore.authorizedFetch(
                `${url}/novels/bookmarked`, {
                    method: 'PUT',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(updatedBookmark),
                });

            if (response === undefined) {
                updatingBookmark.value = false;
                return;
            }

            if (response.status === 200) {
                bookmark.value = await response.json();
            } else {
                const errorResponse = await response.json();

                updateBookmarkError.value = errorResponse.error;
            }
        } catch (error) {
            console.error('Update bookmark error:', error);
        } finally {
            updatingBookmark.value = false;
        }
    };

    const unbookmarkNovel = async (novelId: number): Promise<void> => {
        bookmarkingNovel.value = true;
        bookmarkNovelError.value = null;

        if (user.value === null) {
            fetchingBookmarkedNovel.value = false;
            bookmarkNovelError.value = 'You\'re not logged in!';
            return;
        }

        try {
            const response = await authStore.authorizedFetch(
                `${url}/novels/bookmarked/user/${user.value.ID}/novel/${novelId}`, {
                    method: 'DELETE',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                },
            );

            if (response === undefined) {
                fetchingBookmarkedNovel.value = false;
                return;
            }

            if (response.status === 200) {
                bookmark.value = null;
            } else {
                const errorResponse = await response.json();

                updateBookmarkError.value = errorResponse.error;
            }
        } catch (error) {
            console.error('Unbookmarking Novel Error');
        } finally {
            bookmarkingNovel.value = false;
        }
    };

    const fetchBookmarkedNovelByUser = async (novelId: string): Promise<void> => {
        fetchingBookmarkedNovel.value = true;
        bookmarkedNovelError.value = null;

        try {
            if (user.value === null) {
                fetchingBookmarkedNovel.value = false;
                bookmarkedNovelError.value = 'You\'re not logged in!';

                bookmark.value = null;

                return;
            }

            const response = await authStore.authorizedFetch(`${url}/novels/bookmarked/user/${user.value.ID}/novel/${novelId}`, {
                method: 'GET',
                headers: {
                    'Content-Type': 'application/json',
                },
            });
            if (response === undefined) {
                bookmark.value = null;
                bookmarkedNovelError.value = 'An unexpected error occurred. Please try again later.';

                fetchingBookmarkedNovel.value = false;

                return;
            }

            if (response.status === 200) {
                bookmark.value = await response.json();
            } else {
                bookmark.value = null;
                const errorResponse = await response.json();

                bookmarkedNovelError.value = errorResponse.error;
            }
        } catch (error) {
            console.error('Unbookmarking Novel Error');
            bookmarkedNovelError.value = 'An unexpected error occurred. Please try again later.';
            bookmark.value = null;
        } finally {
            fetchingBookmarkedNovel.value = false;
        }
    };

    const fetchBookmarkedNovelsByUser = async (
        page: number,
        limit: number,
    ): Promise<void> => {
        fetchingBookmarkedNovel.value = true;
        bookmarkedNovelError.value = null;

        try {
            if (user.value === null) {
                fetchingBookmarkedNovel.value = false;
                bookmarkedNovelError.value = 'You\'re not logged in!';
                paginatedBookmarkedNovels.value = null;
                return;
            }

            const response: Response | undefined = await authStore.authorizedFetch(
                `${url}/novels/bookmarked/user/${user.value.ID}?page=${page}&limit=${limit}`, {
                    method: 'GET',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                },
            );

            if (response === undefined) {
                paginatedBookmarkedNovels.value = null;
            }

            if (response.status === 200) {
                paginatedBookmarkedNovels.value = await response.json();
            } else {
                paginatedBookmarkedNovels.value = null;
            }
        } catch (error) {
            console.error('Unbookmarking Novel Error');
            paginatedBookmarkedNovels.value = null;
        } finally {
            fetchingBookmarkedNovel.value = false;
        }
    };

    return {
        bookmark,
        paginatedBookmarkedNovels,
        fetchingBookmarkedNovel,
        bookmarkedNovelError,
        bookmarkNovelError,
        bookmarkingNovel,
        updateBookmarkError,
        updatingBookmark,
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