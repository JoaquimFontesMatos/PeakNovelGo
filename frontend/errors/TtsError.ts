import { ErrorBase } from '~/errors/ErrorBase';

type ErrorName =
    | 'GET_TTS_ERROR'
    | 'CREATE_TTS_ERROR'
    | 'UPDATE_TTS_ERROR'

export class TtsError extends ErrorBase<ErrorName> {
}