import type { ErrorHandler } from '~/interfaces/ErrorHandler';
import type { HttpClient } from '~/interfaces/HttpClient';
import type { ResponseParser } from '~/interfaces/ResponseParser';
import type { NovelService } from '~/interfaces/services/NovelService';
import type { ImportedNovel } from '~/schemas/ImportedNovel';
import type { Novel, NovelSchema } from '~/schemas/Novel';
import type { PaginatedServerResponse } from '~/schemas/PaginatedServerResponse';
import { BaseNovelService } from '~/services/NovelService';

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

  const importingNovel= ref(false);

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
      $errorHandler.handleError(error, { tag: tag, page: page, limit: limit, location: 'novel.ts -> fetchNovelsByTag' });
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
      $errorHandler.handleError(error, { author: author, page: page, limit: limit, location: 'novel.ts -> fetchNovelsByAuthor' });
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
      $errorHandler.handleError(error, { genre: genre, page: page, limit: limit, location: 'novel.ts -> fetchNovelsByGenre' });
      paginatedNovelsDataByGenre.value = null;
      throw error;
    } finally {
      fetchingNovel.value = false;
    }
  };

  const importNovel = async (importedNovel: ImportedNovel): Promise<void> => {
    importingNovel.value = true;

    try {
      novel.value = await $novelService.importNovel(importedNovel);
    } catch (error) {
      $errorHandler.handleError(error, { importedNovel: importedNovel, location: 'novel.ts -> importNovel' });
      novel.value = null;
      throw error;
    } finally {
      importingNovel.value = false;
    }
  };


  const importByNovelUpdatesId = async (novelUpdatesId: string): Promise<void> => {
    importingNovel.value = true;

    try {
      novel.value = await $novelService.importByNovelUpdatesId(novelUpdatesId);
    } catch (error) {
      $errorHandler.handleError(error, { novelUpdatesId: novelUpdatesId, location: 'novel.ts -> importByNovelUpdatesId' });
      novel.value = null;
      throw error;
    } finally {
      importingNovel.value = false;
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
    fetchNovel,
    fetchNovels,
    fetchNovelsByAuthor,
    fetchNovelsByGenre,
    fetchNovelsByTag,
    importNovel,
    importByNovelUpdatesId
  };
});

if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useNovelStore, import.meta.hot));
}
