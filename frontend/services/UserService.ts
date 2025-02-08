import { type SuccessServerResponse, SuccessServerResponseSchema } from '~/schemas/SuccessServerResponse';
import { ProjectError } from '~/errors/ProjectError';
import { type ErrorServerResponse, ErrorServerResponseSchema } from '~/schemas/ErrorServerResponse';
import { AuthError } from '~/errors/AuthError';
import { UserError } from '~/errors/UserError';

export class UserService {
  private readonly baseUrl: string;
  private authStore;

  constructor(baseUrl: string) {
    this.baseUrl = baseUrl;
    this.authStore = useAuthStore();
  }

  async updateUserFields(fields: {}, userId: number): Promise<SuccessServerResponse> {
    let response;
    let errorMessage = 'An unexpected error occurred';

    try {
      response = await this.authStore.authorizedFetch(this.baseUrl + '/user/' + userId + '/fields', {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(fields),
      });
    } catch (error) {
      throw new ProjectError({
        name: 'NETWORK_ERROR',
        message: 'Network request failed',
        cause: error,
      });
    }

    const parsedResponse = await parseJSONPromise(response);

    if (!response.ok) {
      try {
        const validatedResponse: ErrorServerResponse = ErrorServerResponseSchema.parse(parsedResponse);
        errorMessage = validatedResponse.error;
      } catch (validationError) {
        console.log(validationError);
        // If validation fails, keep the default error message.
        // Optionally, you could log validationError for debugging.
      }

      switch (response.status) {
        case 400:
          throw new ProjectError({
            name: 'VALIDATION_ERROR',
            message: errorMessage,
            cause: response,
          });
        case 401:
          throw new AuthError({
            name: 'UNAUTHORIZED_ERROR',
            message: errorMessage,
            cause: response,
          });
        case 403:
          throw new UserError({
            name: 'USER_DEACTIVATED_ERROR',
            message: errorMessage,
            cause: response,
          });
        case 404:
          throw new UserError({
            name: 'USER_NOT_FOUND_ERROR',
            message: errorMessage,
            cause: response,
          });
        default:
          throw new ProjectError({
            name: 'INTERNAL_SERVER_ERROR',
            message: errorMessage,
            cause: response,
          });
      }
    }

    try {
      const validatedSuccessResponse: SuccessServerResponse = SuccessServerResponseSchema.parse(parsedResponse);

      return validatedSuccessResponse as SuccessServerResponse;
    } catch (validationError) {
      console.log(validationError);
      throw new ProjectError({
        name: 'INVALID_RESPONSE_ERROR',
        message: 'Received malformed data',
        cause: validationError,
      });
    }
  }

  async deleteUser(userId: number): Promise<SuccessServerResponse> {
    let response;
    let errorMessage: string = 'An unexpected error occurred';

    try {
      response = await this.authStore.authorizedFetch(`${this.baseUrl}/user/${userId}`, {
        method: 'DELETE',
        headers: {
          'Content-Type': 'application/json',
        },
      });
    } catch (error) {
      throw new ProjectError({
        name: 'NETWORK_ERROR',
        message: 'Network request failed',
        cause: error,
      });
    }

    const parsedResponse = await parseJSONPromise(response);

    if (!response.ok) {
      try {
        const validatedResponse: ErrorServerResponse = ErrorServerResponseSchema.parse(parsedResponse);
        errorMessage = validatedResponse.error;
      } catch (validationError) {
        console.log(validationError);
        // If validation fails, keep the default error message.
        // Optionally, you could log validationError for debugging.
      }

      switch (response.status) {
        case 400:
          throw new ProjectError({
            name: 'VALIDATION_ERROR',
            message: errorMessage,
            cause: response,
          });
        case 401:
          throw new AuthError({
            name: 'UNAUTHORIZED_ERROR',
            message: errorMessage,
            cause: response,
          });
        case 403:
          throw new UserError({
            name: 'USER_DEACTIVATED_ERROR',
            message: errorMessage,
            cause: response,
          });
        case 404:
          throw new UserError({
            name: 'USER_NOT_FOUND_ERROR',
            message: errorMessage,
            cause: response,
          });
        default:
          throw new ProjectError({
            name: 'INTERNAL_SERVER_ERROR',
            message: errorMessage,
            cause: response,
          });
      }
    }

    try {
      const validatedSuccessResponse: SuccessServerResponse = SuccessServerResponseSchema.parse(parsedResponse);

      return validatedSuccessResponse as SuccessServerResponse;
    } catch (validationError) {
      console.log(validationError);
      throw new ProjectError({
        name: 'INVALID_RESPONSE_ERROR',
        message: 'Received malformed data',
        cause: validationError,
      });
    }
  }
}
