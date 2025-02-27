import { ErrorBase } from '~/errors/ErrorBase';

type ChapterErrorName = 'GET_CHAPTER_ERROR' | 'CREATE_CHAPTER_ERROR' | 'UPDATE_CHAPTER_ERROR' | 'NO_CHAPTER_FOUND';

class ChapterError extends ErrorBase<'ChapterError', ChapterErrorName> {
  constructor(params: { type: ChapterErrorName; message: string; cause?: unknown }) {
    super({ name: 'ChapterError', type: params.type, message: params.message, cause: params.cause });
  }
}

export { ChapterError, type ChapterErrorName };
