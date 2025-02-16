import type { AuthService } from '~/interfaces/services/AuthService';
import type { BookmarkService } from '~/interfaces/services/BookmarkService';
import type { ChapterService } from '~/interfaces/services/ChapterService';
import type { NovelService } from '~/interfaces/services/NovelService';
import type { TtsService } from '~/interfaces/services/TtsService';
import type { UserService } from '~/interfaces/services/UserService';
import { BaseAuthService } from '~/services/AuthService';
import { BaseBookmarkService } from '~/services/BookmarkService';
import { BaseChapterService } from '~/services/ChapterService';
import { BaseNovelService } from '~/services/NovelService';
import { BaseTtsService } from '~/services/TtsService';
import { BaseUserService } from '~/services/UserService';

export default defineNuxtPlugin(nuxtApp => {
  const runtimeConfig = useRuntimeConfig();
  const url = runtimeConfig.public.apiUrl;

  const httpClient = new FetchHttpClient();
  const responseParser = new ZodResponseParser();

  const authService: AuthService = new BaseAuthService(url, httpClient, responseParser);
  const userService: UserService = new BaseUserService(url, httpClient, responseParser);
  const novelService: NovelService = new BaseNovelService(url, httpClient, responseParser);
  const chapterService: ChapterService = new BaseChapterService(url, httpClient, responseParser);
  const ttsService: TtsService = new BaseTtsService(url, httpClient, responseParser);
  const bookmarkService: BookmarkService = new BaseBookmarkService(url, httpClient, responseParser);

  const errorHandler = new BaseErrorHandler();

  nuxtApp.provide('authService', authService);
  nuxtApp.provide('novelService', novelService);
  nuxtApp.provide('chapterService', chapterService);
  nuxtApp.provide('ttsService', ttsService);
  nuxtApp.provide('bookmarkService', bookmarkService);
  nuxtApp.provide('userService', userService);
  nuxtApp.provide('errorHandler', errorHandler);
});
