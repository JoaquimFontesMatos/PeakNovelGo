import { ErrorBase } from '~/errors/ErrorBase';

type ErrorName =
    | 'INTERNAL_SERVER_ERROR'
    | 'CREATE_NOVEL_ERROR'
    | 'UPDATE_NOVEL_ERROR'
    | 'NETWORK_ERROR'

export class ProjectError extends ErrorBase<ErrorName> {
}