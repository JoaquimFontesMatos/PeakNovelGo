import { ErrorBase } from '~/errors/ErrorBase';

type ErrorName = 'GET_CHAPTER_ERROR' | 'CREATE_CHAPTER_ERROR' | 'UPDATE_CHAPTER_ERROR' | 'NO_CHAPTER_FOUND';

export class ChapterError extends ErrorBase<ErrorName> {}
