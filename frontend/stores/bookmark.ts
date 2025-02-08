import type { ShallowRef } from 'vue';
import { AuthError } from '~/errors/AuthError';
import type { BookmarkedNovel, BookmarkedNovelSchema } from '~/schemas/BookmarkedNovel';
import type { PaginatedServerResponse } from '~/schemas/PaginatedServerResponse';
import { BookmarkService } from '~/services/BookmarkService';

export const useBookmarkStore = defineStore('Bookmark', () => {
  const runtimeConfig = useRuntimeConfig();
  const url: string = runtimeConfig.public.apiUrl;

  const bookmarkService = new BookmarkService(url);

  const authStore = useAuthStore();

  const { user } = storeToRefs(authStore);

  const bookmark: Ref<BookmarkedNovel | null> = ref<BookmarkedNovel | null>(null);
  const paginatedBookmarkedNovels: ShallowRef<PaginatedServerResponse<typeof BookmarkedNovelSchema> | null> = shallowRef<PaginatedServerResponse<
    typeof BookmarkedNovelSchema
  > | null>(null);
  const fetchingBookmarkedNovel: Ref<boolean | null> = ref<boolean>(true);
  const bookmarkingNovel: Ref<boolean | null> = ref<boolean>(false);
  const updatingBookmark: Ref<boolean | null> = ref<boolean>(false);
  const unbookmarkMessage: Ref<string | null> = ref<string | null>(null);

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

      bookmark.value = await bookmarkService.bookmarkNovel(novelId, user.value.ID);
    } catch (error) {
      handleError(error, { user: user, novelId: novelId, location: 'bookmark.ts -> bookmarkNovel' });
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

      bookmark.value = await bookmarkService.updateBookmark(updatedBookmark);
    } catch (error) {
      handleError(error, { user: user, updatedBookmark: updatedBookmark, location: 'bookmark.ts -> updateBookmark' });
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

      const message = await bookmarkService.unbookmarkNovel(novelId, user.value.ID);
      unbookmarkMessage.value = message.message;
      bookmark.value = null;
    } catch (error) {
      handleError(error, { user: user, novelId: novelId, location: 'bookmark.ts -> unbookmarkNovel' });
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

      bookmark.value = await bookmarkService.fetchBookmarkedNovelByUser(novelId, user.value.ID);
    } catch (error) {
      handleError(error, { user: user, novelId: novelId, location: 'bookmark.ts -> fetchBookmarkedNovelByUser' });
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

      paginatedBookmarkedNovels.value = await bookmarkService.fetchBookmarkedNovelsByUser(user.value.ID, page, limit);
    } catch (error) {
      paginatedBookmarkedNovels.value = null;
      handleError(error, { user: user, page: page, limit: limit, location: 'bookmark.ts -> fetchBookmarkedNovelsByUser' });
      throw error;
    } finally {
      fetchingBookmarkedNovel.value = false;
    }
  };

  return {
    bookmark,
    paginatedBookmarkedNovels,
    fetchingBookmarkedNovel,
    bookmarkingNovel,
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
