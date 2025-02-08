import { ErrorBase } from '~/errors/ErrorBase';

type NovelErrorName = 'GET_NOVEL_ERROR' | 'CREATE_NOVEL_ERROR' | 'UPDATE_NOVEL_ERROR' | 'NOVEL_NOT_FOUND';

class NovelError extends ErrorBase<'NovelError', NovelErrorName> {
  constructor(params: { type: NovelErrorName; message: string; cause?: unknown }) {
    super({ name: 'NovelError', type: params.type, message: params.message, cause: params.cause });
  }
}

export { NovelError, type NovelErrorName };
