import type { AuthErrorName } from './AuthError';
import type { BookmarkErrorName } from './BookmarkError';
import type { ChapterErrorName } from './ChapterError';
import type { NovelErrorName } from './NovelError';
import type { ProjectErrorName } from './ProjectError';
import type { TtsErrorName } from './TtsError';
import type { UserErrorName } from './UserError';

type ErrorNames = {
  AuthError: AuthErrorName;
  UserError: UserErrorName;
  ProjectError: ProjectErrorName;
  BookmarkError: BookmarkErrorName;
  NovelError: NovelErrorName;
  ChapterError: ChapterErrorName;
  TtsError: TtsErrorName;
};

class ErrorBase<N extends keyof ErrorNames, T extends ErrorNames[N]> extends Error {
  override name: N;
  readonly type: T;
  override message: string;
  override cause?: unknown;

  constructor(params: { name: N; type: T; message: string; cause?: unknown }) {
    super(params.message);
    this.name = params.name;
    this.type = params.type;
    this.message = params.message;
    this.cause = params.cause;
  }
}

export { ErrorBase, type ErrorNames };
