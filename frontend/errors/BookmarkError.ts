import { ErrorBase } from '~/errors/ErrorBase';

type ErrorName =
    | 'CREATE_BOOKMARK_ERROR'
    | 'UPDATE_BOOKMARK_ERROR'
    | 'BOOKMARK_NOT_FOUND'

export class BookmarkError extends ErrorBase<ErrorName> {
}