import { AuthError } from '~/errors/AuthError';
import { ProjectError } from '~/errors/ProjectError';
import type { HttpClient } from '~/interfaces/HttpClient';
import type { ResponseParser } from '~/interfaces/ResponseParser';
import type { TtsService } from '~/interfaces/services/TtsService';
import { ErrorServerResponseSchema } from '~/schemas/ErrorServerResponse';
import { ParagraphsSchema, type Paragraph } from '~/schemas/Paragraph';
import type { TTSRequest } from '~/schemas/TTSRequest';

export class BaseTtsService implements TtsService {
  private readonly baseUrl: string;
  private readonly httpClient: HttpClient;
  private readonly responseParser: ResponseParser;

  constructor(baseUrl: string, httpClient: HttpClient, responseParser: ResponseParser) {
    this.baseUrl = baseUrl;
    this.httpClient = httpClient;
    this.responseParser = responseParser;
  }

  async generateTTS(ttsRequest: TTSRequest): Promise<Paragraph[]> {
    let response;
    let errorMessage = 'An unexpected error occurred';

    try {
      response = await this.httpClient.authorizedRequest(this.baseUrl + '/novels/tts', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(ttsRequest),
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
        case 401:
          throw new AuthError({
            type: 'UNAUTHORIZED_ERROR',
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
      const validatedParagraphsResponse = this.responseParser.validateSchema(ParagraphsSchema, parsedResponse);

      return validatedParagraphsResponse;
    } catch (validationError) {
      throw new ProjectError({
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
      response = await this.httpClient.authorizedRequest(this.baseUrl + '/novels/tts/voices', {
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
        case 401:
          throw new AuthError({
            type: 'UNAUTHORIZED_ERROR',
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
      const validatedParagraphsResponse = this.responseParser.validateSchema(ParagraphsSchema, parsedResponse);

      return validatedParagraphsResponse;
    } catch (validationError) {
      throw new ProjectError({
        type: 'VALIDATION_ERROR',
        message: 'Received malformed paragraphs data',
        cause: validationError,
      });
    }
  }
}
