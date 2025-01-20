import type { Chapter } from "~/models/Chapter";
import type { PaginatedServerResponse } from "~/models/PaginatedServerResponse";

export function useChapters(
  novelUpdatesId: string,
  page: number,
  limit: number
) {
  const runtimeConfig = useRuntimeConfig();
  const apiUrl = runtimeConfig.public.apiUrl;

  const paginatedChapterData = ref<PaginatedServerResponse<Chapter> | null>(
    null
  );
  const chapterError = ref<string | null>(null);

  async function fetchChapters(): Promise<PaginatedServerResponse<Chapter> | null> {
    if (
      paginatedChapterData.value === null ||
      paginatedChapterData.value === undefined
    ) {
      chapterError.value = "Failed to fetch chapters";
      return null;
    }
    const res = await fetch(
      `${apiUrl}/novels/chapters/novel/${novelUpdatesId}/chapters?page=${page}&limit=${limit}`
    );
    if (!res.ok) {
      chapterError.value = "Failed to fetch chapters";
      return null;
    }
    paginatedChapterData.value = await res.json();
    return paginatedChapterData.value;
  }

  return {
    paginatedChapterData,
    chapterError,
    fetchChapters,
  };
}
