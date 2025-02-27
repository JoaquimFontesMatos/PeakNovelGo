import { ChapterError } from '~/errors/ChapterError';
import { ProjectError } from '~/errors/ProjectError';
import type { HttpClient } from '~/interfaces/HttpClient';
import type { ResponseParser } from '~/interfaces/ResponseParser';
import type { ChapterService } from '~/interfaces/services/ChapterService';
import { ChapterSchema, PaginatedChaptersSchema, type Chapter } from '~/schemas/Chapter';
import { ErrorServerResponseSchema } from '~/schemas/ErrorServerResponse';
import type { PaginatedServerResponse } from '~/schemas/PaginatedServerResponse';

export class BaseChapterService implements ChapterService {
  private readonly baseUrl: string;
  private readonly httpClient;
  private readonly responseParser;

  constructor(baseUrl: string, httpClient: HttpClient, responseParser: ResponseParser) {
    this.baseUrl = baseUrl;
    this.httpClient = httpClient;
    this.responseParser = responseParser;
  }

  async fetchChapter(novelUpdatesId: string, chaptNo: number): Promise<Chapter> {
    let response;
    let errorMessage = 'An unexpected error occurred';

    try {
      response = await this.httpClient.request(this.baseUrl + '/novels/chapters/novel/' + novelUpdatesId + '/chapter/' + chaptNo, {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
        },
      });
    } catch (error) {
      throw new ProjectError({
        type: 'NETWORK_ERROR',
        message: 'Network request failed',
        cause: error,
      });
    }

    const parsedResponse = await this.responseParser.parseJSON(response);

    if (!response.ok) {
      try {
        const validatedResponse = this.responseParser.validateSchema(ErrorServerResponseSchema, parsedResponse);
        errorMessage = validatedResponse.error;
      } catch (validationError) {
        console.log(validationError);
        // If validation fails, keep the default error message.
        // Optionally, you could log validationError for debugging.
      }

      switch (response.status) {
        case 400:
          throw new ProjectError({
            type: 'VALIDATION_ERROR',
            message: errorMessage,
            cause: response,
          });
        case 404:
          throw new ChapterError({
            type: 'NO_CHAPTER_FOUND',
            message: errorMessage,
            cause: response,
          });
        default:
          throw new ProjectError({
            type: 'INTERNAL_SERVER_ERROR',
            message: errorMessage,
            cause: response,
          });
      }
    }

    try {
      const successResponse = this.responseParser.validateSchema(ChapterSchema, parsedResponse);

      return successResponse;
    } catch (validationError) {
      throw new ProjectError({
        type: 'VALIDATION_ERROR',
        message: 'Received malformed chapter data',
        cause: validationError,
      });
    }
  }

  async fetchChapters(novelUpdatesId: string, page: number, limit: number): Promise<PaginatedServerResponse<typeof ChapterSchema>> {
    let response;
    let errorMessage = 'An unexpected error occurred';

    try {
      response = await this.httpClient.request(this.baseUrl + '/novels/chapters/novel/' + novelUpdatesId + '/chapters?page=' + page + '&limit=' + limit, {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
        },
      });
    } catch (error) {
      throw new ProjectError({
        type: 'NETWORK_ERROR',
        message: 'Network request failed',
        cause: error,
      });
    }

    const parsedResponse = await this.responseParser.parseJSON(response);

    if (!response.ok) {
      try {
        const validatedResponse = this.responseParser.validateSchema(ErrorServerResponseSchema, parsedResponse);
        errorMessage = validatedResponse.error;
      } catch (validationError) {
        console.log(validationError);
        // If validation fails, keep the default error message.
        // Optionally, you could log validationError for debugging.
      }

      switch (response.status) {
        case 400:
          throw new ProjectError({
            type: 'VALIDATION_ERROR',
            message: errorMessage,
            cause: response,
          });
        case 404:
          throw new ChapterError({
            type: 'NO_CHAPTER_FOUND',
            message: errorMessage,
            cause: response,
          });
        default:
          throw new ProjectError({
            type: 'INTERNAL_SERVER_ERROR',
            message: errorMessage,
            cause: response,
          });
      }
    }

    try {
      const successResponse = this.responseParser.validateSchema(PaginatedChaptersSchema, parsedResponse);

      return successResponse;
    } catch (validationError) {
      throw new ProjectError({
        type: 'VALIDATION_ERROR',
        message: 'Received malformed chapter data',
        cause: validationError,
      });
    }
  }
}
