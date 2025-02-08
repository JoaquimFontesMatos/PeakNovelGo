import { ErrorBase } from '~/errors/ErrorBase';

type UserErrorName =
  | 'GET_USER_ERROR'
  | 'CREATE_USER_ERROR'
  | 'UPDATE_USER_ERROR'
  | 'INVALID_USER_DATA'
  | 'USER_DEACTIVATED_ERROR'
  | 'USER_CONFLICT_ERROR'
  | 'USER_NOT_FOUND_ERROR';

class UserError extends ErrorBase<'UserError', UserErrorName> {
  constructor(params: { type: UserErrorName; message: string; cause?: unknown }) {
    super({ name: 'UserError', type: params.type, message: params.message, cause: params.cause });
  }
}

export { UserError, type UserErrorName };
