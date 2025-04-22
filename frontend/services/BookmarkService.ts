import { BookmarkError } from '~/errors/BookmarkError';
import { ProjectError } from '~/errors/ProjectError';
import type { HttpClient } from '~/interfaces/HttpClient';
import type { ResponseParser } from '~/interfaces/ResponseParser';
import type { BookmarkService } from '~/interfaces/services/BookmarkService';
import { BookmarkedNovelSchema, PaginatedBookmarkedNovelsSchema, type BookmarkedNovel } from '~/schemas/BookmarkedNovel';
import { ErrorServerResponseSchema } from '~/schemas/ErrorServerResponse';
import type { PaginatedServerResponse } from '~/schemas/PaginatedServerResponse';
import { SuccessServerResponseSchema, type SuccessServerResponse } from '~/schemas/SuccessServerResponse';
import { type NovelSchema, PaginatedNovelsSchema } from '~/schemas/Novel';

export class BaseBookmarkService implements BookmarkService {
    private readonly baseUrl: string;
    private readonly httpClient: HttpClient;
    private readonly responseParser: ResponseParser;

    constructor(baseUrl: string, httpClient: HttpClient, responseParser: ResponseParser) {
        this.baseUrl = baseUrl;
        this.httpClient = httpClient;
        this.responseParser = responseParser;
    }

    async bookmarkNovel(novelId: number, userId: number): Promise<BookmarkedNovel> {
        let response;
        let errorMessage = 'An unexpected error occurred';

        const createBookmarkedNovel = {
            novelId: novelId,
            userId: userId,
            score: 0,
            status: 'Plan to Read',
            currentChapter: 1,
        };

        try {
            response = await this.httpClient.authorizedRequest(this.baseUrl + '/novels/bookmarked/', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(createBookmarkedNovel),
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
                    throw new BookmarkError({
                        type: 'BOOKMARK_NOT_FOUND',
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
            const successResponse = this.responseParser.validateSchema(BookmarkedNovelSchema, parsedResponse);

            return successResponse;
        } catch (validationError) {
            throw new ProjectError({
                type: 'VALIDATION_ERROR',
                message: 'Received malformed bookmark data',
                cause: validationError,
            });
        }
    }

    async updateBookmark(updatedBookmark: BookmarkedNovel): Promise<BookmarkedNovel> {
        let response;
        let errorMessage = 'An unexpected error occurred';

        try {
            response = await this.httpClient.authorizedRequest(this.baseUrl + '/novels/bookmarked/', {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(updatedBookmark),
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
                    throw new BookmarkError({
                        type: 'BOOKMARK_NOT_FOUND',
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
            const successResponse = this.responseParser.validateSchema(BookmarkedNovelSchema, parsedResponse);

            return successResponse;
        } catch (validationError) {
            throw new ProjectError({
                type: 'VALIDATION_ERROR',
                message: 'Received malformed bookmark data',
                cause: validationError,
            });
        }
    }

    async unbookmarkNovel(novelId: number, userId: number): Promise<SuccessServerResponse> {
        let response;
        let errorMessage = 'An unexpected error occurred';

        try {
            response = await this.httpClient.authorizedRequest(this.baseUrl + '/novels/bookmarked/user/' + userId + '/novel/' + novelId, {
                method: 'DELETE',
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
                    throw new BookmarkError({
                        type: 'BOOKMARK_NOT_FOUND',
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
            const successResponse = this.responseParser.validateSchema(SuccessServerResponseSchema, parsedResponse);

            return successResponse;
        } catch (validationError) {
            throw new ProjectError({
                type: 'VALIDATION_ERROR',
                message: 'Received malformed response data',
                cause: validationError,
            });
        }
    }

    async fetchBookmarkedNovelByUser(novelUpdatesId: string, userId: number): Promise<BookmarkedNovel> {
        let response;
        let errorMessage = 'An unexpected error occurred';

        try {
            response = await this.httpClient.authorizedRequest(this.baseUrl + '/novels/bookmarked/user/' + userId + '/novel/' + novelUpdatesId, {
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
                    throw new BookmarkError({
                        type: 'BOOKMARK_NOT_FOUND',
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
            const successResponse = this.responseParser.validateSchema(BookmarkedNovelSchema, parsedResponse);

            return successResponse;
        } catch (validationError) {
            throw new ProjectError({
                type: 'VALIDATION_ERROR',
                message: 'Received malformed bookmark data',
                cause: validationError,
            });
        }
    }

    async fetchBookmarkedNovelsByUser(userId: number, page: number, limit: number): Promise<PaginatedServerResponse<typeof NovelSchema>> {
        let response;
        let errorMessage = 'An unexpected error occurred';

        try {
            response = await this.httpClient.authorizedRequest(this.baseUrl + '/novels/bookmarked/user/' + userId + '?page=' + page + '&limit=' + limit, {
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
                    throw new BookmarkError({
                        type: 'BOOKMARK_NOT_FOUND',
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
                message: 'Received malformed bookmark data',
                cause: validationError,
            });
        }
    }
}
