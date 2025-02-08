import type { Chapter, ChapterSchema } from '~/schemas/Chapter';
import type { PaginatedServerResponse } from '~/schemas/PaginatedServerResponse';
import { ChapterService } from '~/services/ChapterService';

export const useChapterStore = defineStore('Chapter', () => {
  const runtimeConfig = useRuntimeConfig();
  const url = runtimeConfig.public.apiUrl;

  const chapterService = new ChapterService(url);

  const paginatedChapterData = shallowRef<PaginatedServerResponse<typeof ChapterSchema> | null>(null);
  const chapter: Ref<Chapter | null> = ref<Chapter | null>(null);
  const fetchingChapters = ref(true);

  const fetchChapter = async (novelUpdatesId: string, chaptNo: number) => {
    fetchingChapters.value = true;

    try {
      chapter.value = await chapterService.fetchChapter(novelUpdatesId, chaptNo);
    } catch (error) {
      handleError(error, { novelUpdatesId: novelUpdatesId, chapterNo: chaptNo, location: 'chapter.ts -> fetchChapter' });
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
      handleError(error, { novelUpdatesId: novelUpdatesId, page: page, limit: limit, location: 'chapter.ts -> fetchChapters' });
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
