import { ErrorBase } from '~/errors/ErrorBase';

type ErrorName =
    | 'GET_NOVEL_ERROR'
    | 'CREATE_NOVEL_ERROR'
    | 'UPDATE_NOVEL_ERROR'
    | 'NOVEL_NOT_FOUND'

export class NovelError extends ErrorBase<ErrorName> {
}