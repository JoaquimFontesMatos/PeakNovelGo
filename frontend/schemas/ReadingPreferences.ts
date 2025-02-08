import { z } from 'zod';

// Define the Zod schema for the TTS part of ReadingPreferences
const TtsSchema = z.object({
  autoplay: z.boolean().nullable().optional(),
  voice: z.string().nullable().optional(),
  rate: z.number().nullable().optional(),
});

// Define the Zod schema for ReadingPreferences
const ReadingPreferencesSchema = z.object({
  atomicReading: z.boolean().nullable().optional(),
  font: z.string().nullable().optional(),
  theme: z.string().nullable().optional(),
  tts: TtsSchema,
});

type Tts = z.infer<typeof TtsSchema>;
type ReadingPreferences = z.infer<typeof ReadingPreferencesSchema>;

export { type ReadingPreferences, type Tts, TtsSchema, ReadingPreferencesSchema };
