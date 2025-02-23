import { NovelError } from '~/errors/NovelError';
import { ProjectError } from '~/errors/ProjectError';
import type { HttpClient } from '~/interfaces/HttpClient';
import type { ResponseParser } from '~/interfaces/ResponseParser';
import type { NovelService } from '~/interfaces/services/NovelService';
import { ErrorServerResponseSchema } from '~/schemas/ErrorServerResponse';
import type { ImportedNovel } from '~/schemas/ImportedNovel';
import { NovelSchema, PaginatedNovelsSchema, type Novel } from '~/schemas/Novel';
import type { PaginatedServerResponse } from '~/schemas/PaginatedServerResponse';

export class BaseNovelService implements NovelService {
  private readonly baseUrl: string;
  private readonly httpClient: HttpClient;
  private readonly responseParser: ResponseParser;

  constructor(baseUrl: string, httpClient: HttpClient, responseParser: ResponseParser) {
    this.baseUrl = baseUrl;
    this.httpClient = httpClient;
    this.responseParser = responseParser;
  }

  async fetchNovel(novelUpdatesId: string): Promise<Novel> {
    let response;
    let errorMessage = 'An unexpected error occurred';

    try {
      response = await this.httpClient.request(this.baseUrl + '/novels/title/' + novelUpdatesId, {
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
          throw new NovelError({
            type: 'NOVEL_NOT_FOUND',
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
      const successResponse = this.responseParser.validateSchema(NovelSchema, parsedResponse);

      return successResponse;
    } catch (validationError) {
      throw new ProjectError({
        type: 'VALIDATION_ERROR',
        message: 'Received malformed novel data',
        cause: validationError,
      });
    }
  }

  async fetchNovels(page: number, limit: number): Promise<PaginatedServerResponse<typeof NovelSchema>> {
    let response;
    let errorMessage = 'An unexpected error occurred';

    try {
      response = await this.httpClient.request(this.baseUrl + '/novels/?page=' + page + '&limit=' + limit, {
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
          throw new NovelError({
            type: 'NOVEL_NOT_FOUND',
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
      const successResponse = this.responseParser.validateSchema(PaginatedNovelsSchema, parsedResponse);

      return successResponse;
    } catch (validationError) {
      throw new ProjectError({
        type: 'VALIDATION_ERROR',
        message: 'Received malformed novel data',
        cause: validationError,
      });
    }
  }

  async fetchNovelsByTag(tag: string, page: number, limit: number): Promise<PaginatedServerResponse<typeof NovelSchema>> {
    let response;
    let errorMessage = 'An unexpected error occurred';

    try {
      response = await this.httpClient.request(this.baseUrl + '/novels/tags/' + tag + '/?page=' + page + '&limit=' + limit, {
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
        const validatedResponse = ErrorServerResponseSchema.parse(parsedResponse);
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
          throw new NovelError({
            type: 'NOVEL_NOT_FOUND',
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
      const successResponse = this.responseParser.validateSchema(PaginatedNovelsSchema, parsedResponse);

      return successResponse;
    } catch (validationError) {
      throw new ProjectError({
        type: 'VALIDATION_ERROR',
        message: 'Received malformed novel data',
        cause: validationError,
      });
    }
  }

  async fetchNovelsByGenre(genre: string, page: number, limit: number): Promise<PaginatedServerResponse<typeof NovelSchema>> {
    let response;
    let errorMessage = 'An unexpected error occurred';

    try {
      response = await this.httpClient.request(this.baseUrl + '/novels/genres/' + genre + '/?page=' + page + '&limit=' + limit, {
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
        const validatedResponse = ErrorServerResponseSchema.parse(parsedResponse);
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
          throw new NovelError({
            type: 'NOVEL_NOT_FOUND',
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
      const successResponse = this.responseParser.validateSchema(PaginatedNovelsSchema, parsedResponse);

      return successResponse;
    } catch (validationError) {
      throw new ProjectError({
        type: 'VALIDATION_ERROR',
        message: 'Received malformed novel data',
        cause: validationError,
      });
    }
  }

  async fetchNovelsByAuthor(author: string, page: number, limit: number): Promise<PaginatedServerResponse<typeof NovelSchema>> {
    let response;
    let errorMessage = 'An unexpected error occurred';

    try {
      response = await fetch(this.baseUrl + '/novels/authors/' + author + '/?page=' + page + '&limit=' + limit, {
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
        const validatedResponse = ErrorServerResponseSchema.parse(parsedResponse);
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
          throw new NovelError({
            type: 'NOVEL_NOT_FOUND',
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
      const successResponse = this.responseParser.validateSchema(PaginatedNovelsSchema, parsedResponse);

      return successResponse;
    } catch (validationError) {
      throw new ProjectError({
        type: 'VALIDATION_ERROR',
        message: 'Received malformed novel data',
        cause: validationError,
      });
    }
  }

  async importNovel(importedNovel: ImportedNovel): Promise<Novel> {
    let response;
    let errorMessage = 'An unexpected error occurred';

    try {
      response = await this.httpClient.authorizedRequest(this.baseUrl + '/novels/', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body:  JSON.stringify(importedNovel),
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
          throw new NovelError({
            type: 'NOVEL_NOT_FOUND',
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
      const successResponse = this.responseParser.validateSchema(NovelSchema, parsedResponse);

      return successResponse;
    } catch (validationError) {
      throw new ProjectError({
        type: 'VALIDATION_ERROR',
        message: 'Received malformed novel data',
        cause: validationError,
      });
    }
  }

  async importByNovelUpdatesId(novelUpdatesId: string): Promise<Novel> {
    let response;
    let errorMessage = 'An unexpected error occurred';

    try {
      response = await this.httpClient.authorizedRequest(this.baseUrl + '/novels/' + novelUpdatesId, {
        method: 'POST',
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
          throw new NovelError({
            type: 'NOVEL_NOT_FOUND',
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
      const successResponse = this.responseParser.validateSchema(NovelSchema, parsedResponse);

      return successResponse;
    } catch (validationError) {
      throw new ProjectError({
        type: 'VALIDATION_ERROR',
        message: 'Received malformed novel data',
        cause: validationError,
      });
    }
  }
}
