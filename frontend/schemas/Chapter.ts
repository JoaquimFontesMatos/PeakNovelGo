import { z } from 'zod';
import { PaginatedServerResponseSchema } from './PaginatedServerResponse';

const ChapterSchema = z.object({
  ID: z.number(),
  CreatedAt: z.string(),
  UpdatedAt: z.string(),
  DeletedAt: z.string().nullable().optional(),
  chapterNo: z.number(),
  novelId: z.number(),
  title: z.string(),
  chapterUrl: z.string(),
  body: z.string(),
});

const PaginatedChaptersSchema = PaginatedServerResponseSchema<typeof ChapterSchema>(ChapterSchema);

type Chapter = z.infer<typeof ChapterSchema>;

export { type Chapter, ChapterSchema, PaginatedChaptersSchema };
