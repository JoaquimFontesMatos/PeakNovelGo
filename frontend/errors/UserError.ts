import { ErrorBase } from '~/errors/ErrorBase';

type ErrorName =
    | 'GET_USER_ERROR'
    | 'CREATE_USER_ERROR'
    | 'UPDATE_USER_ERROR'
    | 'INVALID_USER_DATA'
    | 'USER_DEACTIVATED_ERROR'
    | 'USER_CONFLICT_ERROR'
    | 'USER_NOT_FOUND_ERROR'

export class UserError extends ErrorBase<ErrorName> {
}