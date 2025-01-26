import type { Chapter } from '~/models/Chapter';
import type { PaginatedServerResponse } from '~/models/PaginatedServerResponse';

export const useChapterStore = defineStore('Chapter', () => {
  const runtimeConfig = useRuntimeConfig();
  const url = runtimeConfig.public.apiUrl;

  const paginatedChapterData =
    shallowRef<PaginatedServerResponse<Chapter> | null>(null);
  const chapter: Ref<Chapter | null> = ref<Chapter | null>(null);
  const fetchingChapters = ref(true);
  const chapterError = ref<string | null>(null);

  const fetchChapter = async (novelUpdatesId: string, chaptNo: number) => {
    fetchingChapters.value = true;
    chapterError.value = null;

    try {
      const response = await fetch(url + '/novels/chapters/novel/' + novelUpdatesId + '/chapter/' + chaptNo, {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
        },
      });
      if (response.status === 200) {
        chapter.value = await response.json();
      } else {
        chapter.value = null;
        const errorResponse = await response.json();

        chapterError.value = errorResponse.error;
      }
    } catch (error) {
      console.error('Unbookmarking Novel Error');
      chapterError.value = 'An unexpected error occurred. Please try again later.';
      chapter.value = null;
    } finally {
      fetchingChapters.value = false;
    }
  };

  const fetchChapters = async (
    novelUpdatesId: string,
    page: number,
    limit: number,
  ): Promise<void> => {
    fetchingChapters.value = true;
    chapterError.value = null;
    fetch(
      `${url}/novels/chapters/novel/${novelUpdatesId}/chapters?page=${page}&limit=${limit}`,
    )
      .then((response) => response.json())
      .then((data) => {
        paginatedChapterData.value = data;
      })
      .catch((error) => {
        paginatedChapterData.value = null;
        chapterError.value = error;
      })
      .finally(() => {
        fetchingChapters.value = false;
      });
  };

  return {
    chapter,
    fetchChapter,
    fetchChapters,
    fetchingChapters,
    chapterError,
    paginatedChapterData,
  };
});

if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useChapterStore, import.meta.hot));
}
