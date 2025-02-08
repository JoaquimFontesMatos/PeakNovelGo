import { ErrorBase } from '~/errors/ErrorBase';

type ProjectErrorName = 'INTERNAL_SERVER_ERROR' | 'CREATE_NOVEL_ERROR' | 'UPDATE_NOVEL_ERROR' | 'NETWORK_ERROR' | 'VALIDATION_ERROR' | 'INVALID_RESPONSE_ERROR';

class ProjectError extends ErrorBase<'ProjectError', ProjectErrorName> {}

export { ProjectError, type ProjectErrorName };
