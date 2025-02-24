import { acceptHMRUpdate, defineStore } from 'pinia';
import type { ErrorHandler } from '~/interfaces/ErrorHandler';
import type { HttpClient } from '~/interfaces/HttpClient';
import type { ResponseParser } from '~/interfaces/ResponseParser';
import type { ChapterService } from '~/interfaces/services/ChapterService';
import type { Chapter, ChapterSchema } from '~/schemas/Chapter';
import type { PaginatedServerResponse } from '~/schemas/PaginatedServerResponse';
import { BaseChapterService } from '~/services/ChapterService';

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
  const chapterStatuses = ref<Record<number, string>>({}); // Track chapter import statuses
  const fetchChapter = async (novelUpdatesId: string, chaptNo: number) => {
    fetchingChapters.value = true;

    try {
      chapter.value = await $chapterService.fetchChapter(novelUpdatesId, chaptNo);
    } catch (error) {
      $errorHandler.handleError(error, { novelUpdatesId: novelUpdatesId, chapterNo: chaptNo, location: 'chapter.ts -> fetchChapter' });
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
      $errorHandler.handleError(error, { novelUpdatesId: novelUpdatesId, page: page, limit: limit, location: 'chapter.ts -> fetchChapters' });
      paginatedChapterData.value = null;
      throw error;
    } finally {
      fetchingChapters.value = false;
    }
  };

  const importChapters = async (novelUpdatesId: string) => {
    importingChapters.value = true; // Set importing state to true
    chapterStatuses.value = {}; // Reset statuses

    const eventSourceUrl = `${url}/novels/chapters/${novelUpdatesId}/scrape`;

    const eventSource = new EventSource(eventSourceUrl);

    // Debounced handler for status updates
    const throttledStatusHandler = throttle((event: MessageEvent) => {
      try {
        const statuses = JSON.parse(event.data) as Record<number, string>;
        chapterStatuses.value = { ...chapterStatuses.value, ...statuses }; // Update statuses
      } catch (error) {
        console.error('Failed to parse status update:', error);
      }
    }, 1000); // Adjust the interval (in milliseconds) as needed

    eventSource.addEventListener('status', throttledStatusHandler);

    eventSource.addEventListener('error', event => {
      console.error('EventSource failed:', event);
      eventSource.close();
      importingChapters.value = false; // Reset importing state
    });

    eventSource.addEventListener('complete', () => {
      console.log('Chapter import completed');
      eventSource.close();
      importingChapters.value = false; // Reset importing state
    });
  };

  return {
    chapter,
    fetchingChapters,
    paginatedChapterData,
    importingChapters,
    chapterStatuses,
    fetchChapter,
    fetchChapters,
    importChapters,
  };
});

if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useChapterStore, import.meta.hot));
}
