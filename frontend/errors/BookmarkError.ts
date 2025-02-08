import { ErrorBase } from '~/errors/ErrorBase';

type BookmarkErrorName = 'CREATE_BOOKMARK_ERROR' | 'UPDATE_BOOKMARK_ERROR' | 'BOOKMARK_NOT_FOUND';

class BookmarkError extends ErrorBase<'BookmarkError', BookmarkErrorName> {
  constructor(params: { type: BookmarkErrorName; message: string; cause?: unknown }) {
    super({ name: 'BookmarkError', type: params.type, message: params.message, cause: params.cause });
  }
}

export { BookmarkError, type BookmarkErrorName };
