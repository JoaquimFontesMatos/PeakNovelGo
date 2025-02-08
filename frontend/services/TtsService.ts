import { AuthError } from '~/errors/AuthError';
import { ProjectError } from '~/errors/ProjectError';
import { ErrorServerResponseSchema } from '~/schemas/ErrorServerResponse';
import { ParagraphsSchema, type Paragraph } from '~/schemas/Paragraph';
import type { TTSRequest } from '~/schemas/TTSRequest';

export class TtsService {
  private readonly baseUrl: string;
  private authStore;

  constructor(baseUrl: string) {
    this.baseUrl = baseUrl;
    this.authStore = useAuthStore();
  }

  async generateTTS(ttsRequest: TTSRequest): Promise<Paragraph[]> {
    let response;
    let errorMessage = 'An unexpected error occurred';

    try {
      response = await this.authStore.authorizedFetch(this.baseUrl + '/novels/tts', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(ttsRequest),
      });
    } catch (error) {
      throw new ProjectError({
        name: 'ProjectError',
        type: 'NETWORK_ERROR',
        message: 'Network request failed',
        cause: error,
      });
    }

    const parsedResponse = await parseJSONPromise(response);

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
            name: 'ProjectError',
            type: 'VALIDATION_ERROR',
            message: errorMessage,
            cause: response,
          });
        case 401:
          throw new AuthError({
            name: 'AuthError',
            type: 'UNAUTHORIZED_ERROR',
            message: errorMessage,
            cause: response,
          });
        default:
          throw new ProjectError({
            name: 'ProjectError',
            type: 'INTERNAL_SERVER_ERROR',
            message: errorMessage,
            cause: response,
          });
      }
    }

    try {
      const validatedParagraphsResponse = ParagraphsSchema.parse(parsedResponse);

      return validatedParagraphsResponse;
    } catch (validationError) {
      throw new ProjectError({
        name: 'ProjectError',
        type: 'VALIDATION_ERROR',
        message: 'Received malformed paragraphs data',
        cause: validationError,
      });
    }
  }

  async fetchTTSVoices(): Promise<any> {
    let response;
    let errorMessage = 'An unexpected error occurred';

    try {
      response = await this.authStore.authorizedFetch(this.baseUrl + '/novels/tts/voices', {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
        },
      });
    } catch (error) {
      throw new ProjectError({
        name: 'ProjectError',
        type: 'NETWORK_ERROR',
        message: 'Network request failed',
        cause: error,
      });
    }

    const parsedResponse = await parseJSONPromise(response);

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
            name: 'ProjectError',
            type: 'VALIDATION_ERROR',
            message: errorMessage,
            cause: response,
          });
        case 401:
          throw new AuthError({
            name: 'AuthError',
            type: 'UNAUTHORIZED_ERROR',
            message: errorMessage,
            cause: response,
          });
        default:
          throw new ProjectError({
            name: 'ProjectError',
            type: 'INTERNAL_SERVER_ERROR',
            message: errorMessage,
            cause: response,
          });
      }
    }

    try {
      const validatedParagraphsResponse = ParagraphsSchema.parse(parsedResponse);

      return validatedParagraphsResponse;
    } catch (validationError) {
      throw new ProjectError({
        name: 'ProjectError',
        type: 'VALIDATION_ERROR',
        message: 'Received malformed paragraphs data',
        cause: validationError,
      });
    }
  }
}
