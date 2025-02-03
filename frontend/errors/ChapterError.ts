import { ErrorBase } from '~/errors/ErrorBase';

type ErrorName =
    | 'GET_CHAPTER_ERROR'
    | 'CREATE_CHAPTER_ERROR'
    | 'UPDATE_CHAPTER_ERROR'

export class ChapterError extends ErrorBase<ErrorName> {
}