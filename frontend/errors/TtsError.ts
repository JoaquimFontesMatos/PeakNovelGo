import { ErrorBase } from '~/errors/ErrorBase';

type TtsErrorName = 'GET_TTS_ERROR' | 'CREATE_TTS_ERROR' | 'UPDATE_TTS_ERROR';

class TtsError extends ErrorBase<'TtsError', TtsErrorName> {}

export { TtsError, type TtsErrorName };
