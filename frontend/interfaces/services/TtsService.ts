import type { Paragraph } from '~/schemas/Paragraph';
import type { TTSRequest } from '~/schemas/TTSRequest';

export interface TtsService {
  generateTTS(ttsRequest: TTSRequest): Promise<Paragraph[]>;
  fetchTTSVoices(): Promise<any>;
}
