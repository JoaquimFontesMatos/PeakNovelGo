import { ErrorBase } from '~/errors/ErrorBase';

type ErrorName =
    | 'SIGN_UP_ERROR'
    | 'REFRESH_TOKEN_ERROR'
    | 'LOGOUT_ERROR'
    | 'UNAUTHORIZED_ERROR'
    | 'INVALID_CREDENTIALS_ERROR'
    | 'INVALID_SESSION_DATA'

export class AuthError extends ErrorBase<ErrorName> {
}