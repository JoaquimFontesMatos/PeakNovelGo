import { ErrorBase } from '~/errors/ErrorBase';

type AuthErrorName = 'SIGN_UP_ERROR' | 'REFRESH_TOKEN_ERROR' | 'LOGOUT_ERROR' | 'UNAUTHORIZED_ERROR' | 'INVALID_CREDENTIALS_ERROR' | 'INVALID_SESSION_DATA';

class AuthError extends ErrorBase<'AuthError', AuthErrorName> {}

export { AuthError, type AuthErrorName };
