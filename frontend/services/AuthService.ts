import type { LoginForm, SignUpForm } from '~/schemas/Forms';
import { AuthError } from '~/errors/AuthError';
import { ProjectError } from '~/errors/ProjectError';
import { parseJSONPromise } from '~/utils/JsonParser';
import { AuthSessionSchema } from '~/schemas/AuthSession';
import { ErrorServerResponseSchema } from '~/schemas/ErrorServerResponse';
import type { AuthSession } from '~/models/AuthSession';
import type { ErrorServerResponse } from '~/models/ErrorServerResponse';
import { UserError } from '~/errors/UserError';
import type { SuccessServerResponse } from '~/models/SuccessServerResponse';
import { SuccessServerResponseSchema } from '~/schemas/SuccessServerResponse';

export class AuthService {
    private readonly baseUrl: string;

    constructor(baseUrl: string) {
        this.baseUrl = baseUrl;
    }

    async login(form: LoginForm): Promise<AuthSession> {
        let response;
        let errorMessage = 'An unexpected error occurred';

        try {
            response = await fetch(`${(this.baseUrl)}/auth/login`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(form),
                credentials: 'include',
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
                const validatedResponse: ErrorServerResponse = await ErrorServerResponseSchema.validate(parsedResponse);
                errorMessage = validatedResponse.error;
            } catch (validationError) {
                console.log(validationError);
                // If validation fails, keep the default error message.
                // Optionally, you could log validationError for debugging.
            }

            switch (response.status) {
                case 400:
                    throw new AuthError({
                        name: 'INVALID_CREDENTIALS_ERROR',
                        message: errorMessage,
                        cause: response,
                    });
                case 401:
                    throw new AuthError({
                        name: 'UNAUTHORIZED_ERROR',
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
            const validatedAuthSessionResponse = await AuthSessionSchema.validate(parsedResponse);

            return validatedAuthSessionResponse as AuthSession;
        } catch (validationError) {
            console.log(validationError);
            throw new AuthError({
                name: 'INVALID_SESSION_DATA',
                message: 'Received malformed authentication data',
                cause: validationError,
            });
        }
    };

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
                name: 'NETWORK_ERROR',
                message: 'Network request failed',
                cause: error,
            });
        }

        const parsedResponse = await parseJSONPromise(response);

        if (!response.ok) {
            try {
                const validatedResponse: ErrorServerResponse = await ErrorServerResponseSchema.validate(parsedResponse);
                errorMessage = validatedResponse.error;
            } catch (validationError) {
                console.log(validationError);
                // If validation fails, keep the default error message.
                // Optionally, you could log validationError for debugging.
            }

            switch (response.status) {
                case 400:
                    throw new UserError({
                        name: 'INVALID_USER_DATA',
                        message: errorMessage,
                        cause: response,
                    });
                case 409:
                    throw new UserError({
                        name: 'USER_CONFLICT_ERROR',
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
            const validatedResponse: SuccessServerResponse = await SuccessServerResponseSchema.validate(parsedResponse);

            return validatedResponse as SuccessServerResponse;
        } catch (validationError) {
            console.log(validationError);
            throw new UserError({
                name: 'INVALID_USER_DATA',
                message: 'Received malformed user data',
                cause: validationError,
            });
        }
    };

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
                name: 'NETWORK_ERROR',
                message: 'Network request failed',
                cause: error,
            });
        }

        const parsedResponse = await parseJSONPromise(response);

        if (!response.ok) {
            try {
                const validatedResponse: ErrorServerResponse = await ErrorServerResponseSchema.validate(parsedResponse);
                errorMessage = validatedResponse.error;
            } catch (validationError) {
                console.log(validationError);
                // If validation fails, keep the default error message.
                // Optionally, you could log validationError for debugging.
            }

            switch (response.status) {
                case 400:
                    throw new AuthError({
                        name: 'INVALID_CREDENTIALS_ERROR',
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
            const validatedAuthSessionResponse = await AuthSessionSchema.validate(parsedResponse);

            return validatedAuthSessionResponse as AuthSession;
        } catch (validationError) {
            console.log(validationError);
            throw new AuthError({
                name: 'INVALID_SESSION_DATA',
                message: 'Received malformed session data',
                cause: validationError,
            });
        }
    };

    async logout(): Promise<void> {
        try {
            await fetch(`${this.baseUrl}/auth/logout`, {
                method: 'POST',
                credentials: 'include',
            });
        } catch (error) {
            throw new ProjectError({
                name: 'NETWORK_ERROR',
                message: 'Network request failed',
                cause: error,
            });
        }
    }


}