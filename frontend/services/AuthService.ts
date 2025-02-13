import type { LoginForm, SignUpForm } from '~/schemas/Forms';
import { AuthError } from '~/errors/AuthError';
import { ProjectError } from '~/errors/ProjectError';
import { parseJSONPromise } from '~/utils/JsonParser';
import { type AuthSession, AuthSessionSchema } from '~/schemas/AuthSession';
import { type ErrorServerResponse, ErrorServerResponseSchema } from '~/schemas/ErrorServerResponse';
import { UserError } from '~/errors/UserError';
import { type SuccessServerResponse, SuccessServerResponseSchema } from '~/schemas/SuccessServerResponse';

export class AuthService {
  private readonly baseUrl: string;

  constructor(baseUrl: string) {
    this.baseUrl = baseUrl;
  }

  async login(form: LoginForm): Promise<AuthSession> {
    let response;
    let errorMessage = 'An unexpected error occurred';

    try {
      response = await fetch(`${this.baseUrl}/auth/login`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(form),
        credentials: 'include',
      });
    } catch (error) {
      throw new ProjectError({
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
            type: 'VALIDATION_ERROR',
            message: errorMessage,
            cause: response,
          });
        case 401:
          throw new AuthError({
            type: 'INVALID_CREDENTIALS_ERROR',
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
      const validatedAuthSessionResponse = AuthSessionSchema.parse(parsedResponse);

      return validatedAuthSessionResponse;
    } catch (validationError) {
      throw new AuthError({
        type: 'INVALID_SESSION_DATA',
        message: 'Received malformed authentication data',
        cause: validationError,
      });
    }
  }

  async signUp(form: SignUpForm): Promise<SuccessServerResponse> {
    let response;
    let errorMessage = 'An unexpected error occurred';

    const registerData = {
      username: form.username,
      email: form.email,
      password: form.password,
      bio: 'Please edit me',
      profilePicture: 'Please edit me',
      dateOfBirth: form.dateOfBirth,
    };

    try {
      response = await fetch(`${this.baseUrl}/auth/register`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(registerData),
      });
    } catch (error) {
      throw new ProjectError({
        type: 'NETWORK_ERROR',
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
          throw new UserError({
            type: 'INVALID_USER_DATA',
            message: errorMessage,
            cause: response,
          });
        case 409:
          throw new UserError({
            type: 'USER_CONFLICT_ERROR',
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
      const validatedResponse: SuccessServerResponse = SuccessServerResponseSchema.parse(parsedResponse);

      return validatedResponse as SuccessServerResponse;
    } catch (validationError) {
      console.log(validationError);
      throw new UserError({
        type: 'INVALID_USER_DATA',
        message: 'Received malformed user data',
        cause: validationError,
      });
    }
  }

  async refreshAccessToken(): Promise<AuthSession> {
    let response;
    let errorMessage = 'An unexpected error occurred';

    try {
      response = await fetch(`${this.baseUrl}/auth/refresh-token`, {
        method: 'POST',
        credentials: 'include',
      });
    } catch (error) {
      throw new ProjectError({
        type: 'NETWORK_ERROR',
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
          throw new AuthError({
            type: 'INVALID_CREDENTIALS_ERROR',
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
      const validatedAuthSessionResponse = AuthSessionSchema.parse(parsedResponse);

      return validatedAuthSessionResponse as AuthSession;
    } catch (validationError) {
      console.log(validationError);
      throw new AuthError({
        type: 'INVALID_SESSION_DATA',
        message: 'Received malformed session data',
        cause: validationError,
      });
    }
  }

  async verifyToken(token: string): Promise<SuccessServerResponse> {
    let response;
    let errorMessage = 'An unexpected error occurred';

    try {
      response = await fetch(`${this.baseUrl}/auth/verify-email?token=${token}`, {
        method: 'GET',
      });
    } catch (error) {
      throw new ProjectError({
        type: 'NETWORK_ERROR',
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
          throw new AuthError({
            type: 'INVALID_CREDENTIALS_ERROR',
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
      const validatedSuccessResponse = SuccessServerResponseSchema.parse(parsedResponse);

      return validatedSuccessResponse;
    } catch (validationError) {
      console.log(validationError);
      throw new AuthError({
        type: 'INVALID_SESSION_DATA',
        message: 'Received malformed message data',
        cause: validationError,
      });
    }
  }

  async logout(): Promise<SuccessServerResponse> {
    let response;
    let errorMessage = 'An unexpected error occurred';

    try {
      response = await fetch(`${this.baseUrl}/auth/logout`, {
        method: 'POST',
        credentials: 'include',
      });
    } catch (error) {
      throw new ProjectError({
        type: 'NETWORK_ERROR',
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
      const validatedResponse: SuccessServerResponse = SuccessServerResponseSchema.parse(parsedResponse);

      return validatedResponse as SuccessServerResponse;
    } catch (validationError) {
      console.log(validationError);
      throw new ProjectError({
        type: 'VALIDATION_ERROR',
        message: 'Received malformed success data',
        cause: validationError,
      });
    }
  }
}
