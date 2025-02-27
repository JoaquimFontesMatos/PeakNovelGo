import { ErrorBase } from '~/errors/ErrorBase';

type TtsErrorName = 'GET_TTS_ERROR' | 'CREATE_TTS_ERROR' | 'UPDATE_TTS_ERROR';

class TtsError extends ErrorBase<'TtsError', TtsErrorName> {
  constructor(params: { type: TtsErrorName; message: string; cause?: unknown }) {
    super({ name: 'TtsError', type: params.type, message: params.message, cause: params.cause });
  }
}

export { TtsError, type TtsErrorName };
