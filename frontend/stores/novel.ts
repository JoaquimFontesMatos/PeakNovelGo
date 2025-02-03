import type { Novel } from '~/models/Novel';
import type { PaginatedServerResponse } from '~/models/PaginatedServerResponse';

export const useNovelStore = defineStore('Novel', () => {
  const runtimeConfig = useRuntimeConfig();
  const url = runtimeConfig.public.apiUrl;

  const novel = shallowRef<Novel | null>(null);
  const fetchingNovel = ref(true);
  const novelError = ref<string | null>(null);

  const paginatedNovelsData = ref<PaginatedServerResponse<Novel>>({
    data: [],
    limit: 0,
    page: 0,
    total: 0,
    totalPages: 0,
  });

  const paginatedNovelsDataByTag = ref<PaginatedServerResponse<Novel>>({
    data: [],
    limit: 0,
    page: 0,
    total: 0,
    totalPages: 0,
  });

  const paginatedNovelsDataByAuthor = ref<PaginatedServerResponse<Novel>>({
    data: [],
    limit: 0,
    page: 0,
    total: 0,
    totalPages: 0,
  });

  const paginatedNovelsDataByGenre = ref<PaginatedServerResponse<Novel>>({
    data: [],
    limit: 0,
    page: 0,
    total: 0,
    totalPages: 0,
  });

  const fetchNovel = async (novelUpdatesId: string): Promise<void> => {
    fetchingNovel.value = true;
    novelError.value = null;
    fetch(`${url}/novels/title/${novelUpdatesId}`)
      .then((response) => response.json())
      .then((data) => {
        novel.value = data;
      })
      .catch((error) => {
        novelError.value = error;
      })
      .finally(() => {
        fetchingNovel.value = false;
      });
  };

  const fetchNovels = async (page: number, limit: number): Promise<void> => {
    fetchingNovel.value = true;
    novelError.value = null;
    fetch(`${url}/novels/?page=${page}&limit=${limit}`)
      .then((response) => response.json())
      .then((data) => {
        paginatedNovelsData.value = data;
      })
      .catch((error) => {
        novelError.value = error;
      })
      .finally(() => {
        fetchingNovel.value = false;
      });
  };

  const fetchNovelsByTag = async (
    tag: string,
    page: number,
    limit: number,
  ): Promise<void> => {
    fetchingNovel.value = true;
    novelError.value = null;
    fetch(
      `${url}/novels/tags/${encodeURIComponent(
        tag as string,
      )}?page=${page}&limit=${limit}`,
    )
      .then((response) => response.json())
      .then((data) => {
        paginatedNovelsDataByTag.value = data;
      })
      .catch((error) => {
        novelError.value = error;
      })
      .finally(() => {
        fetchingNovel.value = false;
      });
  };

  const fetchNovelsByAuthor = async (
    author: string,
    page: number,
    limit: number,
  ): Promise<void> => {
    fetchingNovel.value = true;
    novelError.value = null;
    fetch(
      `${url}/novels/authors/${encodeURIComponent(
        author as string,
      )}?page=${page}&limit=${limit}`,
    )
      .then((response) => response.json())
      .then((data) => {
        paginatedNovelsDataByAuthor.value = data;
      })
      .catch((error) => {
        novelError.value = error;
      })
      .finally(() => {
        fetchingNovel.value = false;
      });
  };

  const fetchNovelsByGenre = async (
    genre: string,
    page: number,
    limit: number,
  ): Promise<void> => {
    fetchingNovel.value = true;
    novelError.value = null;
    fetch(
      `${url}/novels/genres/${encodeURIComponent(
        genre as string,
      )}?page=${page}&limit=${limit}`,
    )
      .then((response) => response.json())
      .then((data) => {
        paginatedNovelsDataByGenre.value = data;
      })
      .catch((error) => {
        novelError.value = error;
      })
      .finally(() => {
        fetchingNovel.value = false;
      });
  };

  return {
    novel,
    fetchNovel,
    fetchNovels,
    fetchNovelsByAuthor,
    fetchNovelsByGenre,
    fetchNovelsByTag,
    fetchingNovel,
    novelError,
    paginatedNovelsData,
    paginatedNovelsDataByAuthor,
    paginatedNovelsDataByGenre,
    paginatedNovelsDataByTag,
  };
});

if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useNovelStore, import.meta.hot));
}
