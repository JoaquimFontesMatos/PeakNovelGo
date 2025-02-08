import type { Novel, NovelSchema } from '~/schemas/Novel';
import type { PaginatedServerResponse } from '~/schemas/PaginatedServerResponse';
import { NovelService } from '~/services/NovelService';

export const useNovelStore = defineStore('Novel', () => {
  const runtimeConfig = useRuntimeConfig();
  const url = runtimeConfig.public.apiUrl;

  const novelService = new NovelService(url);

  const novel = shallowRef<Novel | null>(null);
  const fetchingNovel = ref(true);

  const paginatedNovelsData = shallowRef<PaginatedServerResponse<typeof NovelSchema> | null>(null);

  const paginatedNovelsDataByTag = shallowRef<PaginatedServerResponse<typeof NovelSchema> | null>(null);

  const paginatedNovelsDataByAuthor = shallowRef<PaginatedServerResponse<typeof NovelSchema> | null>(null);

  const paginatedNovelsDataByGenre = shallowRef<PaginatedServerResponse<typeof NovelSchema> | null>(null);

  const fetchNovel = async (novelUpdatesId: string): Promise<void> => {
    fetchingNovel.value = true;

    try {
      novel.value = await novelService.fetchNovel(novelUpdatesId);
    } catch (error) {
      handleError(error, { novelUpdatesId: novelUpdatesId, location: 'novel.ts -> fetchNovel' });
      novel.value = null;
      throw error;
    } finally {
      fetchingNovel.value = false;
    }
  };

  const fetchNovels = async (page: number, limit: number): Promise<void> => {
    fetchingNovel.value = true;

    try {
      paginatedNovelsData.value = await novelService.fetchNovels(page, limit);
    } catch (error) {
      handleError(error, { page: page, limit: limit, location: 'novel.ts -> fetchNovels' });
      paginatedNovelsData.value = null;
      throw error;
    } finally {
      fetchingNovel.value = false;
    }
  };

  const fetchNovelsByTag = async (tag: string, page: number, limit: number): Promise<void> => {
    fetchingNovel.value = true;

    try {
      paginatedNovelsDataByTag.value = await novelService.fetchNovelsByTag(tag, page, limit);
    } catch (error) {
      handleError(error, { tag: tag, page: page, limit: limit, location: 'novel.ts -> fetchNovelsByTag' });
      paginatedNovelsDataByTag.value = null;
      throw error;
    } finally {
      fetchingNovel.value = false;
    }
  };

  const fetchNovelsByAuthor = async (author: string, page: number, limit: number): Promise<void> => {
    fetchingNovel.value = true;

    try {
      paginatedNovelsDataByAuthor.value = await novelService.fetchNovelsByAuthor(author, page, limit);
    } catch (error) {
      handleError(error, { author: author, page: page, limit: limit, location: 'novel.ts -> fetchNovelsByAuthor' });
      paginatedNovelsDataByAuthor.value = null;
      throw error;
    } finally {
      fetchingNovel.value = false;
    }
  };

  const fetchNovelsByGenre = async (genre: string, page: number, limit: number): Promise<void> => {
    fetchingNovel.value = true;

    try {
      paginatedNovelsDataByGenre.value = await novelService.fetchNovelsByGenre(genre, page, limit);
    } catch (error) {
      handleError(error, { genre: genre, page: page, limit: limit, location: 'novel.ts -> fetchNovelsByGenre' });
      paginatedNovelsDataByGenre.value = null;
      throw error;
    } finally {
      fetchingNovel.value = false;
    }
  };

  return {
    novel,
    fetchingNovel,
    paginatedNovelsData,
    paginatedNovelsDataByAuthor,
    paginatedNovelsDataByGenre,
    paginatedNovelsDataByTag,
    fetchNovel,
    fetchNovels,
    fetchNovelsByAuthor,
    fetchNovelsByGenre,
    fetchNovelsByTag,
  };
});

if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useNovelStore, import.meta.hot));
}
