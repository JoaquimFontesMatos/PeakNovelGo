import { z } from 'zod';

const ParagraphSchema = z.object({
  text: z.string(),
  index: z.number(),
  url: z.string(),
});

const ParagraphsSchema = z.array(ParagraphSchema);

type Paragraph = z.infer<typeof ParagraphSchema>;
export { type Paragraph, ParagraphSchema, ParagraphsSchema };
