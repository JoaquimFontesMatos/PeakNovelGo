import { z } from 'zod';
import { PaginatedServerResponseSchema, type PaginatedServerResponse } from './PaginatedServerResponse';

const TagSchema = z.object({
  id: z.number(),
  name: z.string(),
  description: z.string(),
});

const AuthorSchema = z.object({
  id: z.number(),
  name: z.string(),
});

const GenreSchema = z.object({
  id: z.number(),
  name: z.string(),
  description: z.string(),
});

const NovelSchema = z.object({
  ID: z.number(),
  CreatedAt: z.string(),
  UpdatedAt: z.string(),
  DeletedAt: z.string().nullable().optional(),

  title: z.string(),
  synopsis: z.string(),
  coverUrl: z.string(),
  language: z.string(),
  status: z.string(),
  novelUpdatesUrl: z.string(),
  tags: z.array(TagSchema),
  authors: z.array(AuthorSchema),
  genres: z.array(GenreSchema),
  year: z.string(),
  releaseFrequency: z.string(),
  novelUpdatesId: z.string(),
  latestChapter: z.number(),
});

const PaginatedNovelsSchema = PaginatedServerResponseSchema<typeof NovelSchema>(NovelSchema);

type Tag = z.infer<typeof TagSchema>;
type Author = z.infer<typeof AuthorSchema>;
type Genre = z.infer<typeof GenreSchema>;
type Novel = z.infer<typeof NovelSchema>;
export { type Tag, type Author, type Genre, type Novel, TagSchema, AuthorSchema, GenreSchema, NovelSchema, PaginatedNovelsSchema };
