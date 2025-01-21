import type { Chapter } from "~/models/Chapter";
import type { PaginatedServerResponse } from "~/models/PaginatedServerResponse";

export const useChapterStore = defineStore("Chapter", () => {
  const runtimeConfig = useRuntimeConfig();
  const url = runtimeConfig.public.apiUrl;

  const paginatedChapterData = ref<PaginatedServerResponse<Chapter>>({
    data: [],
    limit: 0,
    page: 0,
    total: 0,
    totalPages: 0,
  });
  const fetchingChapters = ref(true);
  const chapterError = ref<string | null>(null);

  const fetchChapters = async (
    novelUpdatesId: string,
    page: number,
    limit: number
  ): Promise<void> => {
    fetchingChapters.value = true;
    chapterError.value = null;
    fetch(
      `${url}/novels/chapters/novel/${novelUpdatesId}/chapters?page=${page}&limit=${limit}`
    )
      .then((response) => response.json())
      .then((data) => {
        paginatedChapterData.value = data;
      })
      .catch((error) => {
        chapterError.value = error;
      })
      .finally(() => {
        fetchingChapters.value = false;
      });
  };

  return {
    fetchChapters,
    fetchingChapters,
    chapterError,
    paginatedChapterData,
  };
});

if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useChapterStore, import.meta.hot));
}
