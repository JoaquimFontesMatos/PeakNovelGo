import { z } from 'zod';

const PaginatedServerResponseSchema = <T extends z.ZodTypeAny>(itemSchema: T) =>
  z.object({
    data: z.array(itemSchema),
    total: z.number(),
    page: z.number(),
    limit: z.number(),
    totalPages: z.number(),
  });

// Type helper for inference
type PaginatedServerResponse<T extends z.ZodType<any>> = z.infer<ReturnType<typeof PaginatedServerResponseSchema<T>>>;

export { type PaginatedServerResponse, PaginatedServerResponseSchema };
