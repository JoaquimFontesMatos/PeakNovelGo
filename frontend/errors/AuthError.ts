import { ErrorBase } from '~/errors/ErrorBase';

type AuthErrorName = 'SIGN_UP_ERROR' | 'REFRESH_TOKEN_ERROR' | 'LOGOUT_ERROR' | 'UNAUTHORIZED_ERROR' | 'INVALID_CREDENTIALS_ERROR' | 'INVALID_SESSION_DATA';

class AuthError extends ErrorBase<'AuthError', AuthErrorName> {
  constructor(params: { type: AuthErrorName; message: string; cause?: unknown }) {
    super({ name: 'AuthError', type: params.type, message: params.message, cause: params.cause });
  }
}

export { AuthError, type AuthErrorName };
