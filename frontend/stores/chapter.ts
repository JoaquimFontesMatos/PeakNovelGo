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
  const chapterStatuses = ref<Record<string, string>>({}); // Track chapter import statuses
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

    // Use your authorizedRequest function to fetch the SSE stream
    const response = await httpClient.authorizedRequest(eventSourceUrl, {
      headers: {
        Accept: 'text/event-stream', // Ensure the server knows you want SSE
      },
    });

    if (!response.ok) {
      console.error('Failed to establish SSE connection:', response.statusText);
      importingChapters.value = false; // Reset importing state
      return;
    }

    // Read the stream
    const reader = response.body?.getReader();
    const decoder = new TextDecoder();

    if (!reader) {
      console.error('Failed to read SSE stream');
      importingChapters.value = false; // Reset importing state
      return;
    }

    const processStream = async () => {
      let buffer = ''; // Buffer to accumulate incomplete data

      while (true) {
        const { done, value } = await reader.read();

        if (done) {
          console.log('SSE stream closed');
          importingChapters.value = false; // Reset importing state
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
              chapterStatuses.value = { ...chapterStatuses.value, ...statuses }; // Update statuses
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

    processStream();
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
