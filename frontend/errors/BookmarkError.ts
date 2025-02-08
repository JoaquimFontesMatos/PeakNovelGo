import { ErrorBase } from '~/errors/ErrorBase';

type BookmarkErrorName = 'CREATE_BOOKMARK_ERROR' | 'UPDATE_BOOKMARK_ERROR' | 'BOOKMARK_NOT_FOUND';

class BookmarkError extends ErrorBase<'BookmarkError', BookmarkErrorName> {}

export { BookmarkError, type BookmarkErrorName };
