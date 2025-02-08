import { z } from 'zod';

const SuccessServerResponseSchema = z.object({
  message: z.string().nonempty('Message is required'),
});

type SuccessServerResponse = z.infer<typeof SuccessServerResponseSchema>;
export { type SuccessServerResponse, SuccessServerResponseSchema };
