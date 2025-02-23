import z from 'zod';

const ImportedNovelSchema = z.object({
  title: z.string(),
  description: z.string(),
  image: z.string(),
  language: z.object({
    name: z.string(),
  }),
  status: z.string(),
  tags: z.array(
    z.object({
      name: z.string(),
      description: z.string(),
      id: z.number().default(0),
    })
  ),
  authors: z.array(
    z.object({
      name: z.string(),
      description: z.string(),
      id: z.number().default(0),
    })
  ),
  genres: z.array(
    z.object({
      name: z.string(),
      description: z.string(),
      id: z.number().default(0),
    })
  ),
  year: z.string(),
  release_freq: z.string(),
});

type ImportedNovel = z.infer<typeof ImportedNovelSchema>;

export { type ImportedNovel, ImportedNovelSchema };
