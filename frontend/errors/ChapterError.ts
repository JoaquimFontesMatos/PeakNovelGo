import { ErrorBase } from '~/errors/ErrorBase';

type ChapterErrorName = 'GET_CHAPTER_ERROR' | 'CREATE_CHAPTER_ERROR' | 'UPDATE_CHAPTER_ERROR' | 'NO_CHAPTER_FOUND';

class ChapterError extends ErrorBase<'ChapterError', ChapterErrorName> {}

export { ChapterError, type ChapterErrorName };
