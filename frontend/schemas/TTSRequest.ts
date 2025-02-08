import { z } from 'zod';

const TTSRequestSchema = z.object({
  text: z.string(),
  voice: z.string(),
  rate: z.number(),
  novelId: z.number(),
  chapterNo: z.number(),
});

type TTSRequest = z.infer<typeof TTSRequestSchema>;
export { type TTSRequest, TTSRequestSchema };
