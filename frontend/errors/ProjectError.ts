import { ErrorBase } from '~/errors/ErrorBase';

type ProjectErrorName = 'INTERNAL_SERVER_ERROR' | 'CREATE_NOVEL_ERROR' | 'UPDATE_NOVEL_ERROR' | 'NETWORK_ERROR' | 'VALIDATION_ERROR' | 'INVALID_RESPONSE_ERROR';

class ProjectError extends ErrorBase<'ProjectError', ProjectErrorName> {
  constructor(params: { type: ProjectErrorName; message: string; cause?: unknown }) {
    super({ name: 'ProjectError', type: params.type, message: params.message, cause: params.cause });
  }
}

export { ProjectError, type ProjectErrorName };
