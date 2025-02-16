import { type SuccessServerResponse, SuccessServerResponseSchema } from '~/schemas/SuccessServerResponse';
import { ProjectError } from '~/errors/ProjectError';
import { type ErrorServerResponse, ErrorServerResponseSchema } from '~/schemas/ErrorServerResponse';
import { AuthError } from '~/errors/AuthError';
import { UserError } from '~/errors/UserError';
import type { ResponseParser } from '~/interfaces/ResponseParser';
import type { UserService } from '~/interfaces/services/UserService';
import type { HttpClient } from '~/interfaces/HttpClient';

export class BaseUserService implements UserService {
  private readonly baseUrl: string;
  private readonly httpClient: HttpClient;
  private readonly responseParser: ResponseParser;

  constructor(baseUrl: string, httpClient: HttpClient, responseParser: ResponseParser) {
    this.baseUrl = baseUrl;
    this.httpClient = httpClient;
    this.responseParser = responseParser;
  }

  async updateUserFields(fields: {}, userId: number): Promise<SuccessServerResponse> {
    let response;
    let errorMessage = 'An unexpected error occurred';

    try {
      response = await this.httpClient.authorizedRequest(this.baseUrl + '/user/' + userId + '/fields', {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(fields),
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
        const validatedResponse: ErrorServerResponse = this.responseParser.validateSchema(ErrorServerResponseSchema, parsedResponse);
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
        case 401:
          throw new AuthError({
            type: 'UNAUTHORIZED_ERROR',
            message: errorMessage,
            cause: response,
          });
        case 403:
          throw new UserError({
            type: 'USER_DEACTIVATED_ERROR',
            message: errorMessage,
            cause: response,
          });
        case 404:
          throw new UserError({
            type: 'USER_NOT_FOUND_ERROR',
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
      const validatedSuccessResponse: SuccessServerResponse = this.responseParser.validateSchema(SuccessServerResponseSchema, parsedResponse);

      return validatedSuccessResponse as SuccessServerResponse;
    } catch (validationError) {
      console.log(validationError);
      throw new ProjectError({
        type: 'INVALID_RESPONSE_ERROR',
        message: 'Received malformed data',
        cause: validationError,
      });
    }
  }

  async deleteUser(userId: number): Promise<SuccessServerResponse> {
    let response;
    let errorMessage: string = 'An unexpected error occurred';

    try {
      response = await this.httpClient.authorizedRequest(`${this.baseUrl}/user/${userId}`, {
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
        const validatedResponse: ErrorServerResponse = this.responseParser.validateSchema(ErrorServerResponseSchema, parsedResponse);
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
        case 401:
          throw new AuthError({
            type: 'UNAUTHORIZED_ERROR',
            message: errorMessage,
            cause: response,
          });
        case 403:
          throw new UserError({
            type: 'USER_DEACTIVATED_ERROR',
            message: errorMessage,
            cause: response,
          });
        case 404:
          throw new UserError({
            type: 'USER_NOT_FOUND_ERROR',
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
      const validatedSuccessResponse: SuccessServerResponse = this.responseParser.validateSchema(SuccessServerResponseSchema, parsedResponse);

      return validatedSuccessResponse as SuccessServerResponse;
    } catch (validationError) {
      console.log(validationError);
      throw new ProjectError({
        type: 'INVALID_RESPONSE_ERROR',
        message: 'Received malformed data',
        cause: validationError,
      });
    }
  }
}
