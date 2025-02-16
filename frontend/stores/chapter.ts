import type { ErrorHandler } from '~/interfaces/ErrorHandler';
import type { HttpClient } from '~/interfaces/HttpClient';
import type { ResponseParser } from '~/interfaces/ResponseParser';
import type { ChapterService } from '~/interfaces/services/ChapterService';
import type { Chapter, ChapterSchema } from '~/schemas/Chapter';
import type { PaginatedServerResponse } from '~/schemas/PaginatedServerResponse';
import { BaseChapterService } from '~/services/ChapterService';

export const useChapterStore = defineStore('Chapter', () => {
  const runtimeConfig = useRuntimeConfig();
  const url = runtimeConfig.public.apiUrl;

  // Initialize chapter service
  const httpClient: HttpClient = new FetchHttpClient();
  const responseParser: ResponseParser = new ZodResponseParser();
  const chapterService: ChapterService = new BaseChapterService(url, httpClient, responseParser);

  // Initialize error handler
  const errorHandler: ErrorHandler = new BaseErrorHandler();

  const paginatedChapterData = shallowRef<PaginatedServerResponse<typeof ChapterSchema> | null>(null);
  const chapter: Ref<Chapter | null> = ref<Chapter | null>(null);
  const fetchingChapters = ref(true);

  const fetchChapter = async (novelUpdatesId: string, chaptNo: number) => {
    fetchingChapters.value = true;

    try {
      chapter.value = await chapterService.fetchChapter(novelUpdatesId, chaptNo);
    } catch (error) {
      errorHandler.handleError(error, { novelUpdatesId: novelUpdatesId, chapterNo: chaptNo, location: 'chapter.ts -> fetchChapter' });
      chapter.value = null;
      throw error;
    } finally {
      fetchingChapters.value = false;
    }
  };

  const fetchChapters = async (novelUpdatesId: string, page: number, limit: number): Promise<void> => {
    fetchingChapters.value = true;

    try {
      paginatedChapterData.value = await chapterService.fetchChapters(novelUpdatesId, page, limit);
    } catch (error) {
      errorHandler.handleError(error, { novelUpdatesId: novelUpdatesId, page: page, limit: limit, location: 'chapter.ts -> fetchChapters' });
      paginatedChapterData.value = null;
      throw error;
    } finally {
      fetchingChapters.value = false;
    }
  };

  return {
    chapter,
    fetchingChapters,
    paginatedChapterData,
    fetchChapter,
    fetchChapters,
  };
});

if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useChapterStore, import.meta.hot));
}
