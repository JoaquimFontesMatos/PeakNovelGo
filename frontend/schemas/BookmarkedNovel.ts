import { z } from 'zod';
import { PaginatedServerResponseSchema } from './PaginatedServerResponse';

const BookmarkedNovelSchema = z.object({
  ID: z.number(),
  CreatedAt: z.string(),
  UpdatedAt: z.string(),
  DeletedAt: z.string().nullable().optional(),

  novelId: z.number(),
  userId: z.number(),
  status: z.string(),
  score: z.number(),
  currentChapter: z.number(),
});

const PaginatedBookmarkedNovelsSchema = PaginatedServerResponseSchema<typeof BookmarkedNovelSchema>(BookmarkedNovelSchema);

type BookmarkedNovel = z.infer<typeof BookmarkedNovelSchema>;
export { type BookmarkedNovel, BookmarkedNovelSchema, PaginatedBookmarkedNovelsSchema };
