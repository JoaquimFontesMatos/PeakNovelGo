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

  const updatingNovels = ref(false);
  const novelStatuses = ref<Record<string, string>>({});

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

  const batchUpdateNovels = async () => {
    updatingNovels.value = true; // Set importing state to true
    novelStatuses.value = {}; // Reset statuses

    const eventSourceUrl = `${url}/novels/update`;

    // Use your authorizedRequest function to fetch the SSE stream
    const response = await httpClient.authorizedRequest(eventSourceUrl, {
      headers: {
        Accept: 'text/event-stream',
      },
    });

    if (!response.ok) {
      console.error('Failed to establish SSE connection:', response.statusText);
      updatingNovels.value = false;
      return;
    }

    // Read the stream
    const reader = response.body?.getReader();
    const decoder = new TextDecoder();

    if (!reader) {
      console.error('Failed to read SSE stream');
      updatingNovels.value = false; // Reset importing state
      return;
    }

    const processStream = async () => {
      let buffer = ''; // Buffer to accumulate incomplete data

      while (true) {
        const { done, value } = await reader.read();

        if (done) {
          console.log('SSE stream closed');
          updatingNovels.value = false; // Reset importing state
          return;
        }

        const chunk = decoder.decode(value); // Decode the chunk
        buffer += chunk; // Append the chunk to the buffer

        // Split the buffer by double newlines (SSE events are separated by \n\n)
        const events = buffer.split('\n\n');

        // Process all complete events (leave the last one in the buffer if incomplete)
        for (let i = 0; i < events.length - 1; i++) {
          const event = events[i].trim(); // Remove any leading/trailing whitespace

          if (!event) {
            // Skip empty events
            continue;
          }

          // Split the event into individual lines
          const lines = event.split('\n');

          let eventType = ''; // To store the event type (e.g., "status")
          let eventData = ''; // To store the JSON data

          // Process each line in the event
          for (const line of lines) {
            if (line.startsWith('event:')) {
              // Extract the event type
              eventType = line.replace('event:', '').trim();
            } else if (line.startsWith('data:')) {
              // Extract the JSON data
              eventData = line.replace('data:', '').trim();
            }
          }

          // Only process if the event type is "status" and data is present
          if (eventType === 'status' && eventData) {
            try {
              const statuses = JSON.parse(eventData) as Record<string, string>;
              novelStatuses.value = { ...novelStatuses.value, ...statuses }; // Update statuses
            } catch (error) {
              console.error('Failed to parse status update:', error);
              console.error('Invalid JSON data:', eventData); // Log the invalid data
            }
          } else {
            console.warn('Skipping non-status event or empty data:', event);
          }
        }

        // Keep the last (possibly incomplete) event in the buffer
        buffer = events[events.length - 1];
      }
    };

    await processStream();
  };

  return {
    novel,
    fetchingNovel,
    paginatedNovelsData,
    paginatedNovelsDataByAuthor,
    paginatedNovelsDataByGenre,
    paginatedNovelsDataByTag,
    importingNovel,
    updatingNovels,
    novelStatuses,
    fetchNovel,
    fetchNovels,
    fetchNovelsByAuthor,
    fetchNovelsByGenre,
    fetchNovelsByTag,
    importNovel,
    importByNovelUpdatesId,
    batchUpdateNovels
  };
});

if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useNovelStore, import.meta.hot));
}
